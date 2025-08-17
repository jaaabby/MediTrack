package models

type SupplyMovement struct {
	ID              string `json:"id"`
	MedicalSupplyID string `json:"medical_supply_id"`
	LocationID      string `json:"location_id"`
	Status          string `json:"status"`
	Quantity        int    `json:"quantity"`
}
