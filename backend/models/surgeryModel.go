package models

// Surgery representa un tipo de cirugía
type Surgery struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"not null;unique"`
	Duration    float64 `json:"duration" gorm:"not null"` // Duración en horas
	SpecialtyID *int    `json:"specialty_id"`

	// Relaciones
	Specialty           *MedicalSpecialty     `json:"specialty,omitempty" gorm:"foreignKey:SpecialtyID"`
	TypicalSupplies    []SurgeryTypicalSupply `json:"typical_supplies,omitempty" gorm:"foreignKey:SurgeryID"`
}

func (s Surgery) TableName() string {
	return "surgery"
}
