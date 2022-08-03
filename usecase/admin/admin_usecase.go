package admin

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type UseCase interface {
	GetAdmin(ctx context.Context, req dto.GetAdminRequest) (dto.GetAdminResponse, error)
	GetAdmins(ctx context.Context, req dto.GetAdminsRequest) (dto.GetAdminsResponse, error)
	LoginAdmin(ctx context.Context, req dto.LoginAdminRequest) (dto.LoginAdminResponse, error)
}

type UseCaseRepo struct {
	repo Repository
}

func New(r Repository) UseCaseRepo {
	return UseCaseRepo{repo: r}
}

func (u UseCaseRepo) LoginAdmin(ctx context.Context, req dto.LoginAdminRequest) (dto.LoginAdminResponse, error) {
	admin, err := u.repo.LoginAdmin(ctx, req.Email, req.Password)
	if err != nil {
		return dto.LoginAdminResponse{}, err
	}
	return dto.LoginAdminResponse{Admin: admin}, err
}

func (u UseCaseRepo) GetAdmin(ctx context.Context, req dto.GetAdminRequest) (dto.GetAdminResponse, error) {
	admin, err := u.repo.GetAdmin(ctx, req.AdminId)
	if err != nil {
		return dto.GetAdminResponse{}, err
	}

	return dto.GetAdminResponse{Admin: admin}, nil
}

func (u UseCaseRepo) GetAdmins(ctx context.Context, _ dto.GetAdminsRequest) (dto.GetAdminsResponse, error) {
	admins, err := u.repo.GetAdmins(ctx)
	if err != nil {
		return dto.GetAdminsResponse{}, err
	}

	return dto.GetAdminsResponse{Admins: admins}, nil
}
