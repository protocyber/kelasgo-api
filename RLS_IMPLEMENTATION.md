# Row Level Security (RLS) Implementation Guide

## 🎯 Overview

Row Level Security is implemented using PostgreSQL's RLS feature with tenant isolation for the Echo v4 web framework. The tenant context is set using the `app.current_tenant` session variable before each query.

## 🔧 Implementation Points

### 1. **Middleware Level (Entry Point)**

The tenant context is automatically set in middleware for every request:

```go
// In main.go or routes setup
func setupRoutes(e *echo.Echo, db *database.DatabaseConnections) {
    // Apply tenant middleware globally
    e.Use(middleware.TenantMiddleware(db))
    
    // For routes that require tenant
    tenantRoutes := e.Group("/api")
    tenantRoutes.Use(middleware.RequireTenant())
    
    // Your API routes
    tenantRoutes.GET("/users", userHandler.List)
    tenantRoutes.POST("/users", userHandler.Create)
}
```

### 2. **Repository Level**

All repositories should embed BaseRepository and set tenant context:

```go
type userRepository struct {
    BaseRepository
}

func (r *userRepository) Create(user *model.User) error {
    // Set tenant context before database operation
    if err := r.SetTenantContext(user.TenantID); err != nil {
        return err
    }
    return r.db.Write.Create(user).Error
}

func (r *userRepository) GetByUsernameAndTenant(username string, tenantID uuid.UUID) (*model.User, error) {
    // Set tenant context for this query
    if err := r.SetTenantContext(tenantID); err != nil {
        return nil, err
    }
    
    var user model.User
    err := r.db.Read.Preload("Role").Where("username = ? AND tenant_id = ?", username, tenantID).First(&user).Error
    return &user, err
}
```

### 3. **Service Level**

Services get tenant ID from context and pass it to repositories:

```go
func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
    var tenantID uuid.UUID
    if req.TenantID != "" {
        tenantID, err = uuid.Parse(req.TenantID)
        if err != nil {
            return nil, errors.New("invalid tenant ID format")
        }
    }

    // Repository automatically sets tenant context
    user, err := s.userRepo.GetByUsernameAndTenant(req.Username, tenantID)
    if err != nil {
        return nil, errors.New("invalid username or password")
    }
    
    // Continue with authentication...
}
```

### 4. **Handler Level**

Handlers extract tenant ID from middleware context:

```go
func (h *UserHandler) CreateUser(c echo.Context) error {
    // Get tenant ID from middleware context
    tenantID := middleware.GetTenantID(c)
    if tenantID == uuid.Nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Tenant ID required"})
    }

    var req dto.CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
    }

    // Service/repository will handle tenant context
    user, err := h.userService.CreateUser(tenantID, req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
    }

    return c.JSON(http.StatusCreated, user)
}
```

## 🚀 Usage Examples

### **1. Multi-tenant API Calls**

```bash
# Using Header
curl -X GET "https://api.yourdomain.com/api/users" \
  -H "X-Tenant-ID: 123e4567-e89b-12d3-a456-426614174000" \
  -H "Authorization: Bearer your-jwt-token"

# Using Query Parameter
curl -X GET "https://api.yourdomain.com/api/users?tenant_id=123e4567-e89b-12d3-a456-426614174000" \
  -H "Authorization: Bearer your-jwt-token"

# Using Subdomain (if implemented)
curl -X GET "https://tenant-slug.yourdomain.com/api/users" \
  -H "Authorization: Bearer your-jwt-token"
```

### **2. Login with Tenant**

```bash
curl -X POST "https://api.yourdomain.com/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "john.doe",
    "password": "secure-password"
  }'
```

## 🔒 Security Benefits

1. **Database-Level Isolation**: PostgreSQL RLS ensures data isolation at the database level
2. **Automatic Enforcement**: Once set, all queries automatically respect tenant boundaries
3. **Zero Trust**: Even if application logic fails, database prevents cross-tenant data access
4. **Performance**: RLS policies use indexes efficiently for tenant filtering

## ⚡ Performance Considerations

1. **Connection Pooling**: Tenant context is set per connection, so connection pooling works efficiently
2. **Indexing**: All tenant-enabled tables have indexes on `tenant_id` for fast filtering
3. **Policy Optimization**: RLS policies use the `current_tenant_id()` function for efficient tenant checking

## 🔧 Debugging

### Check Current Tenant Context
```sql
SELECT current_setting('app.current_tenant', true);
```

### Test RLS Policies
```sql
-- Set tenant context
SELECT set_config('app.current_tenant', '123e4567-e89b-12d3-a456-426614174000', false);

-- Query should only return data for this tenant
SELECT * FROM users;
```

### Verify Policy Application
```sql
-- Check if RLS is enabled
SELECT schemaname, tablename, rowsecurity 
FROM pg_tables 
WHERE rowsecurity = true;

-- Check policies
SELECT * FROM pg_policies WHERE tablename = 'users';
```

## 🚨 Important Notes

1. **Superuser Bypass**: PostgreSQL superusers bypass RLS by default
2. **Policy Testing**: Always test RLS policies with non-superuser roles
3. **Performance Impact**: RLS adds slight overhead but provides strong security
4. **Migration Safety**: RLS policies are applied after enabling, ensuring data safety during deployment

## 🔄 Migration Path

If you have existing data without tenant_id:

1. **Add tenant_id columns** (already done in your migrations)
2. **Populate tenant_id** for existing data
3. **Enable RLS** on tables (already done)
4. **Apply policies** (already done)
5. **Update application code** to use tenant-aware methods
6. **Test thoroughly** with different tenant contexts

This implementation ensures that your SaaS application has robust tenant isolation at both the application and database levels.
