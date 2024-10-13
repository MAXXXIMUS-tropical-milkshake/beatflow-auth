package auth

import (
	"context"
	"errors"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	authv1.UnimplementedAuthServer
	auth core.AuthService
}

func Register(gRPCServer *grpc.Server, auth core.AuthService) {
	authv1.RegisterAuthServer(gRPCServer, &server{auth: auth})
}

func (s *server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user := core.User{
		Username:     req.GetUsername(),
		PasswordHash: req.GetPassword(),
	}

	token, err := s.auth.Login(ctx, user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to login")
	}

	return &authv1.LoginResponse{Token: *token}, nil
}

func (s *server) Signup(ctx context.Context, req *authv1.SignupRequest) (*authv1.SignupResponse, error) {
	user := core.User{
		Username:     req.GetUsername(),
		PasswordHash: req.GetPassword(),
	}

	err := s.auth.Signup(ctx, user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to signup")
	}

	return &authv1.SignupResponse{}, nil
}

func (s *server) UpdatePassword(ctx context.Context, req *authv1.UpdatePasswordRequest) (*authv1.UpdatePasswordResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	user := core.User{
		ID:           userID,
		PasswordHash: req.GetNewPassword(),
	}

	err = s.auth.UpdatePassword(ctx, user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update password")
	}

	return &authv1.UpdatePasswordResponse{}, nil
}
