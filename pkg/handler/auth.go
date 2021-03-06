package handler

import (
	server "allincecup-server"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"time"
)

const refreshTokenTTL = 1440 * time.Hour

func (h *Handler) signUp(ctx *gin.Context) {
	var input server.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// email, password, phone_number validation
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid email")
		return
	}
	if len(input.Password) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid password")
		return
	}
	if len(input.PhoneNumber) < 10 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid phone_number")
		return
	}

	id, roleId, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"role_id": roleId,
	})
}

func (h *Handler) createModerator(ctx *gin.Context) {
	var input server.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// email, password, phone_number validation
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid email")
		return
	}
	if len(input.Password) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid password")
		return
	}
	if len(input.PhoneNumber) < 10 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid phone_number")
		return
	}

	id, roleId, err := h.services.Authorization.CreateModerator(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"role_id": roleId,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.GenerateTokens(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, userRoleId, err := h.services.ParseToken(accessToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSession, err := h.services.CreateNewSession(&server.Session{
		UserId:       userId,
		RoleId:       userRoleId,
		RefreshToken: refreshToken,
		ClientIp:     c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
		CreatedAt:    time.Now(),
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "unable to create new session: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"session_id":         newSession.Id,
		"refresh_expires_at": newSession.ExpiresAt,
	})

}

type RefreshTokensInput struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refresh(ctx *gin.Context) {
	var input RefreshTokensInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := h.services.Authorization.RefreshAccessToken(input.RefreshToken)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Authorization.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}
