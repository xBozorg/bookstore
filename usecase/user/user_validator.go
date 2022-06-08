package user

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type (
	ValidateCreateUser func(req dto.CreateUserRequest) error
	ValidateGetUser    func(ctx context.Context, req dto.GetUserRequest) error
	ValidateDeleteUser func(ctx context.Context, req dto.DeleteUserRequest) error

	ValidateLoginUser func(ctx context.Context, req dto.LoginUserRequest) error

	ValidateChangePass     func(ctx context.Context, req dto.ChangePassRequest) error
	ValidateChangeUsername func(ctx context.Context, req dto.ChangeUsernameRequest) error

	ValidateAddPhone    func(ctx context.Context, req dto.AddPhoneRequest) error
	ValidateGetPhone    func(ctx context.Context, req dto.GetPhoneRequest) error
	ValidateGetPhones   func(ctx context.Context, req dto.GetPhonesRequest) error
	ValidateDeletePhone func(ctx context.Context, req dto.DeletePhoneRequest) error

	ValidateAddAddress    func(ctx context.Context, req dto.AddAddressRequest) error
	ValidateGetAddress    func(ctx context.Context, req dto.GetAddressRequest) error
	ValidateGetAddresses  func(ctx context.Context, req dto.GetAddressesRequest) error
	ValidateDeleteAddress func(ctx context.Context, req dto.DeleteAddressRequest) error
)
