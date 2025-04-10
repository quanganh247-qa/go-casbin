# Developer Guide: Integrating GORM + SQLite + Casbin (RBAC) in a Gin Project

This guide walks you through the main steps of setting up a Gin web application with GORM for data persistence using SQLite and Casbin for role-based access control (RBAC).

---

## 1. Prerequisites

- Go installed (as defined in the go.mod file, at least 1.23.4).
- A Mac machine with a terminal (or use your IDE’s integrated terminal).
- Basic understanding of Gin, GORM, and Casbin.

---

## 2. Project Structure

Your repository is structured approximately as follows:

```
my-casbin
├─ cmd/api
│  ├─ rbac_model.conf         // Casbin RBAC model configuration
│  └─ policy.csv              // Casbin policy file (if using file-based policies)
├─ internal
│  ├─ database
│  │  └─ database.go          // Database/GORM & Casbin initialization
│  └─ handler                 // HTTP request handlers (e.g., RegisterUser)
│     └─ handler.go
├─ go.mod
├─ Makefile
└─ README.md
```

---

## 3. Configuring GORM with SQLite

In your `/internal/database/database.go` file, the `SetupCasbin()` function demonstrates how to:

- Open an SQLite database connection using GORM.
- Create a Casbin adapter based on the same SQLite file.

Example snippet:
    
````go
// filepath: [database.go](http://_vscodecontentref_/0)
func SetupCasbin() (*casbin.Enforcer, *gorm.DB, error) {
    // Connect to SQLite DB using GORM
    db, err := gorm.Open(sqlite.Open(dburl), &gorm.Config{})
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect database: %w", err)
    }

    // Use xorm adapter for Casbin; Note: adjust the database url and options as needed
    adapter, err := xormadapter.NewAdapter("sqlite3", dburl, true)
    if err != nil {
        return nil, nil, err
    }

    // Load the Casbin model (RBAC model definition)
    m, err := model.NewModelFromFile("config/rbac_model.conf")
    if err != nil {
        return nil, nil, err
    }

    enforcer, err := casbin.NewEnforcer(m, adapter)
    if err != nil {
        return nil, nil, err
    }

    // Load all current policies from persistent storage
    _ = enforcer.LoadPolicy()
    return enforcer, db, nil
}