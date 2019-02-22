package dto

import (
	"time"
)

type BaseDto struct {
	Id        int64             `json:"id"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"auto_now"`

	Extra map[string]interface{} `orm:"-" json:"extra"`

	ExtraJson string `orm:"type(text)" json:"-"`
}