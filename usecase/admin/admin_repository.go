package admin

import (
	"context"

	"github.com/XBozorg/bookstore/entity/admin"
)

type Repository interface {
	GetAdmin(ctx context.Context, adminID string) (admin.Admin, error)
	GetAdmins(ctx context.Context) ([]admin.Admin, error)
	LoginAdmin(ctx context.Context, email, password string) (admin.Admin, error)
}

type ValidatorRepo interface {
	DoesAdminExist(ctx context.Context, adminID string) (bool, error)
}
