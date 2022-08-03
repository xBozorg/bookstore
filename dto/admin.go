package dto

import "github.com/XBozorg/bookstore/entity/admin"

type GetAdminRequest struct {
	AdminId string `json:"adminID"`
}
type GetAdminResponse struct {
	Admin admin.Admin `json:"admin"`
}

type GetAdminsRequest struct {
}
type GetAdminsResponse struct {
	Admins []admin.Admin `json:"admins"`
}

type LoginAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginAdminResponse struct {
	Admin admin.Admin `json:"admin"`
}

type DoesAdminExistRequest struct{}
type DoesAdminExistResponse struct{}
