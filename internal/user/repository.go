package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.First(&user, "email = ?", email).Error
	return &user, err
}
