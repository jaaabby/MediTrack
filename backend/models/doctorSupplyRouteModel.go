package models

type DoctorSupplyRoute struct {
	ID            string `json:"id"`
	DoctorID      string `json:"doctor_id"`
	SupplyRouteID string `json:"supply_route_id"`
}
