package server

import (
	"my-casbin/internal/database"
	"my-casbin/internal/handler"
	"my-casbin/internal/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	// m, _ := model.NewModelFromFile("config/rbac_model.conf")
	// a := fileadapter.NewAdapter("config/policy.csv")
	// e, _ := casbin.NewEnforcer(m, a)
	// e.LoadPolicy()

	e, db, err := database.SetupCasbin()
	if err != nil {
		panic(err)
	}

	r.POST("/register", handler.RegisterUser(db, e))

	r.Use(middleware.JWTMiddleware())

	r.GET("/admin", middleware.Authorize(e), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello Admin!"})
	})

	r.GET("/user", middleware.Authorize(e), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello User!"})
	})

	// ✅ Protected with JWT + Casbin
	protected := r.Group("/")
	protected.Use(middleware.JWTMiddleware())
	protected.Use(middleware.CasbinMiddleware(e))

	protected.GET("/profile", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + username.(string)})
	})

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
