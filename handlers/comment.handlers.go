package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/helpers"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
)

// Get comments
func (m *Repo) HandleGetComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("videoID"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "404 video not found!",
			"errors": err,
		})
		return
	}

	var page, limit int
	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "24"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var comments []models.CommentDTO
	var count int64
	comments, count, err = m.App.DBMethods.GetCommentsByVideoID(id, page, limit)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var hasNextPage bool

	if count > int64(page*limit) {
		hasNextPage = true
	} else {
		hasNextPage = false
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"page":          page,
		"comments":      comments,
		"has_next_page": hasNextPage,
	})
}

func (m *Repo) HandleCreateOrUpdateComment(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.IndentedJSON(400, gin.H{
			"error": "Invalid User ID",
		})
		return
	}

	user, err := helpers.ValidateAndGetUserByID(m.App, user_id)

	if err != nil {
		c.IndentedJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var payload models.Comment
	err = c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, gin.H{
			"error": "Invalid Comment Payload",
		})
		return
	}

	if payload.VideoID <= 0 {
		c.IndentedJSON(400, gin.H{
			"error": "Invalid Video ID",
		})
		return
	}

	if user.ID <= 0 {
		c.IndentedJSON(400, gin.H{
			"error": "Invalid User",
		})
		return
	}

	_, err = m.App.DBMethods.GetVideoByID(int(payload.VideoID))
	if err != nil {
		c.IndentedJSON(400, gin.H{
			"error": "Invalid Video ID",
		})
		return
	}

	// Validate the payload
	validator := validator.New()
	validator.IsLength(payload.Text, "text", 2, 1000)

	if validator.Valid() == false {
		c.IndentedJSON(400, gin.H{
			"error": validator.GetErrMsg(),
		})
		return
	}

	if payload.ID > 0 {
		// update comment
		comment, err := m.App.DBMethods.GetCommnetByID(payload.ID)
		if err != nil {
			c.IndentedJSON(400, gin.H{
				"error": "The comment does not exist",
			})
			return
		}

		if comment.UserID != user.ID || user.UserRole != "admin" {
			c.IndentedJSON(400, gin.H{
				"error": "You are not allowed to update this comment",
			})
			return
		}

		comment.Text = payload.Text

		err = m.App.DBMethods.UpdateComment(comment)
		if err != nil {
			c.IndentedJSON(500, gin.H{
				"error": "Something went wrong. Please try again later",
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "Comment updated successfully",
		})

		return
	} else {
		// create comment
		comment := models.Comment{
			Text:    payload.Text,
			UserID:  user.ID,
			VideoID: payload.VideoID,
		}

		comment_id, err := m.App.DBMethods.CreateComment(&comment)
		if err != nil {
			c.IndentedJSON(500, gin.H{
				"error": "Something went wrong. Please try again later",
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "Comment created successfully",
			"id":      comment_id,
		})

		return
	}
}

func (m *Repo) HandleDeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Params.ByName("commentID"))
	if err != nil {
		c.JSON(400, gin.H{"error": "The comment you are trying to delete does not exist"})
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		c.IndentedJSON(400, gin.H{
			"error": "Invalid User ID",
		})
		return
	}

	user, err := helpers.ValidateAndGetUserByID(m.App, user_id)
	if err != nil {
		c.IndentedJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	comment, err := m.App.DBMethods.GetCommnetByID(commentID)
	if err != nil {
		c.IndentedJSON(400, gin.H{
			"error": "The comment you are trying to delete does not exist",
		})
		return
	}

	if comment.UserID != user.ID || user.UserRole != "admin" {
		c.IndentedJSON(400, gin.H{
			"error": "You are not allowed to delete this comment",
		})
		return
	}

	err = m.App.DBMethods.DeleteCommentByID(commentID)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"error": "Something went wrong. Please try again later",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Comment deleted successfully",
	})
}
