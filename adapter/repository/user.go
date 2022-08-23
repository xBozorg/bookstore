package repository

import (
	"context"
	"errors"
	"time"

	"github.com/XBozorg/bookstore/entity/user"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (storage Storage) CreateUser(ctx context.Context, u user.User) (user.User, error) {

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return user.User{}, err
	}
	u.Password = hashedPassword

	userID := uuid.NewV4().String()

	_, err = storage.MySQL.ExecContext(
		ctx,
		"INSERT INTO user (id, email, password, username, firstname, lastname, regdate) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID,
		u.Email,
		u.Password,
		u.Username,
		u.FirstName,
		u.LastName,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return user.User{}, err
	}

	u.ID = userID
	u.Password = ""

	return u, nil
}

func (storage Storage) LoginUser(ctx context.Context, username, email, password string) (user.User, error) {

	result := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT id, email, password, username, firstname, lastname FROM user WHERE username = ? OR email = ?",
		username,
		email,
	)

	var u user.User
	var passHash string

	err := result.Scan(
		&u.ID,
		&u.Email,
		&passHash,
		&u.Username,
		&u.FirstName,
		&u.LastName,
	)

	if err != nil {
		return user.User{}, err
	}

	isSame := CheckPasswordHash(password, passHash)
	if isSame {
		return u, nil
	}
	return user.User{}, errors.New("password does not match")
}

func (storage Storage) GetUser(ctx context.Context, userID string) (user.User, error) {

	result := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT id, email, username, firstname, lastname FROM user WHERE id = ?",
		userID,
	)

	var u user.User

	err := result.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.FirstName,
		&u.LastName,
	)

	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (storage Storage) GetUsers(ctx context.Context) ([]user.User, error) {

	result, err := storage.MySQL.QueryContext(
		ctx,
		"SELECT id, email, username, firstname, lastname FROM user",
	)

	if err != nil {
		return []user.User{}, err
	}
	defer result.Close()

	users := []user.User{}
	for result.Next() {
		var u user.User

		err := result.Scan(
			&u.ID,
			&u.Email,
			&u.Username,
			&u.FirstName,
			&u.LastName,
		)

		if err != nil {
			return []user.User{}, nil
		}
		users = append(users, u)
	}

	return users, nil
}

func (storage Storage) ChangePassword(ctx context.Context, userID, oldPass, newPass string) error {

	var oldInDB string

	oldQ := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT password FROM user WHERE id = ?",
		userID,
	)

	oldQ.Scan(&oldInDB)

	isSame := CheckPasswordHash(oldPass, oldInDB)
	if isSame {
		new, err := HashPassword(newPass)
		if err != nil {
			return err
		}
		_, err = storage.MySQL.ExecContext(ctx, "UPDATE user SET password = ? WHERE id = ?", new, userID)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("password does not match")
}

func (storage Storage) ChangeUsername(ctx context.Context, userID, username string) error {

	_, err := storage.MySQL.ExecContext(
		ctx,
		"UPDATE user SET username = ? WHERE id = ?",
		username,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddPhone(ctx context.Context, userID string, phone user.PhoneNumber) (user.PhoneNumber, error) {

	noPhonesQuery := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM phone WHERE userID = ?",
		userID,
	)

	var noPhones int
	err := noPhonesQuery.Scan(&noPhones)
	if err != nil {
		return user.PhoneNumber{}, err
	}

	if noPhones >= 3 {
		return user.PhoneNumber{}, errors.New("max number of phones reached (3/3)")
	}
	_, err = storage.MySQL.ExecContext(
		ctx,
		"INSERT INTO phone (code, phonenumber, userID) VALUES (?, ?, ?)",
		phone.Code,
		phone.Number,
		userID,
	)

	if err != nil {
		return user.PhoneNumber{}, err
	}

	return phone, nil
}

func (storage Storage) GetPhone(ctx context.Context, userID string, phoneID uint) (user.PhoneNumber, error) {

	result := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT code, phoneNumber FROM phone WHERE ( userID = ? AND id = ?)",
		userID,
		phoneID,
	)

	var p user.PhoneNumber

	err := result.Scan(
		&p.Code,
		&p.Number,
	)

	if err != nil {
		return user.PhoneNumber{}, err
	}

	p.ID = phoneID

	return p, nil
}
func (storage Storage) GetPhones(ctx context.Context, userID string) ([]user.PhoneNumber, error) {

	result, err := storage.MySQL.QueryContext(
		ctx,
		"SELECT id, code, phonenumber FROM phone WHERE userID = ?",
		userID,
	)

	if err != nil {
		return []user.PhoneNumber{}, err
	}
	phones := []user.PhoneNumber{}
	for result.Next() {
		var phone user.PhoneNumber

		err := result.Scan(
			&phone.ID,
			&phone.Code,
			&phone.Number,
		)

		if err != nil {
			return []user.PhoneNumber{}, err
		}
		phones = append(phones, phone)
	}

	return phones, nil
}

func (storage Storage) DeletePhone(ctx context.Context, userID string, phoneID uint) error {

	_, err := storage.MySQL.ExecContext(ctx,
		"DELETE FROM phone WHERE userID = ? AND id = ?",
		userID,
		phoneID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddAddress(ctx context.Context, userID string, address user.Address) (user.Address, error) {

	noAddressesQuery := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM address WHERE userID = ?",
		userID,
	)

	var noAddresses int
	err := noAddressesQuery.Scan(&noAddresses)
	if err != nil {
		return user.Address{}, err
	}

	if noAddresses >= 3 {
		return user.Address{}, errors.New("max number of addresses reached (3/3)")
	}

	_, err = storage.MySQL.ExecContext(
		ctx,
		"INSERT INTO address (country, province, city, street, postalcode, no, description, userID) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		address.Country,
		address.Province,
		address.City,
		address.Street,
		address.PostalCode,
		address.No,
		address.Description,
		userID,
	)

	if err != nil {
		return user.Address{}, err
	}
	return address, nil
}

func (storage Storage) GetAddress(ctx context.Context, userID string, addressID uint) (user.Address, error) {

	result := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT country, province, city, street, postalCode, no, description FROM address WHERE userID = ? AND id = ?",
		userID,
		addressID,
	)

	var address user.Address

	err := result.Scan(
		&address.Country,
		&address.Province,
		&address.City,
		&address.Street,
		&address.PostalCode,
		&address.No,
		&address.Description,
	)

	if err != nil {
		return user.Address{}, err
	}
	return address, nil
}

func (storage Storage) GetAddresses(ctx context.Context, userID string) ([]user.Address, error) {

	result, err := storage.MySQL.QueryContext(
		ctx,
		"SELECT id, country, province, city, street, postalCode, no, description FROM address WHERE userID = ?",
		userID,
	)

	if err != nil {
		return []user.Address{}, err
	}

	var addresses []user.Address
	for result.Next() {
		var address user.Address

		err := result.Scan(
			&address.ID,
			&address.Country,
			&address.Province,
			&address.City,
			&address.Street,
			&address.PostalCode,
			&address.No,
			&address.Description,
		)

		if err != nil {
			return []user.Address{}, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}
func (storage Storage) DeleteAddress(ctx context.Context, userID string, addressID uint) error {

	_, err := storage.MySQL.ExecContext(
		ctx,
		"DELETE FROM address WHERE userID = ? AND id = ?",
		userID,
		addressID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (storage Storage) DeleteUser(ctx context.Context, userID string) error {

	result, err := storage.MySQL.QueryContext(
		ctx,
		"DELETE FROM user WHERE id = ?",
		userID,
	)

	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (storage Storage) DoesUserExist(ctx context.Context, userID string) (bool, error) {

	result := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)",
		userID,
	)

	var doesExist bool
	err := result.Scan(&doesExist)
	if err != nil {
		return false, err
	}

	return doesExist, nil
}

func (storage Storage) DoesPhoneExist(ctx context.Context, phoneID uint) (bool, error) {

	result := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM phone WHERE id = ?)",
		phoneID,
	)

	var doesExist bool
	err := result.Scan(&doesExist)
	if err != nil {
		return false, err
	}
	return doesExist, nil
}

func (storage Storage) DoesAddressExist(ctx context.Context, addressID uint) (bool, error) {

	result := storage.MySQL.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM address WHERE id = ?)",
		addressID,
	)

	var doesExist bool
	err := result.Scan(&doesExist)
	if err != nil {
		return false, err
	}
	return doesExist, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
