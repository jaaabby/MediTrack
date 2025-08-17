package models

type SupplyRouteMedicalSupply struct {
	ID              string `json:"id"`
	SupplyRouteID   string `json:"supply_route_id"`
	MedicalSupplyID string `json:"medical_supply_id"`
}
