package services

import (
	"errors"

	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/constants"
	"miora-ai/pkg"

	"gorm.io/gorm"
)

type UserService struct {
	repo interfaces.IUserRepository
}

func NewUserService(repo interfaces.IUserRepository) *UserService {

	return &UserService{repo: repo}

}

// FindOrCreateFromFirebase finds a user by Firebase UID, or creates one if not found.
// Also updates name/avatar if they changed.
func (s *UserService) FindOrCreateFromFirebase(firebaseUID, email, name, avatar string) (*entities.User, *pkg.AppError) {

	user, err := s.repo.FindByFirebaseUID(firebaseUID)
	if err == nil {
		// Update name/avatar if changed
		if user.Name != name || user.Avatar != avatar {
			user.Name = name
			user.Avatar = avatar
			s.repo.Update(user)
		}
		return user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrInternal()
	}

	user = &entities.User{
		FirebaseUID: firebaseUID,
		Email:       email,
		Name:        name,
		Avatar:      avatar,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, pkg.ErrInternal()
	}

	return user, nil

}

func (s *UserService) GetByFirebaseUID(firebaseUID string) (*entities.User, *pkg.AppError) {

	user, err := s.repo.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}
	return user, nil

}
