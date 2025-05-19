package votes

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example.com/se/config"
	"example.com/se/entity"
)

func GetAll(c *gin.Context) {
	var votes []entity.Votes

	db := config.DB()
	result := db.Find(&votes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, votes)
}

func Get(c *gin.Context) {
	id := c.Param("id")
	var vote entity.Votes

	db := config.DB()
	result := db.First(&vote, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vote not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, vote)
}

func CreateVote(c *gin.Context) {
	var payload struct {
		UserID      uint      `json:"user_id"`
		CandidateID uint      `json:"candidate_id"`
		ElectionID  uint      `json:"election_id"`
		Timestamp   time.Time `json:"timestamp"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now()
	}

	db := config.DB()

	var existingVote entity.Votes
	err := db.Where("user_id = ? AND election_id = ?", payload.UserID, payload.ElectionID).First(&existingVote).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User has already voted in this election"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	vote := entity.Votes{
		UserID:      payload.UserID,
		CandidateID: payload.CandidateID,
		ElectionID:  payload.ElectionID,
		Timestamp:   payload.Timestamp,
	}

	if err := db.Create(&vote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Vote recorded successfully",
		"vote":    vote,
	})
}
