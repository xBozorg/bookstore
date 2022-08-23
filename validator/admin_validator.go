package validator

import (
	"context"
	"errors"

	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/admin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func doesAdminExist(ctx context.Context, repo admin.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		adminID := value.(string)

		ok, err := repo.DoesAdminExist(ctx, adminID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("admin does not exist")
		}
		return nil
	}
}

func ValidateGetAdmin(storage repository.Storage) admin.ValidateGetAdmin {
	return func(ctx context.Context, req dto.GetAdminRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.AdminId, is.UUIDv4, validation.By(doesAdminExist(ctx, storage))),
		)
	}
}

func ValidateLoginAdmin(storage repository.Storage) admin.ValidateLoginAdmin {
	return func(ctx context.Context, req dto.LoginAdminRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Email, validation.Required, is.Email),
			validation.Field(&req.Password, is.ASCII, validation.Length(6, 60)),
		)
	}
}
