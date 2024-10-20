package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	helper "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/grpc"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model"
	usermodel "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/user"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/validator"
	userv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	userv1.UnimplementedUserServiceServer
	userService core.UserService
}

func Register(gRPCServer *grpc.Server, userService core.UserService) {
	userv1.RegisterUserServiceServer(gRPCServer, &server{userService: userService})
}

func (s *server) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	v := validator.New()
	model.ValidateUpdateUserRequest(v, req)
	if !v.Valid() {
		logger.Log().Debug(ctx, fmt.Sprintf("%+v", v.Errors))
		return nil, helper.ToGRPCError(v)
	}

	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, status.Error(codes.Unauthenticated, core.ErrInvalidCredentials.Error())
	}

	user := usermodel.FromUpdateUserRequest(req, userID)

	retUser, err := s.userService.UpdateUser(ctx, *user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, core.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		} else if errors.Is(err, core.ErrEmailAlreadyExists) || errors.Is(err, core.ErrUsernameAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		} else if errors.Is(err, core.ErrAlreadyDeleted) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return usermodel.ToUpdateUserResponse(*retUser), nil
}

func (s *server) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, status.Error(codes.Unauthenticated, core.ErrUnauthorized.Error())
	}

	err = s.userService.DeleteUser(ctx, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return &userv1.DeleteUserResponse{}, nil
}

func (s *server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	v := validator.New()
	model.ValidateGetUserRequest(v, req)
	if !v.Valid() {
		logger.Log().Debug(ctx, fmt.Sprintf("%+v", v.Errors))
		return nil, helper.ToGRPCError(v)
	}

	getUser := usermodel.FromGetUserRequest(req)

	user, err := s.userService.GetUser(ctx, *getUser)
	if err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		logger.Log().Error(ctx, err.Error())
		return nil, status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return usermodel.ToGetUserResponse(*user), nil
}
