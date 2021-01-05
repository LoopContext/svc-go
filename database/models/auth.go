package models

import (
	"time"

	"github.com/loopcontext/svc-go/utils"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Email    string       `gorm:"uniqueIndex" json:"email"`
	UserName string       `gorm:"uniqueIndex" json:"userName"`
	Password string       `json:"-"`
	FullName string       `json:"fullName"`
	APIKeys  []UserAPIKey `json:"apiKeys"`
	ClientID string       `json:"-"`
	Client   Client       `json:"client"`
}

type UserAPIKey struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	APIKey    string `gorm:"unique" json:"apiKey"`
	UserID    string `json:"-"`
	User      User   `json:"-"`
}

// BeforeCreate hook that runs before entity create
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	err = u.BaseModel.BeforeCreate(tx)
	if err != nil {
		return err
	}
	return u.hashPassword()
}

// BeforeUpdate hook to run before entity update
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return u.hashPassword()
}

// BeforeCreate hook to run before entity create
func (uak *UserAPIKey) BeforeCreate(tx *gorm.DB) (err error) {
	uak.APIKey, err = utils.GenerateGUID()
	return
}

func (u *User) hashPassword() (err error) {
	if u.Password != "" {
		hash, err := utils.EncryptPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
	}

	return err
}
