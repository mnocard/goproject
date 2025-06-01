package handlers

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func CompleteTask(c *gin.Context) {
// 	log.Println("CompleteTask start. Request", c.Request)

// 	var query struct {
// 		TaskId string `json:"value" binding:"taskid"`
// 	}

// 	if c.Bind(&query) == nil {
// 		log.Println("CreateSubTask task", query)

// 		tasks = append(tasks, task)
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	}
// 	log.Println("CreateSubTask end")
// }
