# Smart Desinger

Smart Desinger is a Go-based image storage and transformation service built around PostgreSQL and AWS S3.
It lets authenticated users upload images, fetch them by ID, and apply basic image transformations such as resizing, cropping, rotation, format conversion, and filters using FFmpeg.

## Features

- User registration and login with JWT-based authentication
- Secure password hashing using Argon2
- Image upload to S3 using multipart form requests
- Image retrieval using user-specific access checks
- Image transformation pipeline with FFmpeg
- PostgreSQL persistence for users and image metadata
- Lightweight browser UI served from the `static/` folder

## Tech Stack

- Go
- PostgreSQL
- AWS S3
- JWT
- FFmpeg / FFprobe
- SQLC-style query definitions in `sql/queries`

## Project Structure

```text
.
├── main.go
├── go.mod
├── sql/
│   ├── queries/
│   └── schema/
├── internal/
│   ├── api/
│   ├── app/
│   ├── auth/
│   ├── database/
│   └── responses/
├── static/
│   ├── index.html
│   └── index.js
└── assets/
```

## Prerequisites

Before running the project, make sure you have:

- Go 1.22+ installed
- PostgreSQL running and reachable
- An AWS S3 bucket configured
- AWS CLI credentials or shared profile configured for the app
- FFmpeg and FFprobe installed on the machine

## Environment Variables

Create a `.env` file in the project root with values like:

```env
DB_URL=postgres://username:password@localhost:5432/smart_desinger?sslmode=disable
TOKEN_SECRET=your-super-secret-key
BUCKET_NAME=your-s3-bucket-name
BUCKET_REGION=us-east-1
```

The project reads AWS credentials using the shared profile named `smart-designer`:

```bash
aws configure --profile smart-designer
```

This matches the code in `main.go`:

```go
config.WithSharedConfigProfile("smart-designer")
```

## Database Setup

1. Create a PostgreSQL database.
2. Run the schema files in `sql/schema/`:
   - `001_create_users.sql`
   - `002_create_images.sql`

Example:

```bash
psql -U postgres -d smart_desinger -f sql/schema/001_create_users.sql
psql -U postgres -d smart_desinger -f sql/schema/002_create_images.sql
```

## Running the App

Install dependencies and start the server:

```bash
go mod download
go run .
```

The service listens on:

```text
http://localhost:8080
```

The frontend is served from the `static/` directory and is available at:

```text
http://localhost:8080/
```

## Authentication

Authentication is done using a Bearer token in the `Authorization` header.

### Register

```http
POST /register
Content-Type: application/json
```

Request body:

```json
{
  "username": "alice",
  "password": "secret123"
}
```

Response includes a JWT in the `token` field.

### Login

```http
POST /login
Content-Type: application/json
```

Request body:

```json
{
  "username": "alice",
  "password": "secret123"
}
```

## Image APIs

All image endpoints require a valid Bearer token.

### Upload an image

```http
POST /images
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

Form field:

- `image`: image file (`png` or `jpeg` only)

Example with curl:

```bash
curl -X POST http://localhost:8080/images \
  -H "Authorization: Bearer <token>" \
  -F "image=@/path/to/image.png"
```

### Get an image

```http
GET /images/{id}
Authorization: Bearer <token>
```

### Transform an image

```http
POST /images/{id}/transform
Authorization: Bearer <token>
Content-Type: application/json
```

Request example:

```json
{
  "transformations": {
    "resize": {
      "width": 800,
      "height": 600
    },
    "crop": {
      "width": 400,
      "height": 400,
      "x": 50,
      "y": 50
    },
    "rotate": 90,
    "format": ".png",
    "filters": {
      "grayscale": true,
      "dither": false
    }
  }
}
```

Supported transform options:

- `resize`: width and height
- `crop`: width, height, x, y
- `rotate`: degrees
- `format`: `jpg`, `png`, or `webp`
- `filters.grayscale`: boolean
- `filters.dither`: boolean

## Notes

- The code currently uses `ffmpeg` and `ffprobe` directly from the system PATH.
- Uploaded images are stored with generated object keys in S3.
- Image metadata is stored in PostgreSQL, while actual image bytes are kept in S3.
- This project is intended as a backend service and demo app for image processing workflows.

## License

This project does not currently include a license file. If you plan to distribute it, add an appropriate license before publishing.
