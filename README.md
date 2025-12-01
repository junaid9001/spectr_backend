# Spectr Backend (Admin Dashboard + API)

Spectr is a Go backend (Gin + Gorm) for an e-commerce admin and public API with frontend templates. It provides user management, product management, cart, wishlist, orders, payments, OTP and email services, JWT auth, and admin views.

## Project structure
- Main: [`cmd/main.go`](cmd/main.go)
- Config: [`config/db.go`](config/db.go), [`config/migrate.go`](config/migrate.go)
- Routes: [`routes/register_routes.go`](routes/register_routes.go), [`routes/user_routes.go`](routes/user_routes.go), [`routes/admin_routes.go`](routes/admin_routes.go), [`routes/view_routes.go`](routes/view_routes.go)
- Controllers: [`controllers/auth_controllers.go`](controllers/auth_controllers.go), [`controllers/product_controllers.go`](controllers/product_controllers.go), [`controllers/cart_controllers.go`](controllers/cart_controllers.go), [`controllers/orders_controllers.go`](controllers/orders_controllers.go), [`controllers/whishlist_controllers.go`](controllers/whishlist_controllers.go), [`controllers/payment.go`](controllers/payment.go), [`controllers/user_controllers.go`](controllers/user_controllers.go), [`controllers/userManage_controller.go`](controllers/userManage_controller.go)
- Middlewares: [`middlewares/auth_middlewares.go`](middlewares/auth_middlewares.go)
- Models: [`models/users.go`](models/users.go), [`models/product.go`](models/product.go), [`models/cart_item.go`](models/cart_item.go), [`models/orders.go`](models/orders.go), [`models/order_item.go`](models/order_item.go), [`models/payment.go`](models/payment.go), [`models/whishlist.go`](models/whishlist.go), [`models/refresh_token.go`](models/refresh_token.go), [`models/otp.go`](models/otp.go), [`models/appStats.go`](models/appStats.go)
- Utilities: [`utils/generatetokens.go`](utils/generatetokens.go), [`utils/hash.go`](utils/hash.go), [`utils/validate_jwt.go`](utils/validate_jwt.go), [`utils/getUserId_helper.go`](utils/getUserId_helper.go), [`utils/params_conv.go`](utils/params_conv.go)
- Services: [`services/otp_service.go`](services/otp_service.go), [`services/mail_service.go`](services/mail_service.go)
- Templates: [`templates/login.html`](templates/login.html), [`templates/dashboard.html`](templates/dashboard.html), [`templates/users.html`](templates/users.html), [`templates/products.html`](templates/products.html), [`templates/orders.html`](templates/orders.html)
- Docker & compose: [`Dockerfile`](Dockerfile), [`compose.yaml`](compose.yaml), example env: [`.env.example`](.env.example)
- Uploads: `uploads/` (static file storage)

## Features
- User auth, signin/signup, email verification and OTP flow: see [`controllers/auth_controllers.go`](controllers/auth_controllers.go) and OTP generation/validation in [`services/otp_service.go`](services/otp_service.go).
- JWT-based access tokens and refresh tokens: see [`utils/generatetokens.go`](utils/generatetokens.go) and validation in [`utils/validate_jwt.go`](utils/validate_jwt.go).
- Admin and user route protection using [`middlewares/auth_middlewares.go`](middlewares/auth_middlewares.go).
- Product management (Create/Read/Update/Delete) in [`controllers/product_controllers.go`](controllers/product_controllers.go).
- Cart and wishlist management (`controllers/cart_controllers.go`, [`controllers/whishlist_controllers.go`](controllers/whishlist_controllers.go)).
- Order placement, detail, cancellation, restock, admin order listing & status updates (`controllers/orders_controllers.go`).
- Payments & app statistics updates in [`controllers/payment.go`](controllers/payment.go) and [`models/appStats.go`](models/appStats.go).
- Email sending via [`services/mail_service.go`](services/mail_service.go).
- Database models and relationships defined in `models/*`.
- A server-rendered admin UI using `templates/*.html` for basic admin interactions.

## Main entry points & key symbols
- Start-up: `main()` in [`cmd/main.go`](cmd/main.go) runs [`config.LoadEnv`](config/db.go), [`config.ConnectDB`](config/db.go) and [`config.MigrateAll`](config/migrate.go), registers routes (`routes.RegisterRoutes` in [`routes/register_routes.go`](routes/register_routes.go)) and serves templates and static uploads.
- Authentication helpers: token generation (`utils.GenerateAccessToken`, `utils.GenerateRefreshToken`, `utils.SaveRefreshToken`, `utils.ValidateRefreshToken` in [`utils/generatetokens.go`](utils/generatetokens.go)); hashing (`utils.HashPassword`, `utils.CompareHashAndPass` in [`utils/hash.go`](utils/hash.go)); getting user from context (`utils.GetUserId` in [`utils/getUserId_helper.go`](utils/getUserId_helper.go)).
- Controllers handle HTTP logic: e.g. [`controllers/Login`](controllers/auth_controllers.go), [`controllers.PlaceOrder`](controllers/orders_controllers.go), [`controllers.CreatePayment`](controllers/payment.go).

## API routes
- Public:
  - GET /products — [`controllers.GetAllProducts`](controllers/product_controllers.go)
  - GET /product/:id — [`controllers.GetProductByID`](controllers/product_controllers.go)
- Auth:
  - POST /auth/signup — [`controllers.SignupHandler`](controllers/auth_controllers.go)
  - POST /auth/login — [`controllers.Login`](controllers/auth_controllers.go)
  - POST /auth/refresh — [`controllers.RefreshTokenHandler`](controllers/auth_controllers.go)
  - POST /auth/forgot — [`controllers.ForgotPassword`](controllers/auth_controllers.go)
  - POST /auth/reset — [`controllers.ResetPassword`](controllers/auth_controllers.go)
- User (requires JWT via `UserAuthMiddleware`):
  - GET /user/profile, PUT /user/profile — [`controllers.GetUserProfile`, `UpdateUserProfile`](controllers/user_controllers.go)
  - Cart: POST /user/cart, GET /user/cart — [`controllers.AddProductToCart`, `GetUserCart`](controllers/cart_controllers.go)
  - Order: POST /user/order, GET /user/orders, GET /user/order/:id, DELETE /user/order/:id — [`controllers.PlaceOrder`, `GetOrderHistory`, `GetDetailsOfOrder`, `DeleteOrderById`](controllers/orders_controllers.go)
  - Payments: POST /user/order/:id/payments, POST /user/payment/:payment_id/confirm — [`controllers.CreatePayment`, `ConfirmPayment`](controllers/payment.go)
  - Wishlist: POST /user/wishlist, GET /user/wishlist — [`controllers.AddToWishlist`, `GetWishList`](controllers/whishlist_controllers.go)
- Admin (requires `AdminAuthMiddleware`):
  - GET /admin/users, PUT /admin/users/:id/role, PUT /admin/users/:id/status — [`controllers.AllUsers`, `UpdateUserRole`, `UpdateUserStatus`](controllers/userManage_controller.go)
  - Product: /admin/product (create), /admin/product/:id (update/delete) — [`controllers.CreateProduct`, `UpdateProductByID`, `DeleteProductByID`](controllers/product_controllers.go)
  - Orders: GET /admin/orders, PATCH /admin/order/:id — [`controllers.GetAllOrders`, `UpdateOrderStatus`](controllers/orders_controllers.go)

## Frontend and templates
- Admin UI pages are served under `/view/*` and use templates: [`templates/login.html`](templates/login.html), [`templates/dashboard.html`](templates/dashboard.html), [`templates/users.html`](templates/users.html), [`templates/products.html`](templates/products.html), [`templates/orders.html`](templates/orders.html). The view routes are defined in [`routes/view_routes.go`](routes/view_routes.go).

## Environment variables
Use [`.env.example`](.env.example) as reference. Important vars:
- DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD — used in [`config/db.go`](config/db.go)
- JWT_SECRETKEY — used in [`utils/generatetokens.go`](utils/generatetokens.go)/[`utils/validate_jwt.go`](utils/validate_jwt.go)
- EMAIL, EMAIL_PASS — used in [`services/mail_service.go`](services/mail_service.go)

## Database
- Gorm models in `models/` and migrations done by [`config.MigrateAll`](config/migrate.go).

## Running locally
1. Set up env variables (or copy `.env.example` to `.env`).
2. Build and run:
   ```sh
   go run ./cmd
   ```
3. Open http://localhost:8080 and log in via admin view: `/view/login`.

## Docker & Docker Compose
- Dockerfile: [`Dockerfile`](Dockerfile)
- Compose: [`compose.yaml`](compose.yaml)

Build and run with Docker:
```sh
docker build -t spectr_backend .
docker run --env-file .env --network host spectr_backend
```

Compose:
```sh
docker compose up --build
```

Key steps when using Docker/Compose:
- Ensure `.env` is present and DB credentials map to your Postgres container/environment.
- Migrations run on start: [`config.MigrateAll`](config/migrate.go)

## Notes
- Uploads are stored in `uploads/` and served as static files in `cmd/main.go`.
- OTPs are emailed using `services.SendEmail` (`services/mail_service.go`) and OTP logic in [`services/otp_service.go`](services/otp_service.go).
- Tokens: `GenerateAccessToken` (45-min expiration) and refresh token stored with hashed token in DB (`models/refresh_token.go`) via `utils.SaveRefreshToken` and validated by `utils.ValidateRefreshToken`.

## Contributing
- Use the existing `controllers/*`, `models/*`, `utils/*` and `routes/*` patterns for new features.
- Follow conventions for migrations in [`config/migrate.go`](config/migrate.go).

## References (quick links)
- [`cmd/main.go`](cmd/main.go)
- [`config/db.go`](config/db.go)
- [`config/migrate.go`](config/migrate.go)
- [`routes/register_routes.go`](routes/register_routes.go)
- [`routes/user_routes.go`](routes/user_routes.go)
- [`routes/admin_routes.go`](routes/admin_routes.go)
- [`routes/view_routes.go`](routes/view_routes.go)
- [`controllers/auth_controllers.go`](controllers/auth_controllers.go)
- [`controllers/product_controllers.go`](controllers/product_controllers.go)
- [`controllers/cart_controllers.go`](controllers/cart_controllers.go)
- [`controllers/orders_controllers.go`](controllers/orders_controllers.go)
- [`controllers/whishlist_controllers.go`](controllers/whishlist_controllers.go)
- [`controllers/payment.go`](controllers/payment.go)
- [`controllers/user_controllers.go`](controllers/user_controllers.go)
- [`controllers/userManage_controller.go`](controllers/userManage_controller.go)
- [`middlewares/auth_middlewares.go`](middlewares/auth_middlewares.go)
- Models: [`models/users.go`](models/users.go), [`models/product.go`](models/product.go), [`models/cart_item.go`](models/cart_item.go), [`models/orders.go`](models/orders.go), [`models/order_item.go`](models/order_item.go), [`models/payment.go`](models/payment.go), [`models/whishlist.go`](models/whishlist.go), [`models/refresh_token.go`](models/refresh_token.go), [`models/otp.go`](models/otp.go), [`models/appStats.go`](models/appStats.go)
- Utils: [`utils/generatetokens.go`](utils/generatetokens.go), [`utils/hash.go`](utils/hash.go), [`utils/validate_jwt.go`](utils/validate_jwt.go), [`utils/getUserId_helper.go`](utils/getUserId_helper.go), [`utils/params_conv.go`](utils/params_conv.go)
- Services: [`services/otp_service.go`](services/otp_service.go), [`services/mail_service.go`](services/mail_service.go)
- Templates: [`templates/login.html`](templates/login.html), [`templates/dashboard.html`](templates/dashboard.html), [`templates/users.html`](templates/users.html), [`templates/products.html`](templates/products.html), [`templates/orders.html`](templates/orders.html)
- Docker & compose: [`Dockerfile`](Dockerfile), [`compose.yaml`](compose.yaml), [`.env.example`](.env.example)
