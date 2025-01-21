package session

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	utils "github.com/igortoigildin/goph-keeper/pkg/utils"
)

const sessionFile = "session.json"

type Session struct {
	Email string `json:"email"`
	Token string `json:"token"`
	ExpiresAt time.Time	`json:"expires_at"`
}

func IsSessionValid(tokenSectet string) bool {
	session, err := LoadSession()
	if err != nil {
		return false
	}

	_, err = utils.VeryfyToken(session.Token, []byte(tokenSectet))
	if err != nil {
		return false
	}

	return time.Now().Before(session.ExpiresAt)
}

func LoadSession() (*Session, error) {
	file, err := os.Open(sessionFile)
	if err != nil {
		return nil, fmt.Errorf("session file does not exist of could not be opened: %w", err)
	}
	defer file.Close()

	var session Session
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("could not decode session data: %w", err)
	}

	return &session, nil
}

func SaveSession(session *Session) error {
	file, err := os.Create(sessionFile)
	if err != nil {
		return fmt.Errorf("could not create session file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return fmt.Errorf("could not encode session data: %w", err)
	}

	return nil
}
