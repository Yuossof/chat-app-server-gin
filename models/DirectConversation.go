package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type DirectConversation struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	User1ID   uuid.UUID `gorm:"type:uuid;not null"`
	User2ID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []Message `gorm:"foreignKey:DirectConversationID"`

	User1 User `gorm:"foreignKey:User1ID"`
	User2 User `gorm:"foreignKey:User2ID"`
}

func (dc DirectConversation) Validate() error {
	return validation.ValidateStruct(&dc,
		validation.Field(&dc.User1ID, validation.Required),
		validation.Field(&dc.User2ID, validation.Required),
	)
}
