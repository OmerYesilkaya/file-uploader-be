## Description

A self-hosted file server with:
Per-user authentication (JWT)
Chunked file uploads
Listing and downloading files
Zipping selected files
Private/public access toggle
Signed URL support for private file downloads
Optional image preview rendering

## ✅ Project Tasks

#### 1. 📦 Project Setup

-   ~go mod init github.com/yourname/fileuploader~
-   ~Setup folder structure (see below)~
-   ~Add config loader (env vars, ports, DB path)~
-   ~Set up SQLite DB and migration file~

#### 2. 🔐 Auth System (JWT)

-   ~User model (username, password hash)~
-   ~Signup/login endpoints~
-   ~Password hashing (bcrypt)~
-   ~JWT generation and verification~
-   ~Middleware to protect routes~
-   Add test user seed logic (optional)

#### 3. 📁 File Upload System

-   ~File upload and save to disk (intermediary step)~
-   Create upload directory per user (CONTINUE FROM HERE)
-   POST endpoint for chunked upload:
    -   Accept file metadata (filename, chunkIndex, totalChunks)
    -   Append chunks to temp file
    -   On final chunk, move to permanent user directory
-   Save file metadata in DB (name, path, size, type, visibility, userID)

#### 4. 📜 File Listing & Management

-   GET /files (list user files)
-   Toggle public/private (PATCH /files/:id/visibility)
-   Delete file (DELETE /files/:id)
-   Serve public files directly (/public/:filename)

#### 5. 📦 ZIP Download Support

-   POST /zip with list of file IDs
-   Zip them in memory or temp
-   Serve the zip file with:
    -   Optional signed URL
    -   Optional expiration

#### 6. 🔒 Signed URL Support

-   Create HMAC-signed download link
-   Add expiration to token
-   Validate signature before serving
-   Fallback to 403 if invalid/expired

#### 7. 🖼️ (Optional) Image Preview

-   Detect images (MIME type)
-   Generate low-res preview (optional lib)
-   Serve previews under /preview/:id

#### 8. 🧪 Basic Tests (Optional)

-   Unit test JWT auth
-   Test file upload logic
-   Test signed URL generation/validation

## Projecte Structure

```
/cmd/
  main.go
/internal/
  /api/
    /routes/
    /handlers/
    /middleware/
  /auth/
  /file/
  /zip/
  /signedurl/
  /db/
    db.go
    migrations/
  /config/
/pkg/

go.mod
README.md
```
