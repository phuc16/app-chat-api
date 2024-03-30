package service

type UserService struct {
	UserRepo IUserRepo
}

func NewUserService(userRepo IUserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}
