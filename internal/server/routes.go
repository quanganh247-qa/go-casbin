package server

import (
	"my-casbin/internal/middleware"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
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

	r.Use(middleware.JWTMiddleware())

	m, _ := model.NewModelFromFile("config/rbac_model.conf")
	a := fileadapter.NewAdapter("config/policy.csv")
	e, _ := casbin.NewEnforcer(m, a)
	e.LoadPolicy()

	r.GET("/admin", middleware.Authorize(e), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello Admin!"})
	})

	r.GET("/user", middleware.Authorize(e), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello User!"})
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
