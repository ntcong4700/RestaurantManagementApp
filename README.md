# RestaurantManagementApp

Backend quản lý nhà hàng: bàn, thực đơn, gọi món, xuất hoá đơn (hỗ trợ chia bill) và thanh toán.

## Công nghệ

- **Go** 1.27
- **Gin**: HTTP router
- **GORM** + **PostgreSQL 17**
- **Viper**: đọc cấu hình từ `.env`
- **Docker Compose**: chạy database

## Cấu trúc thư mục

```
backend/
├── cmd/
│   ├── config/        # đọc .env, kết nối DB, AutoMigrate
│   └── server/        # main.go: điểm khởi động
├── internal/
│   ├── models/        # GORM models: Tables, MenuItem, Order, OrderItem, Bill, Payment
│   ├── service/       # business logic
│   ├── controller/    # HTTP handlers (Gin)
│   └── router/        # khai báo route
└── docker-compose.yml # PostgreSQL
data/                  # volume dữ liệu Postgres (bị git ignore)
```

## Cài đặt & chạy

### 1. Tạo file `backend/.env`

```env
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_NAME=restaurant
```

File này đã nằm trong `.gitignore`, đừng commit nó.

### 2. Khởi động database

```bash
cd backend
docker compose up -d
```

Dữ liệu được lưu ở thư mục `data/` ngoài thư mục gốc.

### 3. Chạy server

```bash
cd backend          # bắt buộc: .env được đọc theo thư mục hiện tại
go run ./cmd/server
```

Server chạy ở `http://localhost:8080`, đổi được cổng qua biến môi trường `PORT`. Khi khởi động, bảng được tạo/cập nhật tự động bằng `AutoMigrate`.

## Mô hình dữ liệu

```
Tables 1───* Order 1───* OrderItem *───1 MenuItem
               │             │
               1             * (bill_id, có thể null)
               │             │
               *───────── Bill 1───* Payment
```

| Model       | Ý nghĩa |
|-------------|---------|
| `Tables`    | Bàn. Trạng thái: `available`, `occupied`, `reserved` |
| `MenuItem`  | Món trong thực đơn, `item_number` là duy nhất |
| `Order`     | Lượt phục vụ một bàn. Trạng thái: `open` → `paid` |
| `OrderItem` | Một món trong order (số lượng > 0). `unit_price` = giá lúc gọi món (đổi giá menu sau đó không ảnh hưởng). `bill_id` null = chưa xuất hoá đơn |
| `Bill`      | Hoá đơn của một phần hoặc toàn bộ order. Tổng tiền = Σ `unit_price` × số lượng |
| `Payment`   | Một lần thanh toán cho bill, cho phép trả nhiều lần |

## Luồng nghiệp vụ

1. **Mở order**: `OrderService.CreateOrder`. Bàn phải tồn tại và chưa có order `open`; bàn chuyển sang `occupied`.
2. **Gọi món**: `OrderItemService.CreateOrderItem`. Order phải `open`, món phải có trong thực đơn.
3. **Sửa/xoá món**: chỉ được làm khi order còn `open` và món chưa nằm trong bill nào.
4. **Xuất hoá đơn**: `BillService.CreateBill(orderID, itemIDs)`
   - `itemIDs` rỗng: gom tất cả món chưa vào bill.
   - Truyền danh sách ID: chia bill theo món.
   - Huỷ bill bằng `DeleteBill` khi bill chưa có thanh toán; các món quay lại trạng thái chưa vào bill.
5. **Thanh toán**: `PaymentService.CreatePayment(billID, amount)`. Cho phép trả một phần, không cho trả vượt số còn lại.
6. **Tự động đóng order**: khi mọi món đã vào bill và mọi bill đã trả đủ, order chuyển sang `paid` (ghi `closed_at`) và bàn về `available`.

Các thao tác ghi đều chạy trong transaction và khoá dòng order/bill (`SELECT ... FOR UPDATE`), nên các request đồng thời (ví dụ hai lần thanh toán cùng lúc) không gây trả trùng hay xuất hoá đơn trùng.

## Service layer

| Service            | Hàm chính |
|--------------------|-----------|
| `TableService`     | `CreateTable`, `GetAllTable`, `UpdateTables`, `DeleteTable` |
| `MenuItemService`  | `CreateMenuItem`, `GetAllItems`, `GetMenuItemByID`, `UpdateItems`, `DeleteItems` |
| `OrderService`     | `CreateOrder`, `GetAllOrders`, `GetOrderById`, `DeleteOrder` |
| `OrderItemService` | `CreateOrderItem`, `GetOrderItemByID`, `GetOrderItemsByOrderID`, `UpdateOrderItemQuantity`, `DeleteOrderItem` |
| `BillService`      | `CreateBill`, `GetBillByID`, `GetBillsByOrderID`, `GetRemainingAmount`, `DeleteBill` |
| `PaymentService`   | `CreatePayment`, `GetPaymentByID`, `GetPaymentsByBillID` |

Service trả về các lỗi sentinel (khai báo trong package `service`) để controller dùng `errors.Is` và map sang HTTP status:

| Lỗi | HTTP gợi ý |
|-----|-----------|
| `ErrTableNotFound`, `ErrMenuItemNotFound`, `ErrOrderNotFound`, `ErrOrderItemNotFound`, `ErrBillNotFound`, `ErrPaymentNotFound` | 404 |
| `ErrInvalidQuantity`, `ErrInvalidAmount`, `ErrInvalidTableStatus`, `ErrDuplicateOrderItems`, `ErrItemNotInOrder` | 400 |
| `ErrTableHasOpenOrder`, `ErrOrderNotOpen`, `ErrOrderItemBilled`, `ErrOrderHasItems`, `ErrBillHasPayments`, `ErrBillAlreadyPaid`, `ErrOverpayment`, `ErrNoItemsToBill`, `ErrTableStatusManaged` | 409 |

## API

Prefix: `/api`. Body và response đều là JSON, tên trường dạng `snake_case`. ID là UUID; ID sai định dạng trả `400`.

Lỗi trả về dạng `{"error": "..."}` với status:
- `400`: dữ liệu không hợp lệ
- `404`: không tìm thấy
- `409`: xung đột nghiệp vụ hoặc ràng buộc DB (trùng `table_number`/`item_number`, xoá bản ghi còn được tham chiếu…)
- `500`: lỗi hệ thống

### Bàn

| Method | Path          | Body | Mô tả |
|--------|---------------|------|-------|
| POST   | `/tables`     | `{"table_number": 1}` | Tạo bàn |
| GET    | `/tables`     |  | Danh sách bàn (`/getalltables` vẫn giữ để tương thích) |
| PUT    | `/tables/:id` | `{"table_number": 1, "status": "reserved"}` | Cập nhật. `occupied` do hệ thống tự quản lý: bàn có order đang mở phải giữ `occupied`, bàn không có order không được đặt `occupied` (409) |
| DELETE | `/tables/:id` |  | Xoá bàn (409 nếu bàn đã có order) |

### Thực đơn

| Method | Path              | Body | Mô tả |
|--------|-------------------|------|-------|
| POST   | `/menu-items`     | `{"item_number": 1, "name": "Phở", "category": "main", "price": 50}` | Tạo món |
| GET    | `/menu-items`     |  | Danh sách món (sắp theo `item_number`) |
| GET    | `/menu-items/:id` |  | Chi tiết món |
| PATCH  | `/menu-items/:id` | bất kỳ trường nào ở trên | Cập nhật từng phần |
| DELETE | `/menu-items/:id` |  | Xoá món (409 nếu món đã được gọi) |

### Order & món trong order

| Method | Path                | Body | Mô tả |
|--------|---------------------|------|-------|
| POST   | `/orders`           | `{"table_id": "<uuid>"}` | Mở order, bàn chuyển `occupied` |
| GET    | `/orders`           |  | Danh sách order (mới nhất trước) |
| GET    | `/orders/:id`       |  | Chi tiết order kèm `order_items` |
| DELETE | `/orders/:id`       |  | Xoá order **chưa có món** (mở nhầm), trả bàn |
| POST   | `/orders/:id/items` | `{"menu_item_id": "<uuid>", "quantity": 2}` | Gọi món |
| GET    | `/orders/:id/items` |  | Các món của order |
| GET    | `/order-items/:id`  |  | Chi tiết một món |
| PATCH  | `/order-items/:id`  | `{"quantity": 3}` | Đổi số lượng (chưa vào bill) |
| DELETE | `/order-items/:id`  |  | Huỷ món (chưa vào bill) |

### Hoá đơn & thanh toán

| Method | Path                  | Body | Mô tả |
|--------|-----------------------|------|-------|
| POST   | `/orders/:id/bills`   | *(không bắt buộc)* `{"order_item_ids": ["<uuid>", ...]}` | Xuất bill. Không gửi body: gom tất cả món chưa vào bill; có danh sách: chia bill |
| GET    | `/orders/:id/bills`   |  | Các bill của order |
| GET    | `/bills/:id`          |  | Chi tiết bill kèm `order_items`, `payments`, `paid_amount`, `remaining_amount` |
| DELETE | `/bills/:id`          |  | Huỷ bill chưa có thanh toán |
| POST   | `/bills/:id/payments` | `{"amount": 40}` | Thanh toán (cho phép trả một phần) |
| GET    | `/bills/:id/payments` |  | Lịch sử thanh toán của bill |
| GET    | `/payments/:id`       |  | Chi tiết một lần thanh toán |

### Ví dụ một lượt phục vụ

```bash
API=localhost:8080/api
T=$(curl -s -X POST $API/tables -d '{"table_number":1}' | jq -r .id)
M=$(curl -s -X POST $API/menu-items -d '{"item_number":1,"name":"Phở","category":"main","price":50}' | jq -r .id)

O=$(curl -s -X POST $API/orders -d "{\"table_id\":\"$T\"}" | jq -r .id)
curl -s -X POST $API/orders/$O/items -d "{\"menu_item_id\":\"$M\",\"quantity\":2}"

B=$(curl -s -X POST $API/orders/$O/bills | jq -r .id)   # tổng 100
curl -s -X POST $API/bills/$B/payments -d '{"amount":100}'

curl -s $API/orders/$O | jq .status                      # "paid", bàn đã về "available"
```

## Hạn chế & việc cần làm

- [ ] Viết test tự động (unit test + integration test với Postgres)
- [ ] Tiền đang lưu bằng `float32`; nên chuyển sang số nguyên (đơn vị nhỏ nhất) hoặc decimal. Hiện so sánh tiền dùng sai số 0.005
- [ ] Phân trang / lọc cho các API danh sách (vd. `GET /orders?status=open`)
- [ ] Chưa có hoàn tiền / huỷ payment
- [ ] Chưa có xác thực / phân quyền
