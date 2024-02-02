package create_user

type CreateUserService struct{ userRepository *UserRepository }

func NewUserService(userRepository *UserRepository) *CreateUserService {
	return &CreateUserService{userRepository}
}

func (ur *CreateUserService) CreateUserService(payload DTOCreateUserReq) error {
	return nil
}
