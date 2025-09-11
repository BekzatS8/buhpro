package http

import (
	"github.com/BekzatS8/buhpro/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	uc *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: u}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
	rg.GET("/count", h.Count)
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // client|executor|coach
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
	token, err := h.uc.Register(req.Email, req.Phone, req.FullName, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
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
	token, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Count(c *gin.Context) {
	// simple check using repo via usecase not strictly necessary
	cnt, err := h.uc.RepoCount()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"users_count": cnt})
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

	upd := usecase.UserUpdate{
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
