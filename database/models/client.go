package models

import "gorm.io/gorm"

type Client struct {
	BaseModel
	VAT   string `json:"vat"`
	Name  string `json:"name"`
	Users []User `json:"users"`
}

func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	return c.BaseModel.BeforeCreate(tx)
}
