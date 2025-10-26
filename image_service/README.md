# Image Service with Authentication

Image service yang menyediakan upload dan serving gambar dengan autentikasi menggunakan gRPC.

## Features

1. **gRPC Upload Image**: Upload gambar via gRPC dengan base64 encoding
2. **HTTP Image Serving**: Serving gambar via HTTP dengan autentikasi JWT
3. **Auth Integration**: Validasi token menggunakan auth service via gRPC

## Ports

- **gRPC Server**: Port `50053` (untuk upload image)
- **HTTP Server**: Port `8081` (untuk serving image)

## Usage

### 1. Start the Service

```bash
go run main.go
```

Output:
```
🚀 gRPC Server starting on port 50053
🌐 HTTP Server starting on port 8081
📸 Image serving endpoint: http://localhost:8081/storage/images/...
```

### 2. Upload Image (via gRPC)

Gunakan gRPC client untuk upload:

```go
// Example gRPC upload
imageData := base64.StdEncoding.EncodeToString(imageBytes)
request := &pb.UploadImageRequest{
    ImageData: imageData,
    ImageName: "test.jpg",
}

response, err := client.UploadImage(ctx, request)
// Response: {"image_url": "http://localhost:8081/storage/images/2025/10/26/test_1761478232.jpg"}
```

### 3. Access Image (via HTTP)

#### Dengan Authorization Header:
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8081/storage/images/2025/10/26/test_1761478232.jpg
```

#### Dengan Query Parameter:
```bash
curl "http://localhost:8081/storage/images/2025/10/26/test_1761478232.jpg?token=YOUR_JWT_TOKEN"
```

#### Contoh dari Browser/Frontend:
```html
<!-- Dengan token di query parameter -->
<img src="http://localhost:8081/storage/images/2025/10/26/test_1761478232.jpg?token=YOUR_JWT_TOKEN" />
```

```javascript
// Dengan Authorization header
fetch('http://localhost:8081/storage/images/2025/10/26/test_1761478232.jpg', {
    headers: {
        'Authorization': 'Bearer YOUR_JWT_TOKEN'
    }
})
.then(response => response.blob())
.then(blob => {
    const imageUrl = URL.createObjectURL(blob);
    document.getElementById('myImage').src = imageUrl;
});
```

## Environment Variables

```bash
# .env file
GRPC_PORT=50053                    # gRPC server port
HTTP_PORT=8081                     # HTTP server port
AUTH_SERVICE_URL=localhost:50052   # Auth service gRPC endpoint
```

## Authentication Flow

1. Client mengupload gambar via gRPC → mendapat URL
2. Client mengakses URL dengan JWT token
3. Image service memvalidasi token ke auth service
4. Jika valid, gambar dikembalikan
5. Jika tidak valid, return 401 Unauthorized

## File Structure

```
storage/
└── images/
    └── 2025/
        └── 10/
            └── 26/
                └── testupload_1761478232.jpg
```

## Development Mode

Untuk development, jika auth service tidak tersedia, service akan bypass validasi token dan tetap serving gambar dengan log warning.

## CORS Support

HTTP server sudah dikonfigurasi dengan CORS untuk mendukung akses dari browser:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, OPTIONS`
- `Access-Control-Allow-Headers: Authorization, Content-Type`