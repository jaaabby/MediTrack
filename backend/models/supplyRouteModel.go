package models

type SupplyRoute struct {
	ID              string `json:"id"`
	PatientID       string `json:"patient_id"`
	OperatingRoomID string `json:"operating_room_id"`
	Name            string `json:"name"`
}
