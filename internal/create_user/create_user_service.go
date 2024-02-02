package create_user

type CreateUserService struct{ userRepository CreateUserRepositoryInterface }

func NewUserService(userRepository CreateUserRepositoryInterface) *CreateUserService {
	return &CreateUserService{userRepository}
}

type CreateUserServiceInterface interface {
	CreateUser(payload DTOCreateUserReq) error
}

func (ur *CreateUserService) CreateUser(payload DTOCreateUserReq) error {
	return nil
}
