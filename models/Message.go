package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID                   uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DirectConversationID *uuid.UUID `gorm:"type:uuid;index" json:"direct_conversation_id,omitempty"`
	SenderID             uuid.UUID  `gorm:"type:uuid;not null" json:"sender_id"`
	Content              string     `gorm:"type:text;not null" json:"content"`
	CreatedAt            time.Time

	DirectConversation *DirectConversation `gorm:"foreignKey:DirectConversationID"`
	Sender             *User               `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
