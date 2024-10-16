package response

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

func ToUpdateUserResponse(user core.User) *authv1.UpdateUserResponse {
	return &authv1.UpdateUserResponse{
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
