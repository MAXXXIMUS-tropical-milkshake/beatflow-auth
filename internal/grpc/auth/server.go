package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	helper "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/grpc"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/request"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/response"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/validator"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	authv1.UnimplementedAuthServiceServer
	auth core.AuthService
}

func Register(gRPCServer *grpc.Server, auth core.AuthService) {
	authv1.RegisterAuthServiceServer(gRPCServer, &server{auth: auth})
}

func (s *server) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := s.auth.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrRefreshTokenNotValid) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return response.ToRefreshTokenResponse(*accessToken, *refreshToken), nil
}

func (s *server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	v := validator.New()
	if request.ValidateLoginRequest(v, req); !v.Valid() {
		logger.Log().Debug(ctx, fmt.Sprintf("%+v", v.Errors))
		return nil, helper.ToGRPCError(v)
	}

	user := request.FromLoginRequest(req)

	accessToken, refreshToken, err := s.auth.Login(ctx, *user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return &authv1.LoginResponse{AccessToken: *accessToken, RefreshToken: *refreshToken}, nil
}

func (s *server) Signup(ctx context.Context, req *authv1.SignupRequest) (*authv1.SignupResponse, error) {
	v := validator.New()
	if request.ValidateSignupRequest(v, req); !v.Valid() {
		logger.Log().Debug(ctx, fmt.Sprintf("%+v", v.Errors))
		return nil, helper.ToGRPCError(v)
	}

	user := request.FromSignupRequest(req)

	retUser, err := s.auth.Signup(ctx, *user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return response.ToSignupResponse(*retUser), nil
}

func (s *server) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	v := validator.New()
	if request.ValidateUpdateUserRequest(v, req); !v.Valid() {
		logger.Log().Debug(ctx, fmt.Sprintf("%+v", v.Errors))
		return nil, helper.ToGRPCError(v)
	}

	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, status.Error(codes.Unauthenticated, core.ErrInvalidCredentials.Error())
	}

	user := request.FromUpdateUserRequest(req, userID)

	retUser, err := s.auth.UpdateUser(ctx, *user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else if errors.Is(err, core.ErrEmailAlreadyExists) || errors.Is(err, core.ErrUsernameAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return response.ToUpdateUserResponse(*retUser), nil
}
