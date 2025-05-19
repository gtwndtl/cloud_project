package entity
import "gorm.io/gorm"

type Notifications struct {
	gorm.Model
    ToEmail string `json:"to_email"`
    Subject string `json:"subject"`
    Message string `json:"message"`
    Status  string `json:"status"` // e.g., sent, failed
}
