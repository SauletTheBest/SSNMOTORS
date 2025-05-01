package handler

import (
    "context"
    "user-service/internal/pb"
    "user-service/internal/usecase"
    "user-service/internal/token"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type UserHandler struct {
    pb.UnimplementedUserServiceServer
    uc *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
    return &UserHandler{uc: uc}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
    user, err := h.uc.Register(ctx, req.Email, req.Password, req.Name)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    return &pb.UserResponse{
        Id:    user.ID,
        Email: user.Email,
        Name:  user.Name,
        Role:  user.Role,
    }, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
    user, err := h.uc.Authenticate(ctx, req.Email, req.Password)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, err.Error())
    }
    token, err := token.GenerateToken(user.ID)
    if err != nil {
        return nil, status.Error(codes.Internal, "could not generate token")
    }
    return &pb.AuthResponse{Token: token}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
    user, err := h.uc.GetProfile(ctx, req.UserId)
    if err != nil {
        return nil, status.Error(codes.NotFound, err.Error())
    }
    return &pb.UserResponse{
        Id:    user.ID,
        Email: user.Email,
        Name:  user.Name,
        Role:  user.Role,
    }, nil
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
    user, err := h.uc.UpdateProfile(ctx, req.UserId, req.Name, req.Email)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    return &pb.UpdateResponse{Id: user.ID, Message: "Profile updated"}, nil
}
