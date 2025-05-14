package handler

import (
	"context"
	"regexp"
	"time"

	"user-service/internal/model"
	pb "user-service/internal/pb"
	"user-service/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[A-Za-z]{2,}$`)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	uc        *usecase.UserUsecase
	jwtSecret string // Add field for JWT secret
}

func NewUserHandler(uc *usecase.UserUsecase, jwtSecret string) *UserHandler {
	return &UserHandler{
		uc:        uc,
		jwtSecret: jwtSecret,
	}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if !emailRegex.MatchString(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}
	if len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 6 characters")
	}

	user := &model.User{Username: req.Username, Password: req.Password, Email: req.Email}
	id, err := h.uc.CreateUser(ctx, user)
	if err != nil {
		switch err.Error() {
		case "username already exists":
			return nil, status.Error(codes.AlreadyExists, "username already taken")
		case "failed to hash password":
			return nil, status.Error(codes.Internal, "could not secure password")
		default:
			return nil, status.Errorf(codes.Internal, "registration error: %v", err)
		}
	}
	return &pb.UserResponse{Id: id, Message: "User registered"}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	user, err := h.uc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate a real JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,                               // Subject (user ID)
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
		"iat": time.Now().Unix(),                     // Issued at
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.AuthResponse{Token: tokenString, Message: "Authenticated"}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.UserID) (*pb.UserProfile, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Используем GetUserByID с кэшированием
	user, err := h.uc.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.UserProfile{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
