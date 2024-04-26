package auth

import (
	ssov1 "SSO/contract/gen/go/sso"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(server *grpc.Server) {
	ssov1.RegisterAuthServer(server, &serverAPI{})
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func (s *serverAPI) Login(ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}

func (s *serverAPI) IsAdmin(ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAdmin not implemented")
}
