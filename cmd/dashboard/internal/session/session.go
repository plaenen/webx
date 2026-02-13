package session

import "time"

type UserSession struct {
	UserId            string    `json:"user_id"`
	Email             string    `json:"email"`
	DisplayName       string    `json:"display_name"`
	ControlPlaneRoles []string  `json:"control_plane_roles"`
	ConsentComplete   bool      `json:"consent_complete"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *UserSession) IsLoggedIn() bool {
	return s.Email != ""
}

type SessionStore interface {
	Create(data UserSession) (id string, err error)
	Get(id string) (*UserSession, error)
	Update(id string, data UserSession) error
	Delete(id string) error
	DeleteByEmail(email string) (int, error)
	ListAll() (map[string]UserSession, error)
}
