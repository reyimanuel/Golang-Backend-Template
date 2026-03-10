package api

import (
	user "github.com/reyimanuel/template/internal/api/users"
	"github.com/reyimanuel/template/internal/infrastructures/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	user.RegisterRoutes(api.Group("/users"), db)

	protected := api.Group("/")
	protected.Use(middleware.MiddlewareAuth)

	// Protected routes example
	// api.RegisterRoutes(protected.Group("/api"), db)
}
