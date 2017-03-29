package v11

import "time"

type EdmHistory struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	UserGUID  string     `json:"user_guid"`
	Event     string     `json:"event"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
