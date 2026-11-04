package model

import "time"

// ContactMessage 联系表单提交。首版只入库，不发邮件。
type ContactMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Email     string    `gorm:"size:128;not null" json:"email"`
	Subject   string    `gorm:"size:200" json:"subject"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ContactMessage) TableName() string { return "contact_messages" }
