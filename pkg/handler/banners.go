package handler

import (
	"github/avito/entities"
	"github/avito/pkg/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) createBanner(c *gin.Context) {
	logrus.Println(c.Request.Method, c.Request.URL)
	_, userRole, err := h.userIdentity(c)
	if err != nil {
		return
	}
	if userRole != "admin" {
		newErrorResponse(c, http.StatusForbidden, "not enough rights")
		return
	}

	var bannerPostRequest entities.BannerPostRequest
	if err := c.BindJSON(&bannerPostRequest); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := bannerPostRequest.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	banner := entities.Banner{
		TagIDs:    bannerPostRequest.TagIDs,
		FeatureID: bannerPostRequest.FeatureID,
		Content:   bannerPostRequest.Content,
		IsActive:  bannerPostRequest.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := h.service.CreateBanner(banner)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"banner_id": id,
	})
}

// Получение всех баннеров c фильтрацией по фиче и/или тегу
func (h *Handler) getAllBanners(c *gin.Context) {
	logrus.Println(c.Request.Method, c.Request.URL)
	_, userRole, err := h.userIdentity(c)
	if err != nil {
		return
	}
	if userRole != "admin" {
		newErrorResponse(c, http.StatusForbidden, "not enough rights")
		return
	}

	filters, err := entities.GetAllQueryParams(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	banners, err := h.service.GetAllBanners(filters)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, banners)
}

func (h *Handler) deleteBanners(c *gin.Context) {
	logrus.Println(c.Request.Method, c.Request.URL)
	_, userRole, err := h.userIdentity(c)
	if err != nil {
		return
	}
	if userRole != "admin" {
		newErrorResponse(c, http.StatusForbidden, "not enough rights")
		return
	}

	filters, err := entities.GetAllQueryParams(c)
	if err != nil || filters.Feature_id == -1 && filters.Tags_id == -1 {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.DeleteBanners(filters); err != nil {
		if err == repository.ErrNoRowsDeleted {
			c.Status(http.StatusNotFound)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) userBanner(c *gin.Context) {
	logrus.Println(c.Request.Method, c.Request.URL)
	_, _, err := h.userIdentity(c)
	if err != nil {
		return
	}

	filters, err := entities.UserGetQueryParams(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	bannerContent, err := h.service.UserBanner(filters)
	if err != nil {
		if err == repository.ErrNoRowsSelected {
			c.Status(http.StatusNotFound)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, bannerContent)
}
