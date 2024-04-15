package service

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/trace"
	"app/pkg/utils"
	"app/repository"
	"context"
	"math"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserService struct {
	OtpSvc    IOtpSvc
	UserRepo  IUserRepo
	TokenRepo ITokenRepo
}

func NewUserService(otpSvc IOtpSvc, userRepo IUserRepo, tokenRepo ITokenRepo) *UserService {
	return &UserService{OtpSvc: otpSvc, UserRepo: userRepo, TokenRepo: tokenRepo}
}

func (s *UserService) Login(ctx context.Context, user *entity.User) (accessToken string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	dbUser, err := s.UserRepo.GetUserByUserNameOrEmail(ctx, user.Username, user.Email)
	if err != nil {
		return
	}
	if !dbUser.IsActive {
		return "", errors.UserInactive()
	}
	if !utils.VerifyPassword(user.Password, dbUser.Password) {
		err = errors.PasswordIncorrect()
		return
	}
	accessToken, err = s.CreateToken(ctx, dbUser)
	if err != nil {
		return
	}
	dbUser.LoggedIn()
	_ = s.UserRepo.SaveUser(ctx, dbUser)
	return
}

func (s *UserService) CreateToken(ctx context.Context, user *entity.User) (accessToken string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	duration := time.Duration(int32(config.Cfg.HTTP.AccessTokenDuration)) * time.Minute
	appToken := &entity.Token{
		ID:       utils.NewID(),
		UserID:   user.ID,
		Name:     user.Name,
		UserName: user.Username,
		Email:    user.Email,
		Type:     entity.AccessTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	err = s.TokenRepo.CreateToken(ctx, appToken)
	if err != nil {
		return
	}
	accessToken = appToken.SignedToken(config.Cfg.HTTP.Secret)
	return
}

func (s *UserService) UserLogout(ctx context.Context, tokenStr string) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	token, err := entity.NewTokenFromEncoded(tokenStr, config.Cfg.HTTP.Secret)
	if err != nil {
		err = apperror.NewError(errors.CodeTokenError, err.Error())
		return
	}
	dbToken, err := s.TokenRepo.GetTokenById(ctx, token.ID)
	if err != nil {
		return
	}
	err = s.TokenRepo.DeleteToken(ctx, dbToken)
	return
}

func (s *UserService) Authenticate(ctx context.Context, tokenStr string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	token, err := entity.NewTokenFromEncoded(tokenStr, config.Cfg.HTTP.Secret)
	if err != nil {
		err = apperror.NewError(errors.CodeTokenError, err.Error())
		return
	}
	dbToken, err := s.TokenRepo.GetTokenById(ctx, token.ID)
	if err != nil {
		return
	}
	dbUser, err := s.UserRepo.GetUserById(ctx, dbToken.UserID)
	if err != nil {
		return
	}
	return dbUser, nil
}

func (s *UserService) GetUser(ctx context.Context, e *entity.User) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.UserRepo.GetUserById(ctx, e.ID)
}

func (s *UserService) GetUserList(ctx context.Context, query *repository.QueryParams) (res []*entity.User, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.UserRepo.GetUserList(ctx, query)
}

func (s *UserService) CreateUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	err = s.UserRepo.CheckUserNameAndEmailExist(ctx, e.Username, e.Email)
	if err != nil {
		return
	}

	user := &entity.User{
		ID: utils.NewID(),
	}
	err = user.OnUserCreated(ctx, e, time.Now())
	if err != nil {
		return
	}
	err = s.UserRepo.SaveUser(ctx, user)
	if err != nil {
		return
	}
	_, err = s.OtpSvc.GenerateOtp(ctx, user.Email)
	return
}

func (s *UserService) ActiveUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	otp, err := s.OtpSvc.VerifyOtp(ctx, &entity.Otp{
		Email: e.Email,
		Code:  e.Otp,
	})
	if err != nil {
		return
	}
	dbUser, err := s.UserRepo.GetInactiveUser(ctx, e.Email)
	if err != nil {
		return
	}
	dbUser.OnUserActive(ctx)
	err = s.UserRepo.UpdateUser(ctx, dbUser)
	if err != nil {
		return
	}

	err = s.OtpSvc.DeleteOtp(ctx, otp)
	return
}

func (s *UserService) ResetPassword(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	otp, err := s.OtpSvc.VerifyOtp(ctx, &entity.Otp{
		Email: e.Email,
		Code:  e.Otp,
	})
	if err != nil {
		return
	}
	dbUser, err := s.UserRepo.GetUserByEmail(ctx, e.Email)
	if err != nil {
		return
	}
	dbUser.OnUserUpdated(ctx, e, time.Now())
	err = s.UserRepo.UpdateUser(ctx, dbUser)
	if err != nil {
		return
	}

	err = s.OtpSvc.DeleteOtp(ctx, otp)
	return
}

func (s *UserService) UpdateUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	dbUser, err := s.UserRepo.GetUserById(ctx, e.ID)
	if err != nil {
		return
	}
	if e.GetUserName() != "" || e.GetEmail() != "" {
		err = s.UserRepo.CheckDuplicateUserNameAndEmail(ctx, dbUser, e.GetUserName(), e.GetEmail())
		if err != nil {
			return nil, err
		}
	}

	err = dbUser.OnUserUpdated(ctx, e, time.Now())
	if err != nil {
		return
	}
	err = s.UserRepo.UpdateUser(ctx, dbUser)
	return
}

func (s *UserService) DeleteUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	dbUser, err := s.UserRepo.GetUserById(ctx, e.ID)
	if err != nil {
		return
	}
	err = dbUser.OnUserDeleted(ctx, e, time.Now())
	if err != nil {
		return
	}
	err = s.UserRepo.DeleteUser(ctx, dbUser)
	return
}

func (s *UserService) SendFriendRequest(ctx context.Context, user *entity.User, friend *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	_, err = s.UserRepo.GetUserById(ctx, user.ID)
	if err != nil {
		return
	}
	_, err = s.UserRepo.GetUserById(ctx, friend.ID)
	if err != nil {
		return
	}

	err = s.UserRepo.AddFriendRequest(ctx, user, friend)
	return
}

func (s *UserService) RejectFriendRequest(ctx context.Context, user *entity.User, friend *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	_, err = s.UserRepo.GetUserById(ctx, user.ID)
	if err != nil {
		return
	}
	_, err = s.UserRepo.GetUserById(ctx, friend.ID)
	if err != nil {
		return
	}

	err = s.UserRepo.RemoveFriendRequest(ctx, user, friend)
	return
}

func (s *UserService) AcceptFriendRequest(ctx context.Context, user *entity.User, friend *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	_, err = s.UserRepo.GetUserById(ctx, user.ID)
	if err != nil {
		return
	}
	_, err = s.UserRepo.GetUserById(ctx, friend.ID)
	if err != nil {
		return
	}

	res, err = s.UserRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		err = s.UserRepo.RemoveFriendRequest(ctx, user, friend)
		if err != nil {
			return
		}
		err = s.UserRepo.AddFriend(ctx, user, friend)
		if err != nil {
			return
		}
		err = s.UserRepo.AddFriend(ctx, friend, user)
		if err != nil {
			return
		}
		return
	})
	return
}

func (s *UserService) RemoveFriend(ctx context.Context, user *entity.User, friend *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	_, err = s.UserRepo.GetUserById(ctx, user.ID)
	if err != nil {
		return
	}
	_, err = s.UserRepo.GetUserById(ctx, friend.ID)
	if err != nil {
		return
	}

	res, err = s.UserRepo.ExecTransaction(ctx, func(ctx context.Context) (res any, err error) {
		err = s.UserRepo.RemoveFriend(ctx, user, friend)
		if err != nil {
			return
		}
		err = s.UserRepo.RemoveFriend(ctx, friend, user)
		if err != nil {
			return
		}
		return
	})
	return
}

func (s *UserService) SuggestFriend(ctx context.Context, e *entity.User) (res []*entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	targetUser, err := s.UserRepo.GetUserById(ctx, e.ID)
	if err != nil {
		return
	}

	allUsers, err := s.UserRepo.GetAllUsers(ctx)
	if err != nil {
		return
	}

	// Find nearest neighbors of the target user
	nearestNeighbors := s.FindNearestNeighbors(targetUser, allUsers)

	// Suggest friends of nearest neighbors
	for _, neighbor := range nearestNeighbors {
		for _, friend := range neighbor.Friends {
			if !s.Contains(targetUser.FriendIds, friend.ID) && friend.ID != targetUser.ID {
				res = append(res, friend)
			}
		}
	}
	return
}

// Function to find nearest neighbors based on cosine similarity
func (s *UserService) FindNearestNeighbors(targetUser *entity.User, allUsers []*entity.User) []*entity.User {
	// Calculate similarity with all users
	var similarities []float64
	for _, user := range allUsers {
		similarity := s.CosineSimilarity(targetUser.FriendIds, user.FriendIds)
		similarities = append(similarities, similarity)
	}

	// Find nearest neighbors
	var nearestNeighbors []*entity.User
	for i, similarity := range similarities {
		if similarity > 0.8 && allUsers[i].ID != targetUser.ID { // Adjust threshold as needed
			nearestNeighbors = append(nearestNeighbors, allUsers[i])
		}
	}

	return nearestNeighbors
}

// Function to calculate cosine similarity between two vectors
func (s *UserService) CosineSimilarity(a, b []string) float64 {
	// Convert friend lists to binary vectors
	vectorA := make([]int, len(a))
	vectorB := make([]int, len(b))
	for i, friendID := range a {
		if s.Contains(b, friendID) {
			vectorA[i] = 1
		}
	}
	for i, friendID := range b {
		if s.Contains(a, friendID) {
			vectorB[i] = 1
		}
	}

	// Calculate dot product
	dotProduct := 0.0
	for i := 0; i < len(vectorA); i++ {
		dotProduct += float64(vectorA[i] * vectorB[i])
	}

	// Calculate magnitudes
	magnitudeA := 0.0
	magnitudeB := 0.0
	for _, val := range vectorA {
		magnitudeA += float64(val * val)
	}
	for _, val := range vectorB {
		magnitudeB += float64(val * val)
	}
	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)

	// Calculate cosine similarity
	similarity := dotProduct / (magnitudeA * magnitudeB)
	return similarity
}

// Function to check if a user ID exists in a list of friend IDs
func (s *UserService) Contains(friends []string, userID string) bool {
	for _, friendID := range friends {
		if friendID == userID {
			return true
		}
	}
	return false
}
