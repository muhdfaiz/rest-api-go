package v1_1

// SettingServiceInterface is a contract that defines the method needed for Setting Service.
type SettingServiceInterface interface {
	GetAllSettings() []*Setting
	GetSettingBySlug(slug string) *Setting
}

// SettingRepositoryInterface is a contract that defines the method needed for Setting Repository.
type SettingRepositoryInterface interface {
	GetAll() []*Setting
	GetBySlug(slug string) *Setting
}
