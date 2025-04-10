# Developer Guide: Integrating GORM + SQLite + Casbin (RBAC) in a Gin Project

This guide walks you through the main steps of setting up a Gin web application with GORM for data persistence using SQLite and Casbin for role-based access control (RBAC).

## 1. Prerequisites

- Go installed (version 1.23.4 or higher)
- A Mac machine with a terminal
- Basic understanding of Gin, GORM, and Casbin

## 2. Project Structure

```
my-casbin/
├─ cmd/api/
│  ├─ rbac_model.conf    # Casbin RBAC model configuration
│  └─ policy.csv         # Casbin policy file
├─ internal/
│  ├─ database/
│  │  └─ database.go     # Database/GORM & Casbin initialization
│  └─ handler/           # HTTP request handlers
│     └─ handler.go
├─ go.mod
├─ Makefile
└─ README.md
```

## 3. Configuring GORM with SQLite

The `SetupCasbin()` function in `/internal/database/database.go` handles:

- SQLite database connection using GORM
- Casbin adapter configuration
- RBAC model loading

```go
func SetupCasbin() (*casbin.Enforcer, *gorm.DB, error) {
    // Connect to SQLite DB using GORM
    db, err := gorm.Open(sqlite.Open(dburl), &gorm.Config{})
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect database: %w", err)
    }

    // Initialize Casbin with xorm adapter
    adapter, err := xormadapter.NewAdapter("sqlite3", dburl, true)
    if err != nil {
        return nil, nil, err
    }

    // Load RBAC model
    m, err := model.NewModelFromFile("config/rbac_model.conf")
    if err != nil {
        return nil, nil, err
    }

    enforcer, err := casbin.NewEnforcer(m, adapter)
    if err != nil {
        return nil, nil, err
    }

    enforcer.LoadPolicy()
    return enforcer, db, nil
}
```

## 4. Setting up Casbin for RBAC

The system uses Casbin for role-based access control:

1. Model definition in `rbac_model.conf`
2. Policy management through the adapter
3. Role assignment in handlers

Example user registration handler:

```go
func RegisterUser(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Create user
        user := User{
            Username: req.Username,
            Password: req.Password, // TODO: Add password hashing
            Role:     "user",
        }
        if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }

        // Assign role
        if _, err := enforcer.AddGroupingPolicy(user.Username, "user"); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
            return
        }

        // Grant permissions
        if _, err := enforcer.AddPolicy(user.Username, "/profile", "GET"); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign policy"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
    }
}
```

## 5. Integrating with Gin

Server initialization and route configuration:

```go
func StartServer() {
    // Initialize Casbin and database
    enforcer, db, err := database.SetupCasbin()
    if err != nil {
        log.Fatalf("Error initializing Casbin and DB: %v", err)
    }

    // Configure Gin router
    r := gin.Default()
    
    // Register routes
    r.POST("/register", handler.RegisterUser(db, enforcer))

    // Start server
    r.Run(":8080")
}
```