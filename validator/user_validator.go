package validator

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	repository "github.com/XBozorg/bookstore/adapter/repository/mysql"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/user"
)

func doesUserExist(ctx context.Context, repo user.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		userID := value.(string)

		ok, err := repo.DoesUserExist(ctx, userID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("user does not exist")
		}
		return nil
	}
}

func doesPhoneExist(ctx context.Context, repo user.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		phoneID := value.(uint)

		ok, err := repo.DoesPhoneExist(ctx, phoneID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("phone does not exist")
		}
		return nil
	}
}

func doesAddressExist(ctx context.Context, repo user.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		addressID := value.(uint)

		ok, err := repo.DoesAddressExist(ctx, addressID)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("address does not exist")
		}
		return nil
	}
}

func ValidateCreateUser(req dto.CreateUserRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required, is.ASCII, validation.Length(6, 60)),
		validation.Field(&req.Username, validation.Required, is.Alphanumeric, validation.Length(6, 40)),
		validation.Field(&req.FirstName, validation.Required, is.Alpha, validation.Length(1, 80)),
		validation.Field(&req.LastName, validation.Required, is.Alpha, validation.Length(1, 80)),
	)

}

func ValidateGetUser(repo repository.MySQLRepo) user.ValidateGetUser {
	return func(ctx context.Context, req dto.GetUserRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateDeleteUser(repo repository.MySQLRepo) user.ValidateDeleteUser {
	return func(ctx context.Context, req dto.DeleteUserRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateLoginUser(repo repository.MySQLRepo) user.ValidateLoginUser {
	return func(ctx context.Context, req dto.LoginUserRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Email, validation.When(req.Username == "", validation.Required.Error("Either Username or Email is required"), is.Email)),
			validation.Field(&req.Username, validation.When(req.Email == "", validation.Required.Error("Either Username or Email is required"), is.Alphanumeric, validation.Length(6, 40))),
			validation.Field(&req.Password, is.ASCII, validation.Length(6, 60)),
		)
	}
}

func ValidateChangePass(repo repository.MySQLRepo) user.ValidateChangePass {
	return func(ctx context.Context, req dto.ChangePassRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.OldPass, validation.Required, is.ASCII, validation.Length(6, 60)),
			validation.Field(&req.NewPass, validation.Required, is.ASCII, validation.Length(6, 60)),
		)
	}
}

func ValidateChangeUsername(repo repository.MySQLRepo) user.ValidateChangeUsername {
	return func(ctx context.Context, req dto.ChangeUsernameRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.Username, validation.Required, is.Alphanumeric, validation.Length(6, 40)),
		)
	}
}

func ValidateAddPhone(repo repository.MySQLRepo) user.ValidateAddPhone {
	return func(ctx context.Context, req dto.AddPhoneRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.Code, validation.Required, is.Digit, validation.Length(1, 10)),
			validation.Field(&req.PhoneNumber, validation.Required, is.Digit, validation.Length(5, 20)),
		)
	}
}

func ValidateGetPhone(repo repository.MySQLRepo) user.ValidateGetPhone {
	return func(ctx context.Context, req dto.GetPhoneRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.PhoneID, validation.Required, validation.By(doesPhoneExist(ctx, repo))),
		)
	}
}

func ValidateGetPhones(repo repository.MySQLRepo) user.ValidateGetPhones {
	return func(ctx context.Context, req dto.GetPhonesRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateDeletePhone(repo repository.MySQLRepo) user.ValidateDeletePhone {
	return func(ctx context.Context, req dto.DeletePhoneRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.PhoneID, validation.Required, validation.By(doesPhoneExist(ctx, repo))),
		)
	}
}

func ValidateAddAddress(repo repository.MySQLRepo) user.ValidateAddAddress {
	return func(ctx context.Context, req dto.AddAddressRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.Country, validation.Required, is.CountryCode2),
			validation.Field(&req.City, validation.Required, is.ASCII, validation.Length(1, 50)),
			validation.Field(&req.Street, validation.Required, is.ASCII, validation.Length(1, 50)),
			validation.Field(&req.PostalCode, validation.Required, is.Digit, validation.Length(3, 20)),
			validation.Field(&req.No, is.Digit, validation.Length(1, 5)),
			validation.Field(&req.Description, is.Alphanumeric, validation.Length(1, 50)),
		)
	}
}

func ValidateGetAddress(repo repository.MySQLRepo) user.ValidateGetAddress {
	return func(ctx context.Context, req dto.GetAddressRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.AddressID, validation.Required, validation.By(doesAddressExist(ctx, repo))),
		)
	}
}

func ValidateGetAddresses(repo repository.MySQLRepo) user.ValidateGetAddresses {
	return func(ctx context.Context, req dto.GetAddressesRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateDeleteAddress(repo repository.MySQLRepo) user.ValidateDeleteAddress {
	return func(ctx context.Context, req dto.DeleteAddressRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.AddressID, validation.Required, validation.By(doesAddressExist(ctx, repo))),
		)
	}
}
