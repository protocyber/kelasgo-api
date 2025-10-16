# Request ID Implementation

## Overview
This implementation adds a unique request ID to every HTTP request and logs it only when errors occur (status code >= 400).

## Components

### 1. Request ID Middleware (`request_id.go`)
- **Purpose**: Generates and tracks a unique ID for each request
- **Location**: `internal/server/middleware/request_id.go`

**Features:**
- Checks for existing `X-Request-ID` header in incoming requests
- Generates a new UUID if no request ID exists
- Stores request ID in Gin context for easy access throughout the request lifecycle
- Adds `X-Request-ID` to response headers
- Provides `GetRequestID()` helper function to retrieve the request ID from context

**Constants:**
- `RequestIDHeader`: Header key "X-Request-ID"
- `RequestIDKey`: Context key "request_id"

### 2. Logger Middleware Updates (`logger.go`)
Updated two logging functions to include request ID only when errors occur:

#### `Logger()` Function
- Logs request ID when `statusCode >= 400`
- Uses `param.Request.Header.Get(RequestIDHeader)` to retrieve the ID

#### `RequestLogger()` Function
- Logs request ID when `statusCode >= 400`
- Uses `GetRequestID(c)` helper function to retrieve the ID from context
- Works in both development and production environments

### 3. Routes Configuration (`routes.go`)
- Added `middleware.RequestID()` as the **first** middleware in the chain
- Order is important: RequestID must be set before Logger can access it

**Middleware Order:**
```go
r.Use(middleware.RequestID())        // 1. Generate/set request ID
r.Use(middleware.Logger())           // 2. Basic logger
r.Use(middleware.RequestLogger(cfg)) // 3. Enhanced logger
r.Use(middleware.AppContextMiddleware(cfg))
r.Use(middleware.CORSMiddleware(cfg.App.CORS))
```

## Usage Examples

### 1. Client Sends Request ID
```bash
curl -H "X-Request-ID: my-custom-id-123" http://localhost:8080/v1/health
```
Response will include the same ID in the header.

### 2. Auto-Generated Request ID
```bash
curl http://localhost:8080/v1/health
```
Response will include a new UUID in the `X-Request-ID` header.

### 3. Accessing Request ID in Handlers
```go
func MyHandler(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    // Use requestID in your business logic or error responses
}
```

## Log Output Examples

### Successful Request (200 OK)
```json
{
  "level": "info",
  "method": "GET",
  "uri": "/v1/health",
  "remote_ip": "127.0.0.1",
  "status": 200,
  "latency": 1234567,
  "message": "HTTP Request"
}
```
Note: No request_id in the log because the request was successful.

### Error Request (400+)
```json
{
  "level": "error",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "uri": "/v1/users",
  "remote_ip": "127.0.0.1",
  "status": 400,
  "latency": 2345678,
  "message": "HTTP Request"
}
```
Note: request_id is included because status code >= 400.

## Benefits

1. **Traceability**: Each request has a unique identifier for debugging
2. **Error Tracking**: Easy to correlate logs, errors, and user reports
3. **Conditional Logging**: Request ID only logged when needed (errors), reducing log noise
4. **Client Flexibility**: Clients can provide their own request IDs for end-to-end tracing
5. **Response Headers**: Request ID returned to client for their reference

## Dependencies

- `github.com/google/uuid` (already in go.mod)
- `github.com/gin-gonic/gin`
- `github.com/rs/zerolog`

## Context Logger Utility for Business Logic

### Overview
A reusable utility (`internal/util/context_logger.go`) has been created for consistent logging with automatic request ID inclusion in handlers and services.

### Features

1. **Automatic Request ID Inclusion**: Request ID is automatically added to error and warning logs
2. **Context Propagation**: Request ID can be passed from handlers to services
3. **Additional Context**: Support for tenant ID and user ID
4. **Convenience Methods**: Helper methods for common logging patterns
5. **Type-safe**: Handles different field types automatically

### Usage in Handlers

```go
package handler

import (
    "github.com/protocyber/kelasgo-api/internal/util"
    "github.com/rs/zerolog/log"
)

func (h *AuthHandler) Login(c *gin.Context) {
    // Create context logger - automatically captures request_id
    logger := util.NewContextLogger(c)
    
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Error log will include request_id automatically
        logger.Error().
            Err(err).
            Str("remote_ip", c.ClientIP()).
            Msg("Failed to bind login request")
        
        c.JSON(http.StatusBadRequest, dto.Response{
            Success: false,
            Message: "Invalid request",
        })
        return
    }
    
    // Pass request ID to service layer
    requestID := logger.GetRequestID()
    result, err := h.authService.Login(requestID, req)
    if err != nil {
        logger.Warn().
            Err(err).
            Str("email", req.Email).
            Msg("Login failed")
        // ... handle error
        return
    }
    
    // Success log - use standard log (no request_id)
    log.Info().
        Str("user_id", result.User.ID.String()).
        Msg("Login successful")
    
    c.JSON(http.StatusOK, dto.Response{Success: true, Data: result})
}
```

### Usage in Services

```go
package service

import (
    "github.com/protocyber/kelasgo-api/internal/util"
    "github.com/rs/zerolog/log"
)

func (s *studentService) Create(requestID string, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error) {
    // Create logger with request ID and tenant ID
    logger := util.NewContextLoggerWithRequestID(requestID).
        WithTenantID(tenantID.String())
    
    // Validation error
    if err := s.validateStudent(req); err != nil {
        // Warning log will include request_id and tenant_id
        logger.Warn().
            Err(err).
            Str("student_number", req.StudentNumber).
            Msg("Validation failed")
        return nil, err
    }
    
    // Database error
    if err := s.repo.Create(student); err != nil {
        // Error log will include request_id and tenant_id
        logger.Error().
            Err(err).
            Str("student_number", req.StudentNumber).
            Msg("Failed to create student")
        return nil, err
    }
    
    // Success - use standard log (no request_id)
    log.Info().
        Str("student_id", student.ID.String()).
        Msg("Student created successfully")
    
    return student, nil
}
```

### Convenience Methods

```go
// Simple error logging with multiple fields
logger.LogError(err, "Database query failed", map[string]interface{}{
    "table": "students",
    "operation": "insert",
    "duration_ms": 150,
})

// Warning with fields
logger.LogWarn("Validation failed", map[string]interface{}{
    "field": "email",
    "value": req.Email,
})

// Info logging
logger.LogInfo("Operation completed", map[string]interface{}{
    "records_processed": 100,
})
```

### Migration Strategy

#### Step 1: Update Handler Methods
1. Add `logger := util.NewContextLogger(c)` at the start of each handler
2. Replace error logs with context logger methods
3. Keep success logs using standard `log` package
4. Extract request ID for passing to services: `requestID := logger.GetRequestID()`

#### Step 2: Update Service Interfaces
Add `requestID string` as the first parameter to all service methods:

```go
// Before
type StudentService interface {
    Create(tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error)
}

// After
type StudentService interface {
    Create(requestID string, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error)
}
```

#### Step 3: Update Service Implementations
1. Accept `requestID` parameter in all methods
2. Create context logger: `logger := util.NewContextLoggerWithRequestID(requestID)`
3. Add tenant/user context if available: `.WithTenantID(tenantID.String())`
4. Use context logger for error/warning logs
5. Use standard log for success logs

### Best Practices

1. ✅ **Use ContextLogger for errors/warnings**: Ensures request ID is tracked
2. ✅ **Use standard log for success**: Keeps logs clean
3. ✅ **Pass request ID to services**: Enable tracing throughout the application
4. ✅ **Add context when available**: tenant_id, user_id provide valuable debugging info
5. ❌ **Don't log request ID on success**: Unnecessary noise in logs
6. ❌ **Don't create multiple loggers**: Create once per handler/service method

### Example Log Outputs

#### Error in Service (with request_id and tenant_id)
```json
{
  "level": "error",
  "request_id": "c783a782-425e-4194-853a-5aec377c2be8",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440000",
  "error": "tenant user not found",
  "tenant_user_id": "123e4567-e89b-12d3-a456-426614174000",
  "time": "2025-10-16T12:45:03+07:00",
  "message": "Tenant user not found during student creation"
}
```

#### Success Log (clean, no request_id)
```json
{
  "level": "info",
  "student_id": "789e4567-e89b-12d3-a456-426614174000",
  "student_number": "STD-2025-001",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440000",
  "time": "2025-10-16T12:45:03+07:00",
  "message": "Student created successfully"
}
```

## Additional Documentation

- **Detailed Usage Guide**: See `docs/CONTEXT_LOGGER_USAGE.md`
- **Handler Examples**: See `docs/examples/handler_with_context_logger.go`
- **Service Examples**: See `docs/examples/service_with_context_logger.go`

## Files

### New Files
- `internal/server/middleware/request_id.go` - Request ID middleware
- `internal/util/context_logger.go` - Context logger utility for business logic

### Modified Files
- `internal/server/middleware/logger.go` - Updated to log request ID on errors only
- `internal/server/routes.go` - Added RequestID middleware

### Documentation Files
- `docs/REQUEST_ID_IMPLEMENTATION.md` - This file
- `docs/CONTEXT_LOGGER_USAGE.md` - Detailed usage guide
- `docs/examples/handler_with_context_logger.go` - Handler examples
- `docs/examples/service_with_context_logger.go` - Service examples
