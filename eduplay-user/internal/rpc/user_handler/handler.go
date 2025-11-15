package sign_up_user

import (
	"context"
	"eduplay-user/internal/model"
	"eduplay-user/internal/rpc/converters"
	"eduplay-user/internal/storage"
	"errors"
	"fmt"

	dto "eduplay-user/internal/generated"

	"google.golang.org/grpc"
)

type UseCase interface {
	SignUpUser(ctx context.Context, name string, surname string, email string, organization string, phone string, password string) (*model.Session, error)
	SignInUser(ctx context.Context, email string, password string) (*model.Session, error)
	RefreshSession(ctx context.Context, refreshToken string) (*model.Tokens, error)
	GetUserAccess(ctx context.Context, authToken string) (*dto.GetUserAccessOut, error)
	GetUserInfo(ctx context.Context, accessToken string) (*model.UserInfo, error)
	ChangeUserInfo(ctx context.Context, in *dto.ChangeUserInfoIn) (*model.UserInfo, error)
	ChangeUserPassword(ctx context.Context, in *dto.ChangePasswordIn) error
	DeleteUserAccount(ctx context.Context, accessToken string) error
	SignOutUser(ctx context.Context, accessToken string) error
}

type Handler struct {
	dto.UnimplementedUsersServer
	uc UseCase
}

func Register(gRPCServer *grpc.Server, uc UseCase) {
	dto.RegisterUsersServer(gRPCServer, &Handler{uc: uc})
}

func (h *Handler) SignUp(
	ctx context.Context,
	in *dto.SignUpIn,
) (*dto.SignUpOut, error) {
	op := "SignUpUser.Handler"

	session, err := h.uc.SignUpUser(ctx, in.Name, in.Surname, in.Email, in.Organization, in.Phone, in.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return &dto.SignUpOut{ErrorMessage: storage.ErrUserAlreadyExists.Error()}, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &dto.SignUpOut{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		Role:         converters.StringToDto(session.Role),
		AccessLevel:  int64(session.AccessLevel)}, nil
}

func (h *Handler) SignIn(ctx context.Context, in *dto.SignInIn) (*dto.SignUpOut, error) {
	op := "SignInUser.Handler"

	session, err := h.uc.SignInUser(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return &dto.SignUpOut{ErrorMessage: storage.ErrUserNotFound.Error()}, nil
		}
		if errors.Is(err, storage.ErrIncorrectPassword) {
			return &dto.SignUpOut{ErrorMessage: storage.ErrIncorrectPassword.Error()}, nil
		}
		if errors.Is(err, storage.ErrIsActive) {
			return &dto.SignUpOut{ErrorMessage: storage.ErrIsActive.Error()}, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &dto.SignUpOut{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		Role:         converters.StringToDto(session.Role),
		AccessLevel:  int64(session.AccessLevel)}, nil
}

func (h *Handler) Refresh(ctx context.Context, in *dto.RefreshIn) (*dto.RefreshOut, error) {
	op := "RefreshSession.Handler"

	tokens, err := h.uc.RefreshSession(ctx, in.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if tokens.Message != "" {
		return &dto.RefreshOut{Message: tokens.Message}, nil
	}

	return &dto.RefreshOut{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (h *Handler) GetUserAccess(ctx context.Context, in *dto.GetUserAccessIn) (*dto.GetUserAccessOut, error) {
	op := "GetUserAccess.Handler"

	userAccess, err := h.uc.GetUserAccess(ctx, in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userAccess, nil
}

func (h *Handler) GetUserInfo(ctx context.Context, in *dto.GetUserAccessIn) (*dto.GetUserInfoOut, error) {
	op := "GetUserInfo.Handler"

	info, err := h.uc.GetUserInfo(ctx, in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp := &dto.GetUserInfoOut{
		Name:                   info.Name,
		Surname:                info.Surname,
		JobTitle:               info.JobTitle,
		Organisation:           info.Organisation,
		Phone:                  info.Phone,
		Email:                  info.Email,
		City:                   info.City,
		ShortOrganisationTitle: info.ShortOrgTitle,
		INN:                    info.INN,
		OrganisationType:       info.OrganisationType,
		CurrentTarrif:          info.CurrentTarrif,
	}

	return resp, nil
}

func (h *Handler) ChangeUserInfo(ctx context.Context, in *dto.ChangeUserInfoIn) (*dto.GetUserInfoOut, error) {
	op := "ChangeUserInfo.Handler"

	info, err := h.uc.ChangeUserInfo(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp := &dto.GetUserInfoOut{
		Name:                   info.Name,
		Surname:                info.Surname,
		JobTitle:               info.JobTitle,
		Organisation:           info.Organisation,
		Phone:                  info.Phone,
		Email:                  info.Email,
		City:                   info.City,
		ShortOrganisationTitle: info.ShortOrgTitle,
		INN:                    info.INN,
		OrganisationType:       info.OrganisationType,
		CurrentTarrif:          info.CurrentTarrif,
	}

	return resp, nil
}

func (h *Handler) ChangePassword(ctx context.Context, in *dto.ChangePasswordIn) (*dto.Empty, error) {
	op := "ChangePassword.Handler"

	err := h.uc.ChangeUserPassword(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.Empty{}, nil
}

func (h *Handler) DeleteAccount(ctx context.Context, in *dto.DeleteAccountIn) (*dto.Empty, error) {
	op := "ChangePassword.Handler"

	err := h.uc.DeleteUserAccount(ctx, in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.Empty{}, nil
}

func (h *Handler) SignOut(ctx context.Context, in *dto.DeleteAccountIn) (*dto.Empty, error) {
	op := "SignOutUser.Handler"

	err := h.uc.SignOutUser(ctx, in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.Empty{}, nil
}
