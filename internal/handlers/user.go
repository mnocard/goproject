package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//	@Summary	GetRating
//	@Tags		User
//	@Security	BasicAuth
//	@Accept		json
//	@Produce	json
//	@Router		/user/getRating [get]
func (h *Handler) GetRating(c *gin.Context) {
	log.Println("GetRating start. Request", c.Request)

	user := c.MustGet(gin.AuthUserKey).(string)
	log.Println("GetRating user", user)

	rating, err := h.uService.GetRating(c.Request.Context(), user)
	if err != nil {
		log.Println("GetRating err", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rating": rating})
	log.Println("GetRating end")
}

//	@Summary	CompleteTask
//	@Tags		User
//	@Security	BasicAuth
//	@Param		taskId	path	int	true	"Task ID"
//	@Router		/user/completeTask/{taskId} [get]
func (h *Handler) CompleteTask(c *gin.Context) {
	log.Println("CompleteTask start. Request", c.Request)

	user := c.MustGet(gin.AuthUserKey).(string)
	log.Println("CompleteTask user", user)

	ctx := c.Request.Context()
	taskIdQuery := c.Param("taskId")
	taskId, err := strconv.Atoi(taskIdQuery)
	if err != nil {
		log.Println("CompleteTask err, can't parse task id: " + taskIdQuery)
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "error": "can't parse task id"})
		return
	}

	if c.Bind(&taskId) == nil {
		log.Println("CompleteTask taskToComplete", taskId)

		userData, err := h.uService.FindByName(ctx, user)
		if err != nil {
			log.Println("CompleteTask err", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "error": err.Error()})
			return
		}

		points, err := h.tService.CompleteTask(ctx, userData.Id, taskId)
		if err != nil {
			log.Println("CompleteTask err", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "error": err.Error()})
			return
		}

		err = h.uService.UpdateRating(ctx, userData.Id, userData.Rating+points)
		if err != nil {
			log.Println("CompleteTask err", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"is_completed": true})
	}

	log.Println("CompleteTask end")
}
