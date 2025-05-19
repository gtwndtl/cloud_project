package entity

import (
	"gorm.io/gorm"
	"time"
)

type Elections struct {
	gorm.Model

	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`

	CandidateID uint       `json:"candidate_id"`
	Candidate   *Candidates `gorm:"foreignKey: candidate_id" json:"candidate"`
}
