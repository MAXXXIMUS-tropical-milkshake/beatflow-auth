package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	helper "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/grpc"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/auth"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/validator"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	authv1.UnimplementedAuthServiceServer
	authService core.AuthService
	authConfig  core.AuthConfig
}

func Register(gRPCServer *grpc.Server, authService core.AuthService, authConfig core.AuthConfig) {
	authv1.RegisterAuthServiceServer(gRPCServer, &server{authService: authService, authConfig: authConfig})
}

func (s *server) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := s.authService.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrRefreshTokenNotValid) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return auth.ToRefreshTokenResponse(*accessToken, *refreshToken), nil
}

func (s *server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	v := validator.New()
	if model.ValidateLoginRequest(v, req); !v.Valid() {
		logger.Log().Debug(ctx, fmt.Sprintf("%+v", v.Errors))
		return nil, helper.ToGRPCError(v)
	}

	user := auth.FromLoginRequest(req)

	accessToken, refreshToken, err := s.authService.Login(ctx, *user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrInvalidCredentials) || errors.Is(err, core.ErrUserNotFound) || errors.Is(err, core.ErrAlreadyDeleted) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return &authv1.LoginResponse{AccessToken: *accessToken, RefreshToken: *refreshToken}, nil
}

func (s *server) Signup(ctx context.Context, req *authv1.SignupRequest) (*authv1.SignupResponse, error) {
	v := validator.New()
	if model.ValidateSignupRequest(v, req); !v.Valid() {
		logger.Log().Debug(ctx, fmt.Sprintf("%+v", v.Errors))
		return nil, helper.ToGRPCError(v)
	}

	user := auth.FromSignupRequest(req)

	retUser, err := s.authService.Signup(ctx, *user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrEmailAlreadyExists) || errors.Is(err, core.ErrUsernameAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return auth.ToSignupResponse(*retUser), nil
}

func (s *server) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	userID, err := helper.ValidToken(ctx, req.GetToken(), s.authConfig.Secret)
	if err != nil {
		if errors.Is(err, core.ErrUnauthorized) {
			return auth.ToValidateTokenResponse(false, 0), nil
		}
		logger.Log().Error(ctx, err.Error())
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return auth.ToValidateTokenResponse(true, *userID), nil
}
