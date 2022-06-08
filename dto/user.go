package dto

import "github.com/XBozorg/bookstore/entity/user"

type CreateUserRequest struct {
	Email        string             `json:"email"`
	Password     string             `json:"password"`
	Username     string             `json:"username"`
	FirstName    string             `json:"firstName"`
	LastName     string             `json:"lastName"`
	PhoneNumbers []user.PhoneNumber `json:"phoneNumbers"`
	Addresses    []user.Address     `json:"addresses"`
}

type CreateUserResponse struct {
	User user.User `json:"user"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginUserResponse struct {
	User user.User `json:"user"`
}

type GetUserRequest struct {
	UserID string `json:"userID"`
}
type GetUserResponse struct {
	User user.User `json:"user"`
}

type GetUsersRequest struct {
}
type GetUsersResponse struct {
	Users []user.User `json:"users"`
}

type DeleteUserRequest struct {
	UserID string `json:"userID"`
}
type DeleteUserResponse struct {
}

type ChangePassRequest struct {
	UserID  string `json:"userID"`
	OldPass string `json:"oldPass"`
	NewPass string `json:"newPass"`
}
type ChangePassResponse struct {
}

type ChangeUsernameRequest struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
}
type ChangeUsernameResponse struct {
}

type AddPhoneRequest struct {
	UserID      string `json:"userID"`
	Code        string `json:"code"`
	PhoneNumber string `json:"phoneNumber"`
}
type AddPhoneResponse struct {
	Phone user.PhoneNumber `json:"phone"`
}

type GetPhoneRequest struct {
	UserID  string `json:"userID"`
	PhoneID uint   `json:"phoneID"`
}
type GetPhoneResponse struct {
	Phone user.PhoneNumber `json:"phone"`
}

type GetPhonesRequest struct {
	UserID string `json:"userID"`
}
type GetPhonesResponse struct {
	Phones []user.PhoneNumber `json:"phones"`
}

type DeletePhoneRequest struct {
	UserID  string `json:"userID"`
	PhoneID uint   `json:"phoneID"`
}
type DeletePhoneResponse struct {
}

type AddAddressRequest struct {
	UserID      string `json:"userID"`
	Country     string `json:"country"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
	No          string `json:"no"`
	Description string `json:"description"`
}
type AddAddressResponse struct {
	Address user.Address `json:"address"`
}

type GetAddressRequest struct {
	UserID    string `json:"userID"`
	AddressID uint   `json:"addressID"`
}
type GetAddressResponse struct {
	Address user.Address `json:"address"`
}

type GetAddressesRequest struct {
	UserID string `json:"userID"`
}
type GetAddressesResponse struct {
	Addresses []user.Address `json:"addresses"`
}

type DeleteAddressRequest struct {
	UserID    string `json:"userID"`
	AddressID uint   `json:"addressID"`
}
type DeleteAddressResponse struct {
}
