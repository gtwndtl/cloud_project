package elections

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/metrics"
)

const (
	pushGatewayURL = "http://pushgateway:9091"
	jobName        = "elections_job"
)

// pushCounter ส่ง counter เดียวไปยัง Pushgateway
func pushCounter(counter prometheus.Counter) {
	_ = push.New(pushGatewayURL, jobName).
		Collector(counter).
		Push()
}

// GetAll ดึงรายการทั้งหมด
func GetAll(c *gin.Context) {
	var elections []entity.Elections
	db := config.DB()
	if err := db.Find(&elections).Error; err != nil {
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, elections)
}

// Get ดึง election ตาม ID
func Get(c *gin.Context) {
	id := c.Param("id")
	var election entity.Elections
	db := config.DB()
	if err := db.First(&election, id).Error; err != nil {
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, election)
}

// Create สร้าง election ใหม่
func Create(c *gin.Context) {
	var payload struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Status      string    `json:"status"`
		CandidateID uint      `json:"candidate_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

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
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create election", "details": err.Error()})
		return
	}

	metrics.ElectionsCreatedTotal.Inc()
	pushCounter(metrics.ElectionsCreatedTotal)
	c.JSON(http.StatusCreated, gin.H{"message": "Election created successfully", "election": election})
}

// Update แก้ไข election
func Update(c *gin.Context) {
	id := c.Param("id")
	var election entity.Elections
	db := config.DB()
	if err := db.First(&election, id).Error; err != nil {
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
		return
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Status      string    `json:"status"`
		CandidateID uint      `json:"candidate_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request", "details": err.Error()})
		return
	}

	// อัปเดต field ที่ต้องการ
	election.Title = input.Title
	election.Description = input.Description
	election.StartTime = input.StartTime
	election.EndTime = input.EndTime
	election.Status = input.Status
	election.CandidateID = input.CandidateID

	if err := db.Save(&election).Error; err != nil {
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update election"})
		return
	}

	metrics.ElectionsUpdatedTotal.Inc()
	pushCounter(metrics.ElectionsUpdatedTotal)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

// Delete ลบ election
func Delete(c *gin.Context) {
	id := c.Param("id")
	db := config.DB()
	if tx := db.Delete(&entity.Elections{}, id); tx.RowsAffected == 0 {
		metrics.ElectionsFailuresTotal.Inc()
		pushCounter(metrics.ElectionsFailuresTotal)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found"})
		return
	}

	metrics.ElectionsDeletedTotal.Inc()
	pushCounter(metrics.ElectionsDeletedTotal)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
