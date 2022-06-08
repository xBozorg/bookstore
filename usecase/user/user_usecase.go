package user

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/entity/user"
)

type UseCase interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
	GetUser(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	GetUsers(ctx context.Context, req dto.GetUsersRequest) (dto.GetUsersResponse, error)
	DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.DeleteUserResponse, error)

	LoginUser(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, error)

	ChangePassword(ctx context.Context, req dto.ChangePassRequest) (dto.ChangePassResponse, error)

	ChangeUsername(ctx context.Context, req dto.ChangeUsernameRequest) (dto.ChangeUsernameResponse, error)

	AddPhone(ctx context.Context, req dto.AddPhoneRequest) (dto.AddPhoneResponse, error)
	GetPhone(ctx context.Context, req dto.GetPhoneRequest) (dto.GetPhoneResponse, error)
	GetPhones(ctx context.Context, req dto.GetPhonesRequest) (dto.GetPhonesResponse, error)
	DeletePhone(ctx context.Context, req dto.DeletePhoneRequest) (dto.DeletePhoneResponse, error)

	AddAddress(ctx context.Context, req dto.AddAddressRequest) (dto.DeleteAddressResponse, error)
	GetAddress(ctx context.Context, req dto.GetAddressRequest) (dto.GetAddressResponse, error)
	GetAddresses(ctx context.Context, req dto.GetAddressesRequest) (dto.GetAddressesResponse, error)
	DeleteAddress(ctx context.Context, req dto.DeleteAddressRequest) (dto.DeleteAddressResponse, error)
}

type UseCaseRepo struct {
	repo Repository
}

func New(r Repository) UseCaseRepo {
	return UseCaseRepo{repo: r}
}

func (u UseCaseRepo) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	user := user.User{
		Email:        req.Email,
		Password:     req.Password,
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumbers: req.PhoneNumbers,
		Addresses:    req.Addresses,
	}
	createdUser, err := u.repo.CreateUser(ctx, user)

	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{User: createdUser}, nil
}

func (u UseCaseRepo) LoginUser(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	user, err := u.repo.LoginUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}
	return dto.LoginUserResponse{User: user}, err
}

func (u UseCaseRepo) GetUser(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	user, err := u.repo.GetUser(ctx, req.UserID)
	if err != nil {
		return dto.GetUserResponse{}, err
	}

	return dto.GetUserResponse{User: user}, nil
}

func (u UseCaseRepo) GetUsers(ctx context.Context, _ dto.GetUsersRequest) (dto.GetUsersResponse, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return dto.GetUsersResponse{}, err
	}

	return dto.GetUsersResponse{Users: users}, nil
}

func (u UseCaseRepo) DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.DeleteUserResponse, error) {
	err := u.repo.DeleteUser(ctx, req.UserID)
	if err != nil {
		return dto.DeleteUserResponse{}, err
	}
	return dto.DeleteUserResponse{}, nil
}

func (u UseCaseRepo) ChangePassword(ctx context.Context, req dto.ChangePassRequest) (dto.ChangePassResponse, error) {
	err := u.repo.ChangePassword(ctx, req.UserID, req.OldPass, req.NewPass)
	if err != nil {
		return dto.ChangePassResponse{}, err
	}
	return dto.ChangePassResponse{}, nil
}

func (u UseCaseRepo) ChangeUsername(ctx context.Context, req dto.ChangeUsernameRequest) (dto.ChangeUsernameResponse, error) {
	err := u.repo.ChangeUsername(ctx, req.UserID, req.Username)
	if err != nil {
		return dto.ChangeUsernameResponse{}, err
	}
	return dto.ChangeUsernameResponse{}, nil
}

func (u UseCaseRepo) AddPhone(ctx context.Context, req dto.AddPhoneRequest) (dto.AddPhoneResponse, error) {
	p := user.PhoneNumber{
		Code:   req.Code,
		Number: req.PhoneNumber,
	}
	p, err := u.repo.AddPhone(ctx, req.UserID, p)
	if err != nil {
		return dto.AddPhoneResponse{}, err
	}
	return dto.AddPhoneResponse{Phone: p}, nil
}

func (u UseCaseRepo) GetPhone(ctx context.Context, req dto.GetPhoneRequest) (dto.GetPhoneResponse, error) {
	p, err := u.repo.GetPhone(ctx, req.UserID, req.PhoneID)
	if err != nil {
		return dto.GetPhoneResponse{}, err
	}

	return dto.GetPhoneResponse{Phone: p}, nil
}

func (u UseCaseRepo) GetPhones(ctx context.Context, req dto.GetPhonesRequest) (dto.GetPhonesResponse, error) {
	phones, err := u.repo.GetPhones(ctx, req.UserID)
	if err != nil {
		return dto.GetPhonesResponse{}, err
	}

	return dto.GetPhonesResponse{Phones: phones}, nil
}

func (u UseCaseRepo) DeletePhone(ctx context.Context, req dto.DeletePhoneRequest) (dto.DeletePhoneResponse, error) {
	err := u.repo.DeletePhone(ctx, req.UserID, req.PhoneID)
	if err != nil {
		return dto.DeletePhoneResponse{}, err
	}
	return dto.DeletePhoneResponse{}, nil
}

func (u UseCaseRepo) AddAddress(ctx context.Context, req dto.AddAddressRequest) (dto.AddAddressResponse, error) {
	address := user.Address{
		Country:     req.Country,
		Province:    req.Province,
		City:        req.City,
		Street:      req.Street,
		PostalCode:  req.PostalCode,
		No:          req.No,
		Description: req.Description,
	}
	userAddress, err := u.repo.AddAddress(ctx, req.UserID, address)
	if err != nil {
		return dto.AddAddressResponse{}, err
	}
	return dto.AddAddressResponse{Address: userAddress}, nil
}

func (u UseCaseRepo) GetAddress(ctx context.Context, req dto.GetAddressRequest) (dto.GetAddressResponse, error) {
	address, err := u.repo.GetAddress(ctx, req.UserID, req.AddressID)
	if err != nil {
		return dto.GetAddressResponse{}, err
	}
	return dto.GetAddressResponse{Address: address}, nil
}

func (u UseCaseRepo) GetAddresses(ctx context.Context, req dto.GetAddressesRequest) (dto.GetAddressesResponse, error) {
	addresses, err := u.repo.GetAddresses(ctx, req.UserID)
	if err != nil {
		return dto.GetAddressesResponse{}, err
	}

	return dto.GetAddressesResponse{Addresses: addresses}, nil
}
func (u UseCaseRepo) DeleteAddress(ctx context.Context, req dto.DeleteAddressRequest) (dto.DeleteAddressResponse, error) {
	err := u.repo.DeleteAddress(ctx, req.UserID, req.AddressID)
	if err != nil {
		return dto.DeleteAddressResponse{}, err
	}
	return dto.DeleteAddressResponse{}, nil
}
