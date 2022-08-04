package user

import (
	"context"

	"github.com/XBozorg/bookstore/entity/user"
)

type Repository interface {
	CreateUser(ctx context.Context, user user.User) (user.User, error)
	GetUser(ctx context.Context, userID string) (user.User, error)
	GetUsers(ctx context.Context) ([]user.User, error)
	DeleteUser(ctx context.Context, userID string) error

	LoginUser(ctx context.Context, username, email, password string) (user.User, error)

	ChangePassword(ctx context.Context, userID, oldPass, NewPass string) error

	ChangeUsername(ctx context.Context, userID, username string) error

	AddPhone(ctx context.Context, userID string, phone user.PhoneNumber) (user.PhoneNumber, error)
	GetPhone(ctx context.Context, userID string, phoneID uint) (user.PhoneNumber, error)
	GetPhones(ctx context.Context, userID string) ([]user.PhoneNumber, error)
	DeletePhone(ctx context.Context, userID string, phoneID uint) error

	AddAddress(ctx context.Context, userID string, address user.Address) (user.Address, error)
	GetAddress(ctx context.Context, userID string, addressID uint) (user.Address, error)
	GetAddresses(ctx context.Context, userID string) ([]user.Address, error)
	DeleteAddress(ctx context.Context, userID string, addressID uint) error
}

type ValidatorRepo interface {
	DoesUserExist(ctx context.Context, userID string) (bool, error)
	DoesPhoneExist(ctx context.Context, phoneID uint) (bool, error)
	DoesAddressExist(ctx context.Context, addressID uint) (bool, error)
}
