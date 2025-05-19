package elections

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"example.com/se/config"
	"example.com/se/entity"
)

func GetAll(c *gin.Context) {
	var elections []entity.Elections

	db := config.DB()
	result := db.Find(&elections)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, elections)
}

func Get(c *gin.Context) {
	ID := c.Param("id")
	var election entity.Elections

	db := config.DB()
	result := db.First(&election, ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	if election.ID == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, election)
}

func Update(c *gin.Context) {
	var election entity.Elections
	ElectionID := c.Param("id")

	db := config.DB()
	result := db.First(&election, ElectionID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
		return
	}

	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, unable to map payload"})
		return
	}

	result = db.Save(&election)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	if tx := db.Exec("DELETE FROM elections WHERE id = ?", id); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func Create(c *gin.Context) {
	var payload struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Status      string    `json:"status"`
		CandidateID uint      `json:"candidate_id"`
	}

	// ตรวจสอบความถูกต้องของ input
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// สร้าง instance ของ Election โดยกำหนดค่าชัดเจน
	election := entity.Elections{
		Title:       payload.Title,
		Description: payload.Description,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		Status:      payload.Status,

		CandidateID: payload.CandidateID,
	}

	db := config.DB()
	if err := db.Create(&election).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create election", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Election created successfully",
		"election": election,
	})
}
