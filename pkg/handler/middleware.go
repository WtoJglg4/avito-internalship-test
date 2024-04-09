package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userIdCtx           = "userId"
	userRoleCtx         = "userRole"
)

func (h *Handler) userIdentity(c *gin.Context) (int, string, error) {

	header := c.GetHeader(authorizationHeader)
	if len(header) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return 0, "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return 0, "", errors.New("invalid auth header")
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return 0, "", errors.New("token is empty")
	}

	logrus.Println(headerParts)

	userId, userRole, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		logrus.Errorln(err)
		newErrorResponse(c, http.StatusUnauthorized, "failed to parse token")
		return 0, "", errors.New("failed to parse token")
	}

	if userId < 0 || userRole != "admin" && userRole != "default" {
		logrus.Errorln(err)
		newErrorResponse(c, http.StatusUnauthorized, "invalid userId or userRole")
		return 0, "", errors.New("invalid userId or userRole")
	}
	// logrus.Printf("userId: %v, userRole: %v", userId, userRole)

	// c.Header(userIdCtx, fmt.Sprintf("%d", userId))
	// c.Header(userRoleCtx, userRole)

	return userId, userRole, nil
}

// func getUserId(c *gin.Context) (int, error) {
// 	// userId := w.Header().Get(userIdCtx)
// 	userId := c.GetHeader(userIdCtx)
// 	id, err := strconv.Atoi(userId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, "user id is not found")
// 		logrus.Error("user id is not found")
// 		return 0, errors.New("user id is not found")
// 	}
// 	return id, nil
// }

// func getUserRole(c *gin.Context) (string, error) {
// 	userRole := c.GetHeader(userRoleCtx)
// 	logrus.Println(userRole, userRole == "admin")
// 	if userRole != "admin" {
// 		newErrorResponse(c, http.StatusForbidden, "not enough rights")
// 		return "", errors.New("not enough rights")
// 	}
// 	return userRole, nil
// }
