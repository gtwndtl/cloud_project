package votes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/metrics"

	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	pushGatewayURL = "http://pushgateway:9091"
	jobName        = "votes_job"
)

// pushVoteMetrics จะ push ทั้ง VotesCreatedTotal, VotesRejectedTotal และ VotesPerElectionTotal
func pushVoteMetrics(electionID uint) {
	// อัปเดต votes_per_election_total label ตาม election_id
	count := getVoteCount(electionID)
	metrics.VotesPerElectionTotal.WithLabelValues(fmt.Sprint(electionID)).Set(float64(count))

	// Push ทุก metric ที่เกี่ยวกับ vote
	if err := push.New(pushGatewayURL, jobName).
		Collector(metrics.VotesCreatedTotal).
		Collector(metrics.VotesRejectedTotal).
		Collector(metrics.VotesPerElectionTotal).
		Push(); err != nil {
		log.Printf("Could not push vote metrics to Pushgateway: %v", err)
	}
}

// getVoteCount ดึงจำนวน vote ของ election หนึ่ง ๆ จาก DB
func getVoteCount(electionID uint) int {
	var count int64
	config.DB().Model(&entity.Votes{}).
		Where("election_id = ?", electionID).
		Count(&count)
	return int(count)
}

func GetAll(c *gin.Context) {
	var votes []entity.Votes
	db := config.DB()
	if err := db.Find(&votes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, votes)
}

func Get(c *gin.Context) {
	id := c.Param("id")
	var vote entity.Votes
	db := config.DB()
	if err := db.First(&vote, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vote not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		metrics.VotesRejectedTotal.Inc()
		pushVoteMetrics(payload.ElectionID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now()
	}

	db := config.DB()

	// ตรวจสอบว่ามีการโหวตซ้ำหรือไม่
	var existingVote entity.Votes
	err := db.
		Where("user_id = ? AND election_id = ?", payload.UserID, payload.ElectionID).
		First(&existingVote).Error

	if err == nil {
		// โหวตซ้ำ
		metrics.VotesRejectedTotal.Inc()
		pushVoteMetrics(payload.ElectionID)
		c.JSON(http.StatusConflict, gin.H{"error": "User has already voted in this election"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// DB error
		metrics.VotesRejectedTotal.Inc()
		pushVoteMetrics(payload.ElectionID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	// สร้าง vote ใหม่
	vote := entity.Votes{
		UserID:      payload.UserID,
		CandidateID: payload.CandidateID,
		ElectionID:  payload.ElectionID,
		Timestamp:   payload.Timestamp,
	}

	if err := db.Create(&vote).Error; err != nil {
		metrics.VotesRejectedTotal.Inc()
		pushVoteMetrics(payload.ElectionID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote", "details": err.Error()})
		return
	}

	// สำเร็จ
	metrics.VotesCreatedTotal.Inc()
	pushVoteMetrics(payload.ElectionID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Vote recorded successfully",
		"vote":    vote,
	})
}
