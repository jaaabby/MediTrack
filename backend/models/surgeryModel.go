package models

// Surgery representa un tipo de cirugía
type Surgery struct {
	ID       int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string  `json:"name" gorm:"not null;unique"`
	Duration float64 `json:"duration" gorm:"not null"` // Duración en horas
}

func (s Surgery) TableName() string {
	return "surgery"
}
