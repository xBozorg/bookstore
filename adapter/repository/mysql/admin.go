package repository

import (
	"context"
	"errors"

	"github.com/XBozorg/bookstore/entity/admin"
)

func (m MySQLRepo) LoginAdmin(ctx context.Context, email, password string) (admin.Admin, error) {

	result := m.db.QueryRowContext(ctx, "SELECT id, email, password, phonenumber FROM admin WHERE email = ?", email)

	var a admin.Admin
	var passHash string

	err := result.Scan(&a.ID, &a.Email, &passHash, &a.PhoneNumber)
	if err != nil {
		return admin.Admin{}, err
	}

	isSame := CheckPasswordHash(password, passHash)
	if isSame {
		return a, nil
	}
	return admin.Admin{}, errors.New("password does not match")
}

func (m MySQLRepo) DoesAdminExist(ctx context.Context, adminID string) (bool, error) {

	result := m.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM admin WHERE id = ?)", adminID)

	var doesExist bool
	err := result.Scan(&doesExist)
	if err != nil {
		return false, err
	}

	return doesExist, nil
}

func (m MySQLRepo) GetAdmin(ctx context.Context, adminID string) (admin.Admin, error) {

	result := m.db.QueryRowContext(ctx, "SELECT id, email, phonenumber FROM admin WHERE id = ?", adminID)

	a := admin.Admin{}

	err := result.Scan(&a.ID, &a.Email, &a.PhoneNumber)
	if err != nil {
		return admin.Admin{}, err
	}

	return a, nil
}

func (m MySQLRepo) GetAdmins(ctx context.Context) ([]admin.Admin, error) {

	result, err := m.db.QueryContext(ctx, "SELECT id, email, phonenumber FROM admin")
	if err != nil {
		return []admin.Admin{}, err
	}
	defer result.Close()

	admins := []admin.Admin{}
	for result.Next() {
		a := admin.Admin{}
		err := result.Scan(&a.ID, &a.Email, &a.PhoneNumber)
		if err != nil {
			return []admin.Admin{}, nil
		}
		admins = append(admins, a)
	}

	return admins, nil
}
