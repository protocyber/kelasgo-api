# Common Use Cases for `app.url` Configuration

The `app.url` configuration represents the **external URL** that clients use to access your API. This is different from `server.host` which controls where the server binds internally.

## Real-World Scenarios:

### Scenario 1: Development vs Production
```yaml
# Development config.yaml
app:
  url: 'http://10.201.0.25:8080'  # Your local network IP
server:
  host: '0.0.0.0'                 # Bind to all interfaces
  port: '8080'

# Production config.yaml  
app:
  url: 'https://api.kelasgo.com'   # Public domain
server:
  host: '0.0.0.0'                 # Bind to all interfaces  
  port: '8080'                    # Behind load balancer/reverse proxy
```

### Scenario 2: Behind Reverse Proxy/Load Balancer
```yaml
app:
  url: 'https://api.kelasgo.com'   # External URL clients use
server:
  host: '127.0.0.1'               # Only bind locally
  port: '3000'                    # Nginx forwards to this port
```

## Common Use Cases in Code:

### 1. **API Response Links (HATEOAS)**
```go
// GET /v1/students/123 response
{
  "success": true,
  "data": {
    "id": 123,
    "name": "John Doe",
    "links": {
      "self": "http://10.201.0.25:8080/v1/students/123",
      "edit": "http://10.201.0.25:8080/v1/students/123", 
      "enrollments": "http://10.201.0.25:8080/v1/students/123/enrollments",
      "grades": "http://10.201.0.25:8080/v1/students/123/grades"
    }
  }
}
```

### 2. **Email Templates**
```go
// Password reset email
resetURL := fmt.Sprintf("%s/v1/auth/reset-password?token=%s", appURL, token)

emailBody := fmt.Sprintf(`
  Hi %s,
  
  Click here to reset your password: %s
  
  This link expires in 1 hour.
`, user.Name, resetURL)
```

### 3. **Webhook Registration**  
```go
// Register webhook with payment gateway
webhookURL := fmt.Sprintf("%s/v1/webhooks/payment-callback", appURL)
paymentGateway.RegisterWebhook(webhookURL)
```

### 4. **File Download Links**
```go  
// File upload response
{
  "success": true,
  "data": {
    "file_id": "abc123",
    "download_url": "http://10.201.0.25:8080/v1/files/abc123/download",
    "preview_url": "http://10.201.0.25:8080/v1/files/abc123/preview"
  }
}
```

### 5. **Pagination Links**
```go
// List API with pagination
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 2,
    "limit": 10,
    "total": 100
  },
  "links": {
    "first": "http://10.201.0.25:8080/v1/students?page=1",
    "prev": "http://10.201.0.25:8080/v1/students?page=1", 
    "self": "http://10.201.0.25:8080/v1/students?page=2",
    "next": "http://10.201.0.25:8080/v1/students?page=3",
    "last": "http://10.201.0.25:8080/v1/students?page=10"
  }
}
```

### 6. **QR Codes & Deep Links**
```go
// Generate QR code for student profile
profileURL := fmt.Sprintf("%s/v1/students/%d/profile", appURL, studentID)
qrCode := generateQRCode(profileURL)
```

### 7. **API Documentation**
```go  
// Swagger/OpenAPI base URL
{
  "openapi": "3.0.0",
  "servers": [
    {
      "url": "http://10.201.0.25:8080/v1",
      "description": "Development server"
    }
  ]
}
```

## Implementation in Your Handlers:

```go
func (h *StudentHandler) GetStudent(c *gin.Context) {
    studentID := c.Param("id")
    appURL, _ := c.Get("app_url")
    
    student, err := h.studentService.GetByID(studentID)
    if err != nil {
        c.JSON(404, gin.H{"error": "Student not found"})
        return
    }
    
    response := gin.H{
        "success": true,
        "data": gin.H{
            "student": student,
            "links": gin.H{
                "self": fmt.Sprintf("%s/v1/students/%s", appURL, studentID),
                "enrollments": fmt.Sprintf("%s/v1/students/%s/enrollments", appURL, studentID),
                "attendance": fmt.Sprintf("%s/v1/students/%s/attendance", appURL, studentID),
                "grades": fmt.Sprintf("%s/v1/students/%s/grades", appURL, studentID),
                "fees": fmt.Sprintf("%s/v1/students/%s/fees", appURL, studentID),
            },
        },
    }
    
    c.JSON(200, response)
}
```

## Benefits:

1. **Environment Flexibility**: Same code works in dev/staging/production with different URLs
2. **Client Convenience**: Clients get ready-to-use URLs in responses  
3. **SEO/Discoverability**: Proper canonical URLs for your API
4. **Integration**: External services can callback to correct URLs
5. **Documentation**: API docs show correct base URLs

## Your Current Setup:

- `server.host: '0.0.0.0'` ✅ Server binds to all interfaces
- `server.port: '8080'` ✅ Server listens on port 8080  
- `app.url: 'http://10.201.0.25:8080'` ✅ External URL for API responses

This means:
- Your server accepts connections from anywhere (0.0.0.0:8080)
- API responses include links using your network IP (10.201.0.25:8080)
- Your local PC can access the API and get proper links in responses
