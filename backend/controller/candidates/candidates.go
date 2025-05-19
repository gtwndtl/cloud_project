package candidates

import (
	"net/http"

	"example.com/se/config"
	"example.com/se/entity"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	var candidates []entity.Candidates

	db := config.DB()
	result := db.Find(&candidates)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, candidates)
}

func Get(c *gin.Context) {
	ID := c.Param("id")
	var candidate entity.Candidates

	db := config.DB()
	result := db.First(&candidate, ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	if candidate.ID == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, candidate)
}

func Create(c *gin.Context) {
	var payload struct {
		Name       string `json:"name"`
		ElectionID uint   `json:"election_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	candidate := entity.Candidates{
		Name:       payload.Name,
		ElectionID: payload.ElectionID,
	}

	db := config.DB()
	if err := db.Create(&candidate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create candidate", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Candidate created successfully",
		"candidate": candidate,
	})
}

func Update(c *gin.Context) {
	var candidate entity.Candidates
	candidateID := c.Param("id")

	db := config.DB()
	result := db.First(&candidate, candidateID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
		return
	}

	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	result = db.Save(&candidate)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update candidate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Candidate updated successfully"})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()

	if tx := db.Exec("DELETE FROM candidates WHERE id = ?", id); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Candidate not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Candidate deleted successfully"})
}
