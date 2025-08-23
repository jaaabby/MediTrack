package models

type SupplyCode struct {
	Code         int    `json:"code" gorm:"primaryKey"`
	Name         string `json:"name" db:"name" gorm:"not null"`
	CodeSupplier int    `json:"code_supplier" db:"code_supplier" gorm:"not null"`
}

func (s SupplyCode) TableName() string {
	return "supply_code"
}
