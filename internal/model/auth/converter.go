package auth

import (
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
)

func ToSignupResponse(user core.User) *authv1.SignupResponse {
	return &authv1.SignupResponse{
		UserId:   int64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}
}

func ToRefreshTokenResponse(accessToken, refreshToken string) *authv1.RefreshTokenResponse {
	return &authv1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func FromLoginRequest(req *authv1.LoginRequest) *core.User {
	return &core.User{
		Email:        req.GetEmail(),
		PasswordHash: req.GetPassword(),
	}
}

func FromSignupRequest(req *authv1.SignupRequest) *core.User {
	return &core.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPassword(),
	}
}

func ToValidateTokenResponse(isValid bool, userID int) *authv1.ValidateTokenResponse {
	return &authv1.ValidateTokenResponse{
		IsValid: isValid,
		UserId:  int64(userID),
	}
}
