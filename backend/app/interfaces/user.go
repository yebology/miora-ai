package interfaces

import (
	"miora-ai/app/entities"
	"miora-ai/pkg"
)

// IUserService defines the business logic contract for user operations.
type IUserService interface {
	FindOrCreateFromFirebase(firebaseUID, email, name, avatar string) (*entities.User, *pkg.AppError)
	GetByFirebaseUID(firebaseUID string) (*entities.User, *pkg.AppError)
}

// IUserRepository defines the data access contract for user operations.
type IUserRepository interface {
	FindByFirebaseUID(firebaseUID string) (*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
}
