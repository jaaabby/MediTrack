package models

import "time"

type BatchHistory struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	DateTime       time.Time `json:"date_time" db:"date_time" gorm:"not null"`
	ChangeDetails  string    `json:"change_details" db:"change_details" gorm:"not null"`
	PreviousValues string    `json:"previous_values" db:"previous_values" gorm:"type:jsonb"`
	NewValues      string    `json:"new_values" db:"new_values" gorm:"type:jsonb"`
	UserName       string    `json:"user_name" db:"user_name" gorm:"not null"`
	BatchID        *int      `json:"batch_id" db:"batch_id" gorm:"default:null"`
	UserRUT        string    `json:"user_rut" db:"user_rut" gorm:"not null"`
	BatchNumber    int       `json:"batch_number" db:"batch_number" gorm:"not null"`
}

func (b BatchHistory) TableName() string {
	return "batch_history"
}
