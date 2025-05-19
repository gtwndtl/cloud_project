package entity

import (
	"gorm.io/gorm"
	"time"
)

type Votes struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	User   *Users `gorm:"foreignKey: user_id" json:"user"`

	CandidateID uint       `json:"candidate_id"`
	Candidate   *Candidates `gorm:"foreignKey: candidate_id" json:"candidate"`

	ElectionID uint      `json:"election_id"`
	Election   *Elections `gorm:"foreignKey: election_id" json:"election"`
	Timestamp  time.Time `json:"timestamp"`
}
