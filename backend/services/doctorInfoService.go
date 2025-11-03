package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type DoctorInfoService struct {
	DB *gorm.DB
}

func NewDoctorInfoService(db *gorm.DB) *DoctorInfoService {
	return &DoctorInfoService{DB: db}
}

// CreateDoctorInfo crea información extendida de un doctor
func (s *DoctorInfoService) CreateDoctorInfo(doctorInfo *models.DoctorInfo) error {
	return s.DB.Create(doctorInfo).Error
}

// GetDoctorInfoByRUT obtiene información de un doctor por RUT
func (s *DoctorInfoService) GetDoctorInfoByRUT(rut string) (*models.DoctorInfo, error) {
	var doctorInfo models.DoctorInfo
	if err := s.DB.Preload("User").Preload("Specialty").Where("user_rut = ?", rut).First(&doctorInfo).Error; err != nil {
		return nil, err
	}
	return &doctorInfo, nil
}

// GetAllDoctorInfo obtiene información de todos los doctores
func (s *DoctorInfoService) GetAllDoctorInfo() ([]models.DoctorInfo, error) {
	var doctorsInfo []models.DoctorInfo
	if err := s.DB.Preload("User").Preload("Specialty").
		Where("is_available = ?", true).
		Order("user_rut ASC").
		Find(&doctorsInfo).Error; err != nil {
		return nil, err
	}
	return doctorsInfo, nil
}

// GetDoctorsBySpecialtyID obtiene todos los doctores de una especialidad
func (s *DoctorInfoService) GetDoctorsBySpecialtyID(specialtyID int) ([]models.DoctorInfo, error) {
	var doctorsInfo []models.DoctorInfo
	if err := s.DB.Preload("User").Preload("Specialty").
		Where("specialty_id = ? AND is_available = ?", specialtyID, true).
		Order("user_rut ASC").
		Find(&doctorsInfo).Error; err != nil {
		return nil, err
	}
	return doctorsInfo, nil
}

// GetDoctorsPaginated obtiene doctores con paginación
func (s *DoctorInfoService) GetDoctorsPaginated(page int, pageSize int, search *string, specialtyID *int) ([]models.DoctorInfo, int64, error) {
	var doctorsInfo []models.DoctorInfo
	var total int64

	query := s.DB.Model(&models.DoctorInfo{}).Where("is_available = ?", true)

	// Filtrar por especialidad si se proporciona
	if specialtyID != nil {
		query = query.Where("specialty_id = ?", *specialtyID)
	}

	// Aplicar búsqueda si se proporciona
	if search != nil && *search != "" {
		query = query.Joins("JOIN \"user\" u ON doctor_info.user_rut = u.rut").
			Where("u.name ILIKE ? OR doctor_info.medical_license ILIKE ? OR doctor_info.specialization ILIKE ?",
				"%"+*search+"%", "%"+*search+"%", "%"+*search+"%")
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Specialty").
		Order("user_rut ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&doctorsInfo).Error; err != nil {
		return nil, 0, err
	}

	return doctorsInfo, total, nil
}

// UpdateDoctorInfo actualiza información de un doctor
func (s *DoctorInfoService) UpdateDoctorInfo(rut string, doctorInfo *models.DoctorInfo) (*models.DoctorInfo, error) {
	var existingDoctorInfo models.DoctorInfo
	if err := s.DB.Where("user_rut = ?", rut).First(&existingDoctorInfo).Error; err != nil {
		return nil, err
	}

	existingDoctorInfo.MedicalLicense = doctorInfo.MedicalLicense
	if doctorInfo.LicenseExpirationDate != nil {
		existingDoctorInfo.LicenseExpirationDate = doctorInfo.LicenseExpirationDate
	}
	existingDoctorInfo.Specialization = doctorInfo.Specialization
	existingDoctorInfo.SpecialtyID = doctorInfo.SpecialtyID
	existingDoctorInfo.YearsOfExperience = doctorInfo.YearsOfExperience
	existingDoctorInfo.Phone = doctorInfo.Phone
	existingDoctorInfo.EmergencyContact = doctorInfo.EmergencyContact
	existingDoctorInfo.EmergencyPhone = doctorInfo.EmergencyPhone
	existingDoctorInfo.IsAvailable = doctorInfo.IsAvailable
	existingDoctorInfo.Notes = doctorInfo.Notes

	if err := s.DB.Save(&existingDoctorInfo).Error; err != nil {
		return nil, err
	}

	return &existingDoctorInfo, nil
}

// DeleteDoctorInfo elimina información de un doctor
func (s *DoctorInfoService) DeleteDoctorInfo(rut string) error {
	return s.DB.Where("user_rut = ?", rut).Delete(&models.DoctorInfo{}).Error
}

// GetDoctorsBySpecialtyCode obtiene todos los doctores de una especialidad por código
func (s *DoctorInfoService) GetDoctorsBySpecialtyCode(specialtyCode string) ([]models.DoctorInfo, error) {
	var doctorsInfo []models.DoctorInfo
	if err := s.DB.Preload("User").Preload("Specialty").
		Joins("JOIN medical_specialty ms ON doctor_info.specialty_id = ms.id").
		Where("ms.code = ? AND doctor_info.is_available = ?", specialtyCode, true).
		Order("user_rut ASC").
		Find(&doctorsInfo).Error; err != nil {
		return nil, err
	}
	return doctorsInfo, nil
}

