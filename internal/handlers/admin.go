package handlers

import (
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
	Create(uName, password string) (int, error)
	Get(id int) (*uService.User, error)
	FindByName(uName string) (*uService.User, error)
	Update(uName, password string) (int, error)
	Delete(uName string) (bool, error)
}

func New(uService userService) *Handler {
	return &Handler{uService: uService}
}

func (h *Handler) ChangeAdminPassword(c *gin.Context) {
	log.Println("ChangeAdminPassword start. Request", c.Request)
	user := c.MustGet(gin.AuthUserKey).(string)

	var json struct {
		Value string `json:"value" binding:"required"`
	}

	if c.Bind(&json) == nil {
		log.Println("ChangeAdminPassword json.Value", json.Value)

		h.uService.Update(user, json.Value)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	log.Println("ChangeAdminPassword end")
}

func (h *Handler) CreateUser(c *gin.Context) {
	log.Println("CreateUser start. Request", c.Request)

	var json struct {
		UserName string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if c.Bind(&json) == nil {
		log.Println("CreateUser json", json)

		h.uService.Create(json.UserName, json.Password)
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
