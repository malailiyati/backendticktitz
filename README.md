# 🎬 Tickitz Backend

Backend untuk aplikasi **Tickitz – Movie Ticketing Web App**.  
Dibangun menggunakan **Go + Gin** dengan **PostgreSQL** sebagai database, serta support **Docker** untuk containerization.  
Menyediakan REST API untuk autentikasi, manajemen film, jadwal, kursi, hingga pemesanan tiket dan pembayaran.

---

## 🔧 Tech Stack

- [Go](https://go.dev/) – Backend language
- [Gin](https://gin-gonic.com/) – HTTP web framework
- [PostgreSQL](https://www.postgresql.org/) – Database
- [Docker](https://www.docker.com/) – Containerization
- [JWT](https://jwt.io/) – Authentication

---

## 🌱 Environment

Buat file `.env` di folder `backend` dengan isi seperti berikut:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=tickitz

# App
APP_PORT=8080
JWT_SECRET=yourjwtsecret
```

## Instalation

1. Clone Repo

```
git clone https://github.com/malailiyati/tickitz.git
cd tickitz-backend
```

2. Setup Environment

```
cp .env.example .env
```

3. Jalankan Secara Lokal

```
go run main.go
```

## License

https://dbdiagram.io/d/sistem-ticketing-68afbd78777b52b76ce3fc24

## License

- MIT License
- Copyright (c) 2025 Tickitz

# Contact Info

- Author: Ma'la Iliyati
- Email: malailiyati107@gmail.com

# Related Project

Ticketing-react

```

```
