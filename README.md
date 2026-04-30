# Product Management App (React + Go + PostgreSQL)

Implementasi untuk **Question 1** pada dokumen tes full stack.

## Stack
- Frontend: React (Vite)
- Backend: Golang (`chi`, `database/sql`, `pgx`)
- Database: PostgreSQL

## Fitur yang tersedia
- Lihat daftar produk
- Tambah produk
- Ubah produk
- Hapus produk
- Search produk berdasarkan nama/SKU (`q`)
- Filter berdasarkan status (`active`/`inactive`)
- Validasi input API
- Migration SQL database
- Error response konsisten
- Pagination
- Login sederhana (JWT)
- Docker Compose
- Unit test (sample pada validasi payload)
- Layered structure (handler/service/repository)

## Struktur folder
- `backend/` API Golang
- `backend/migrations/` SQL schema/init data
- `frontend/` aplikasi React
- `docs/answers.md` jawaban Question 2 & Question 3

## Kredensial login default
- Username: `admin`
- Password: `password`

## Menjalankan dengan Docker Compose
Prasyarat: Docker Desktop aktif.

```bash
docker compose up --build
```

Akses:
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`
- Health check: `http://localhost:8080/health`

## Menjalankan manual (tanpa Docker)

### 1) Jalankan PostgreSQL
Buat database `product_management` lalu execute SQL di:
- `backend/migrations/001_init.sql`

### 2) Jalankan backend
```bash
cd backend
copy .env.example .env
go mod tidy
go run ./cmd/server
```

### 3) Jalankan frontend
```bash
cd frontend
npm install
npm run dev
```

## Endpoint utama
Base URL: `http://localhost:8080/api/v1`

1. Login
- `POST /auth/login`
- Body:
```json
{
  "username": "admin",
  "password": "password"
}
```

2. List produk
- `GET /products?q=&status=&page=1&limit=10`
- Header: `Authorization: Bearer <token>`

3. Tambah produk
- `POST /products`

4. Update produk
- `PUT /products/{id}`

5. Hapus produk
- `DELETE /products/{id}`

Body create/update:
```json
{
  "sku": "SKU-001",
  "name": "Keyboard Mechanical",
  "description": "Switch tactile",
  "price": 899000,
  "status": "active"
}
```

## Format response
Sukses:
```json
{
  "success": true,
  "data": {},
  "meta": {}
}
```

Error:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "input validation failed",
    "details": {
      "name": "name harus 2-120 karakter"
    }
  }
}
```

## Verifikasi yang sudah dijalankan
Backend:
```bash
go test ./...
go build ./cmd/server
```

Frontend:
```bash
npm run build
```
