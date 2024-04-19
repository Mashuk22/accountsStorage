package model

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                    uuid.UUID `json:"id,omitempty"`
	Name                  string    `json:"name,omitempty"`
	AccountType           string    `json:"account_type,omitempty"`
	Login                 string    `json:"login,omitempty"`
	Password              string    `json:"password,omitempty"`
	Email                 string    `json:"email,omitempty"`
	EmailPassword         string    `json:"emailPassword,omitempty"`
	RecoveryEmail         string    `json:"recovery_email,omitempty"`
	RecoveryEmailPassword string    `json:"recovery_email_password,omitempty"`
	Cookie                string    `json:"cookie,omitempty"`
	Status                string    `json:"status,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
}

type AccountCreate struct {
	Name                  string `json:"name,omitempty"`
	AccountType           string `json:"account_type,omitempty"`
	Login                 string `json:"login,omitempty"`
	Password              string `json:"password,omitempty"`
	Email                 string `json:"email,omitempty"`
	EmailPassword         string `json:"emailPassword,omitempty"`
	RecoveryEmail         string `json:"recovery_email,omitempty"`
	RecoveryEmailPassword string `json:"recovery_email_password,omitempty"`
	Cookie                string `json:"cookie,omitempty"`
	Status                string `json:"status,omitempty"`
}

type AccountUpdate struct {
	Name                  string `json:"name,omitempty"`
	AccountType           string `json:"account_type,omitempty"`
	Login                 string `json:"login,omitempty"`
	Password              string `json:"password,omitempty"`
	Email                 string `json:"email,omitempty"`
	EmailPassword         string `json:"emailPassword,omitempty"`
	RecoveryEmail         string `json:"recovery_email,omitempty"`
	RecoveryEmailPassword string `json:"recovery_email_password,omitempty"`
	Cookie                string `json:"cookie,omitempty"`
	Status                string `json:"status,omitempty"`
}
