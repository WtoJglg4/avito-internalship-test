package handler

import (
	"github/avito/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// add auth endpoints to API
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	router.GET("/user_banner", h.userBanner)

	banners := router.Group("/banner")
	{
		banners.GET("/", h.getAllBanners)
		banners.POST("/", h.createBanner)
		banners.DELETE("/", h.deleteBanners)

		banner := banners.Group("/:id")
		{
			banner.PATCH("/", h.updateBanner)
			banner.DELETE("/", h.deleteBanner)

			versions := banner.Group("/version")
			{
				versions.GET("/", h.getAllVersions)
				versions.GET("/:version", h.getVersion)
			}
		}
	}

	return router
}
