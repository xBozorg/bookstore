package admin

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type (
	ValidateGetAdmin   func(ctx context.Context, req dto.GetAdminRequest) error
	ValidateGetAdmins  func(ctx context.Context, req dto.GetAdminsRequest) error
	ValidateLoginAdmin func(ctx context.Context, req dto.LoginAdminRequest) error
)
