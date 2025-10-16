# Context Logger Implementation

This document explains how to use the Context Logger utility for consistent logging with request ID tracking across handlers and services.

## Overview

The `ContextLogger` utility automatically includes request ID (and optionally tenant ID and user ID) in all log entries when errors occur. This makes it easy to trace requests throughout the system.

## Features

- **Automatic Request ID inclusion**: Request ID is automatically added to error logs
- **Context propagation**: Request context can be passed from handlers to services
- **Flexible fields**: Support for tenant ID, user ID, and custom fields
- **Type-safe logging**: Helper methods for common field types
- **Consistent format**: Ensures all logs follow the same structure

## Usage

### In Handlers

```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/protocyber/kelasgo-api/internal/util"
)

func (h *AuthHandler) Login(c *gin.Context) {
    // Create context logger
    logger := util.NewContextLogger(c)
    
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Error log will automatically include request_id
        logger.Error().
            Err(err).
            Str("remote_ip", c.ClientIP()).
            Str("user_agent", c.Request.UserAgent()).
            Msg("Failed to bind login request JSON")
            
        c.JSON(http.StatusBadRequest, dto.Response{
            Success: false,
            Message: "Invalid request body",
            Error:   err.Error(),
        })
        return
    }
    
    // Pass request ID to service layer
    requestID := logger.GetRequestID()
    response, err := h.authService.Login(requestID, req)
    if err != nil {
        logger.Warn().
            Err(err).
            Str("email", req.Email).
            Str("remote_ip", c.ClientIP()).
            Msg("Login attempt failed")
            
        c.JSON(http.StatusUnauthorized, dto.Response{
            Success: false,
            Message: "Login failed",
            Error:   err.Error(),
        })
        return
    }
    
    // Success log (without request_id)
    log.Info().
        Str("user_id", response.User.ID.String()).
        Str("email", response.User.Email).
        Msg("User logged in successfully")
    
    c.JSON(http.StatusOK, dto.Response{
        Success: true,
        Message: "Login successful",
        Data:    response,
    })
}
```

### In Services

```go
package service

import (
    "github.com/protocyber/kelasgo-api/internal/util"
    "github.com/rs/zerolog/log"
)

func (s *studentService) Create(requestID string, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error) {
    // Create logger with request ID
    logger := util.NewContextLoggerWithRequestID(requestID).
        WithTenantID(tenantID.String())
    
    // Check if tenant user exists
    tenantUser, err := s.tenantUserRepo.GetByID(req.TenantUserID)
    if err != nil {
        // Error log will include request_id and tenant_id
        logger.Error().
            Err(err).
            Str("tenant_user_id", req.TenantUserID.String()).
            Msg("Tenant user not found during student creation")
        return nil, errors.New("tenant user not found")
    }
    
    // Validation warning
    if tenantUser.TenantID != tenantID {
        logger.Warn().
            Str("tenant_user_id", req.TenantUserID.String()).
            Str("expected_tenant", tenantID.String()).
            Str("actual_tenant", tenantUser.TenantID.String()).
            Msg("Tenant user does not belong to the specified tenant")
        return nil, errors.New("tenant user does not belong to this tenant")
    }
    
    // Create student
    student := &model.Student{
        TenantID:      tenantID,
        TenantUserID:  req.TenantUserID,
        StudentNumber: req.StudentNumber,
    }
    
    err = s.studentRepo.Create(student)
    if err != nil {
        logger.Error().
            Err(err).
            Str("student_number", req.StudentNumber).
            Msg("Failed to create student in database")
        return nil, errors.New("failed to create student")
    }
    
    // Success log (without request_id - use standard log)
    log.Info().
        Str("student_id", student.ID.String()).
        Str("student_number", student.StudentNumber).
        Str("tenant_id", tenantID.String()).
        Msg("Student created successfully")
    
    return student, nil
}
```

### Using Convenience Methods

The logger also provides convenience methods for simpler logging:

```go
// Log error with fields
logger.LogError(err, "Failed to create student", map[string]interface{}{
    "student_number": req.StudentNumber,
    "tenant_id": tenantID.String(),
    "remote_ip": c.ClientIP(),
})

// Log warning with fields
logger.LogWarn("Validation failed", map[string]interface{}{
    "field": "email",
    "value": req.Email,
})

// Log info with fields
logger.LogInfo("Operation completed", map[string]interface{}{
    "operation": "create_student",
    "duration_ms": 150,
})
```

## Service Interface Updates

Update service interfaces to accept `requestID` as the first parameter:

```go
type StudentService interface {
    Create(requestID string, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error)
    Update(requestID string, tenantID uuid.UUID, id uuid.UUID, req dto.UpdateStudentRequest) (*model.Student, error)
    Delete(requestID string, tenantID uuid.UUID, id uuid.UUID) error
    GetByID(requestID string, tenantID uuid.UUID, id uuid.UUID) (*model.Student, error)
    // ... other methods
}
```

## Best Practices

1. **Handlers**: Always create a `ContextLogger` from the gin context
2. **Services**: Accept `requestID` as first parameter and create logger at the start
3. **Error Logs**: Use context logger for errors (includes request_id automatically)
4. **Success Logs**: Use standard `log` package for success (no request_id needed)
5. **Warnings**: Use context logger for warnings about validation/business logic issues
6. **Propagation**: Pass request ID through all service method calls

## Migration Strategy

1. Update handler methods to create context logger
2. Extract request ID and pass to service methods
3. Update service interfaces to accept requestID parameter
4. Update service implementations to use context logger for errors
5. Test thoroughly to ensure request IDs are properly logged

## Example Log Output

### Error Log (with request_id)
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

### Success Log (without request_id)
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
