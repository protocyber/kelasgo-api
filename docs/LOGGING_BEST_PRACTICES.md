# Logging Best Practices

## Summary of Changes

This document outlines the logging best practices implemented in the KelasGo API and the refactoring performed to eliminate duplicate logging.

## Problem Identified

Previously, the same error was being logged at **three different layers**:
1. **Handler Layer** - Logged service failures
2. **Service Layer** - Logged business logic errors  
3. **Repository Layer** - Logged database errors

This resulted in **triple logging** for every error, making logs noisy and difficult to debug.

## Solution: Layer-Based Logging Strategy

### **Handler Layer** (Minimal Logging)
**Purpose**: Log only HTTP-specific concerns

**What TO Log:**
- ✅ Request binding/parsing errors (`log.Error`)
- ✅ Request validation failures (`log.Warn`)
- ✅ Missing authentication/authorization (`log.Error`)
- ✅ Missing tenant context (`log.Error`)

**What NOT to Log:**
- ❌ Service errors (already logged in service layer)
- ❌ Business logic failures
- ❌ Database errors

**Example:**
```go
// ✅ GOOD - Log HTTP-specific error
if err := c.ShouldBindJSON(&req); err != nil {
    log.Error().
        Err(err).
        Msg("Failed to bind request JSON")
    c.JSON(http.StatusBadRequest, dto.Response{...})
    return
}

// ✅ GOOD - Log validation failure
if err := h.validator.Struct(req); err != nil {
    log.Warn().
        Err(err).
        Interface("params", params).
        Msg("Request validation failed")
    c.JSON(http.StatusBadRequest, dto.Response{...})
    return
}

// ✅ GOOD - Don't log service errors (already logged)
student, err := h.studentService.Create(tenantID, req)
if err != nil {
    c.JSON(http.StatusBadRequest, dto.Response{
        Success: false,
        Message: "Failed to create student",
        Error:   err.Error(),
    })
    return
}
```

### **Service Layer** (Primary Logging)
**Purpose**: Main diagnostic layer - log all business logic decisions and errors

**What TO Log:**
- ✅ Business rule violations (`log.Warn`)
- ✅ Database operation failures (`log.Error`)
- ✅ External service failures (`log.Error`)
- ✅ Validation errors with full context (`log.Warn`)
- ✅ Unexpected errors with stack traces (`log.Error`)

**What NOT to Log:**
- ❌ Successful operations (avoid noise)

**Example:**
```go
// ✅ GOOD - Log business rule violation
if existingStudent != nil {
    log.Warn().
        Str("student_number", req.StudentNumber).
        Str("tenant_id", tenantID.String()).
        Msg("Student number already exists")
    return nil, errors.New("student number already exists")
}

// ✅ GOOD - Log database failure with context
err = s.studentRepo.Create(student)
if err != nil {
    log.Error().
        Err(err).
        Str("student_number", req.StudentNumber).
        Str("tenant_id", tenantID.String()).
        Msg("Failed to create student")
    return nil, errors.New("failed to create student")
}
```

### **Repository Layer** (Minimal Logging)
**Purpose**: Log only database-specific critical errors

**What TO Log:**
- ✅ Database connection errors (`log.Error`)
- ✅ Constraint violations (`log.Error`)
- ✅ Transaction failures (`log.Error`)
- ✅ Unexpected database errors (`log.Error`)

**What NOT to Log:**
- ❌ "Record not found" - Return error silently (let service decide if it's critical)
- ❌ Successful operations
- ❌ Business logic errors

**Example:**
```go
// ✅ GOOD - Return "not found" silently
func (r *studentRepository) GetByID(id uuid.UUID) (*model.Student, error) {
    var student model.Student
    err := r.db.Read.First(&student, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("student not found")
        }
        log.Error().
            Err(err).
            Str("operation", "get_student_by_id").
            Msg("Database query failed")
        return nil, err
    }
    return &student, nil
}

// ✅ GOOD - Log critical DB errors with operation name
func (r *studentRepository) Create(student *model.Student) error {
    if err := r.SetTenantContext(student.TenantID); err != nil {
        return err  // Don't log, let service handle it
    }
    err := r.db.Write.Create(student).Error
    if err != nil {
        log.Error().
            Err(err).
            Str("operation", "create_student").
            Msg("Database write operation failed")
    }
    return err
}

// ❌ BAD - Too verbose, duplicate info in error
func (r *studentRepository) Create(student *model.Student) error {
    if err := r.SetTenantContext(student.TenantID); err != nil {
        log.Error().
            Err(err).
            Str("student_number", student.StudentNumber).
            Str("tenant_id", student.TenantID.String()).
            Msg("Failed to set tenant context")
        return err
    }
    // ...
}
```

## Log Levels

| Level | Usage |
|-------|-------|
| `log.Error()` | System errors, database failures, unexpected conditions |
| `log.Warn()` | Business rule violations, validation failures, expected but notable events |
| `log.Info()` | Application lifecycle events (startup, shutdown, migrations) |
| `log.Debug()` | Detailed debugging information (disabled in production) |

## Structured Logging

Always use structured fields instead of string formatting:

```go
// ✅ GOOD - Structured logging
log.Error().
    Err(err).
    Str("student_id", id.String()).
    Str("tenant_id", tenantID.String()).
    Interface("params", params).
    Msg("Failed to create student")

// ❌ BAD - String formatting
log.Error().Msgf("Failed to create student %s in tenant %s: %v", 
    id.String(), tenantID.String(), err)
```

## Benefits of This Approach

1. **Reduced Log Noise**: Each error logged only once at the appropriate layer
2. **Better Debugging**: Clear separation of concerns makes it easier to identify the source
3. **Improved Performance**: Fewer log operations
4. **Cleaner Code**: Handlers focus on HTTP concerns, not logging implementation details
5. **Consistent Patterns**: Developers know exactly where to look for specific types of errors

## Files Refactored

### Handlers (Removed duplicate service error logging)
- `internal/domain/handler/student_handler.go`
- `internal/domain/handler/user_handler.go`
- `internal/domain/handler/auth_handler.go`

### Repositories (Removed debug "not found" logs)
- `internal/domain/repository/student_repository.go`
- `internal/domain/repository/user_repository.go`
- `internal/domain/repository/role_repository.go`

### Services (Kept as primary logging layer)
- `internal/domain/service/student_service.go`
- `internal/domain/service/user_service.go`
- `internal/domain/service/auth_service.go`

## Quick Reference

```
┌─────────────┬──────────────┬─────────────────────────────────────┐
│ Layer       │ Log Level    │ What to Log                         │
├─────────────┼──────────────┼─────────────────────────────────────┤
│ Handler     │ Error/Warn   │ HTTP errors, validation, auth       │
│ Service     │ Error/Warn   │ Business logic, main diagnostics    │
│ Repository  │ Error only   │ Critical DB errors                  │
└─────────────┴──────────────┴─────────────────────────────────────┘
```

## Migration Notes

When adding new endpoints or services:

1. **Handlers**: Only log request-specific errors
2. **Services**: Log all business logic and error paths with context
3. **Repositories**: Log only unexpected database errors
4. Use `Interface("params", params)` instead of logging individual fields
5. Return errors up the stack; let the service layer decide what to log
