package v1

import "time"

type Event struct {
	ID        int        `json:"id"`
	GUID      string     `json:"guid"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	Img       string     `json:"img"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
