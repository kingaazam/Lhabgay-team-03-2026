# Lhabgay Backend Setup

## Required Tables

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(150) UNIQUE NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  title VARCHAR(200) NOT NULL,
  author VARCHAR(150) NOT NULL,
  description TEXT NOT NULL,
  category VARCHAR(100) NOT NULL,
  cover_image_path TEXT NOT NULL,
  book_file_path TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Run

```bash
go run backend/main.go
```

Open:

```text
http://localhost:8080/Login.html
```

The server prints:

```text
Database connected successfully
Server running at http://localhost:8080
```

## Admin Login

- Email: `admin@lhabgay.com`
- Password: `admin123`
- Redirects to `admin.html`

## API Endpoints

- `POST /auth/login`
- `POST /user/signup`
- `GET /logout`
- `POST /admin/books/upload`
- `GET /books`
- `GET /books/{id}`

## Test Flow

1. Start the server with `go run backend/main.go`.
2. Open `UserSignUp.html`, create a user, then log in from `Login.html`.
3. Log in as admin with `admin@lhabgay.com` and `admin123`.
4. Open `bookUpload.html`, fill the form, upload a cover image and a book file.
5. Open `home.html` or `search.html` to see uploaded books.
6. Click a book to open `details.html?id=BOOK_ID` and use the Read Now button to open the uploaded file.
