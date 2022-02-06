package money

import (
	"net/http"

	"github.com/google/uuid"
)

type SessionID = uuid.UUID

type Session struct {
	Id      SessionID `json:"id"`
	Credits int64     `json:"credits"`
}

func (s *Session) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
