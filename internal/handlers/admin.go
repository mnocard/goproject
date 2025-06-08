package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uService "github.com/mnocard/go-project/internal/services/user"
)

var Rating = map[string]int{}

var tasks []task

type task struct {
	Id         int    `json:"username" binding:"required"`
	MainTaskId int    `json:"maintaskid" binding:"required"`
	Points     int    `json:"points"`
	User       string `json:"user" binding:"required"`
	IsComplete bool   `json:"iscomplete"`
}

type Handler struct {
	uService userService
}

type userService interface {
	Create(context.Context, string, string, bool) (int, error)
	Get(context.Context, int) (*uService.User, error)
	FindByName(context.Context, string) (*uService.User, error)
	Update(context.Context, string, string) (int, error)
	Delete(context.Context, string) (bool, error)
}

func New(uService userService) *Handler {
	return &Handler{uService: uService}
}

// @Summary	ChangeAdminPassword
// @Schemes	http
// @Accept		json
// @Produce	json
// @Param		input	body	handlers.Password	true	"new password"
// @Router		/admin/changeAdminPassword [post]
func (h *Handler) ChangeAdminPassword(c *gin.Context) {
	log.Println("ChangeAdminPassword start. Request", c.Request)
	user := c.MustGet(gin.AuthUserKey).(string)

	var json struct {
		NewPassword string `json:"password" binding:"required"`
	}

	err := c.Bind(&json)
	if err != nil {
		log.Println("ChangeAdminPassword err", err)
		return
	}

	log.Println("ChangeAdminPassword json", json)
	h.uService.Update(c.Request.Context(), user, json.NewPassword)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	log.Println("ChangeAdminPassword end")
}

func (h *Handler) CreateUser(c *gin.Context) {
	log.Println("CreateUser start. Request", c.Request)

	var json struct {
		UserName string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		IsAdmin  bool   `json:"is_admin" `
	}

	if c.Bind(&json) == nil {
		log.Println("CreateUser json", json)

		h.uService.Create(c.Request.Context(), json.UserName, json.Password, json.IsAdmin)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	log.Println("CreateUser end")
}

func CreateTask(c *gin.Context) {
	log.Println("CreateTask start. Request", c.Request)

	var task task

	if c.Bind(&task) == nil {
		log.Println("CreateTask task", task)

		tasks = append(tasks, task)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	log.Println("CreateTask end")
}

func CreateSubTask(c *gin.Context) {
	log.Println("CreateSubTask start. Request", c.Request)

	var task task

	if c.Bind(&task) == nil {
		log.Println("CreateSubTask task", task)

		tasks = append(tasks, task)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	log.Println("CreateSubTask end")
}
