package http

import (
	"net/http"

	"github.com/BekzatS8/buhpro/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *services.UserUsecase
}

func NewUserHandler(u *services.UserUsecase) *UserHandler {
	return &UserHandler{uc: u}
}

// RegisterRoutes not used in router (we register handlers directly), but keep for parity
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
	rg.POST("/refresh", h.Refresh)
	rg.POST("/logout", h.Logout)
	rg.GET("/count", h.Count)
}

// Requests
type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Role == "" {
		req.Role = "executor"
	}
	access, refresh, err := h.uc.Register(req.Email, req.Phone, req.FullName, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"access_token": access, "refresh_token": refresh})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := h.uc.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *UserHandler) Logout(c *gin.Context) {
	// This endpoint is protected by AuthMiddleware (access token).
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := uidVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
		return
	}
	if err := h.uc.Logout(c.Request.Context(), uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *UserHandler) Count(c *gin.Context) {
	cnt, err := h.uc.RepoCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users_count": cnt})
}

func (h *UserHandler) Me(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := uidVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
		return
	}

	u, err := h.uc.GetProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

// UpdateMe — обновление профиля текущего пользователя
type updateMeReq struct {
	FullName *string                `json:"full_name,omitempty"`
	Phone    *string                `json:"phone,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := uidVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
		return
	}

	var req updateMeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	upd := services.UserUpdate{
		FullName: req.FullName,
		Phone:    req.Phone,
		Metadata: req.Metadata,
	}

	u, err := h.uc.UpdateProfile(c.Request.Context(), uid, upd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}
