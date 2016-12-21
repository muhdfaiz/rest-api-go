package v1

type SettingService struct {
	SettingRepository SettingRepositoryInterface
}

// GetAllSettings function used to retrieve all settings from database through
// Setting Repository.
func (ss *SettingService) GetAllSettings() []*Setting {
	settings := ss.SettingRepository.GetAll()

	return settings
}
