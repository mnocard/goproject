package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tService "github.com/mnocard/go-project/internal/services/task"
	uService "github.com/mnocard/go-project/internal/services/user"
)

type Handler struct {
	uService userService
	tService taskService
}

type userService interface {
	Create(context.Context, *uService.User) (int, error)
	Get(context.Context, int) (*uService.User, error)
	FindByName(context.Context, string) (*uService.User, error)
	Update(context.Context, string, string) (int, error)
	Delete(context.Context, string) (bool, error)
}

type taskService interface {
	Create(context.Context, *tService.Task) (int, error)
	Get(context.Context, int) (*tService.Task, error)
	FindTaskByUserId(context.Context, int) (*tService.Task, error)
	Update(context.Context, *tService.Task) (int, error)
	Delete(context.Context, int) (bool, error)
}

func New(uService userService, tService taskService) *Handler {
	return &Handler{
		uService: uService,
		tService: tService,
	}
}

type NewPassword struct {
	Value string `json:"value" binding:"required"`
}

// @Summary	ChangeAdminPassword
// @Tags		Admin
// @Security	BasicAuth
// @Accept		json
// @Produce	json
// @Param		input	body	handlers.NewPassword	true	"new password"
// @Router		/admin/changeAdminPassword [post]
func (h *Handler) ChangeAdminPassword(c *gin.Context) {
	log.Println("ChangeAdminPassword start. Request", c.Request)
	user := c.MustGet(gin.AuthUserKey).(string)
	isAdmin := c.MustGet("isAdmin").(bool)
	if !isAdmin {
		log.Println("ChangeAdminPassword not admin")
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": "not admin"})
		return
	}

	var json NewPassword

	err := c.Bind(&json)
	if err != nil {
		log.Println("ChangeAdminPassword err", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}

	log.Println("ChangeAdminPassword json", json)
	h.uService.Update(c.Request.Context(), user, json.Value)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	log.Println("ChangeAdminPassword end")
}

// @Summary		CreateUser
// @Tags		Admin
// @Security	BasicAuth
// @Accept		json
// @Produce		json
// @Param		input	body	user.User	true	"user data"
// @Router		/admin/createUser [post]
func (h *Handler) CreateUser(c *gin.Context) {
	log.Println("CreateUser start. Request", c.Request)

	var user uService.User

	if c.Bind(&user) == nil {
		log.Println("CreateUser user", user)

		h.uService.Create(c.Request.Context(), &user)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	log.Println("CreateUser end")
}

// @Summary		CreateTask
// @Tags		Admin
// @Security	BasicAuth
// @Accept		json
// @Produce		json
// @Param		input	body	task.Task	true	"task data"
// @Router		/admin/createTask [post]
func (h *Handler) CreateTask(c *gin.Context) {
	log.Println("CreateTask start. Request", c.Request)

	var task tService.Task

	if c.Bind(&task) == nil {
		log.Println("CreateTask task", task)

		h.tService.Create(c.Request.Context(), &task)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	log.Println("CreateTask end")
}
