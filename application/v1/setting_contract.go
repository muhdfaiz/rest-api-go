package v1

// SettingServiceInterface is a contract that defines the method needed for Setting Service.
type SettingServiceInterface interface {
	GetAllSettings() []*Setting
}

// SettingRepositoryInterface is a contract that defines the method needed for Setting Repository.
type SettingRepositoryInterface interface {
	GetAll() []*Setting
}
