package repositories

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {

	return &UserRepository{db: db}

}

func (r *UserRepository) FindByFirebaseUID(firebaseUID string) (*entities.User, error) {

	var user entities.User
	if err := r.db.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *UserRepository) Create(user *entities.User) error {

	return r.db.Create(user).Error

}

func (r *UserRepository) Update(user *entities.User) error {

	return r.db.Save(user).Error

}
