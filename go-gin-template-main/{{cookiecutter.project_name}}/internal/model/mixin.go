package model

import (
	"time"
)

type DeleteMixin struct {
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

type TimestampMixin struct {
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

//func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
//	if u.Role == "admin" {
//		return errors.New("admin user not allowed to update")
//	}
//	return
//}
//func (m *BaseModel) BeforeInsert(db gorm.DB) error {
//	now := time.Now()
//	if m.CreatedAt.IsZero() {
//		m.CreatedAt = now
//	}
//	if m.UpdatedAt.IsZero() {
//		m.UpdatedAt = now
//	}
//	return nil
//}
//
//func (m *BaseModel) BeforeUpdate(db gorm.DB) error {
//	m.UpdatedAt = time.Now()
//	return nil
//}
