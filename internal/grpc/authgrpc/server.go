package authgrpc

import (
	ssov1 "SSO/contract/gen/go/sso"
	"SSO/internal/services/auth"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int32,
	) (token string, err error)
	RegisterUser(ctx context.Context,
		email string,
		password string,
	) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

const (
	emptyValue = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(server *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(server, &serverAPI{auth: auth})
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "user hasn't registered")
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	// if validation pass successful
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "user hasn't login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")

		}

		return nil, status.Error(codes.Internal, "internal msg")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

// validateLogin checks for user login input correctly
func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

// validateRegister checks for user register input correctly
func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}
