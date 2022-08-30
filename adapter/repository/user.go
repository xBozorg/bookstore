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

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`INSERT INTO user 
		(id, email, password, username, firstname, lastname, regdate) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return user.User{}, err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		userID,
		u.Email,
		u.Password,
		u.Username,
		u.FirstName,
		u.LastName,
		time.Now().Format("2006-01-02 15:04:05"),
	); err != nil {
		return user.User{}, err
	}

	u.ID = userID
	u.Password = ""

	return u, nil
}

func (storage Storage) LoginUser(ctx context.Context, username, email, password string) (user.User, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id, email, password, username, firstname, lastname FROM user WHERE username = ? OR email = ?",
	)
	if err != nil {
		return user.User{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		username,
		email,
	)

	var u user.User
	var passHash string

	if err = result.Scan(
		&u.ID,
		&u.Email,
		&passHash,
		&u.Username,
		&u.FirstName,
		&u.LastName,
	); err != nil {
		return user.User{}, err
	}

	isSame := CheckPasswordHash(password, passHash)
	if isSame {
		return u, nil
	}
	return user.User{}, errors.New("password does not match")
}

func (storage Storage) GetUser(ctx context.Context, userID string) (user.User, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id, email, username, firstname, lastname FROM user WHERE id = ?",
	)
	if err != nil {
		return user.User{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, userID)

	var u user.User

	if err = result.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.FirstName,
		&u.LastName,
	); err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (storage Storage) GetUsers(ctx context.Context) ([]user.User, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id, email, username, firstname, lastname FROM user",
	)
	if err != nil {
		return []user.User{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []user.User{}, err
	}
	defer result.Close()

	users := []user.User{}
	for result.Next() {
		var u user.User

		if err = result.Scan(
			&u.ID,
			&u.Email,
			&u.Username,
			&u.FirstName,
			&u.LastName,
		); err != nil {
			return []user.User{}, nil
		}
		users = append(users, u)
	}

	return users, nil
}

func (storage Storage) ChangePassword(ctx context.Context, userID, oldPass, newPass string) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT password FROM user WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	oldQ := stmt.QueryRowContext(ctx, userID)

	var oldInDB string
	if err = oldQ.Scan(&oldInDB); err != nil {
		return err
	}

	isSame := CheckPasswordHash(oldPass, oldInDB)
	if isSame {
		new, err := HashPassword(newPass)
		if err != nil {
			return err
		}

		stmt, err := storage.MySQL.PrepareContext(ctx,
			"UPDATE user SET password = ? WHERE id = ?",
		)
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err = stmt.ExecContext(ctx, new, userID); err != nil {
			return err
		}

		return nil
	}

	return errors.New("password does not match")
}

func (storage Storage) ChangeUsername(ctx context.Context, userID, username string) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"UPDATE user SET username = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		username,
		userID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddPhone(ctx context.Context, userID string, phone user.PhoneNumber) (user.PhoneNumber, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT COUNT(*) FROM phone WHERE userID = ?",
	)
	if err != nil {
		return user.PhoneNumber{}, err
	}
	defer stmt.Close()

	var noPhones int

	noPhonesQuery := stmt.QueryRowContext(ctx, userID)
	if err = noPhonesQuery.Scan(&noPhones); err != nil {
		return user.PhoneNumber{}, err
	}

	if noPhones >= 3 {
		return user.PhoneNumber{}, errors.New("max number of phones reached (3/3)")
	}

	stmt, err = storage.MySQL.PrepareContext(ctx,
		"INSERT INTO phone (code, phonenumber, userID) VALUES (?, ?, ?)",
	)
	if err != nil {
		return user.PhoneNumber{}, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		phone.Code,
		phone.Number,
		userID,
	); err != nil {
		return user.PhoneNumber{}, err
	}

	return phone, nil
}

func (storage Storage) GetPhone(ctx context.Context, userID string, phoneID uint) (user.PhoneNumber, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT code, phoneNumber FROM phone WHERE ( userID = ? AND id = ?)",
	)
	if err != nil {
		return user.PhoneNumber{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		userID,
		phoneID,
	)

	var p user.PhoneNumber

	if err = result.Scan(
		&p.Code,
		&p.Number,
	); err != nil {
		return user.PhoneNumber{}, err
	}

	p.ID = phoneID

	return p, nil
}

func (storage Storage) GetPhones(ctx context.Context, userID string) ([]user.PhoneNumber, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id, code, phonenumber FROM phone WHERE userID = ?",
	)
	if err != nil {
		return []user.PhoneNumber{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return []user.PhoneNumber{}, err
	}

	phones := []user.PhoneNumber{}
	for result.Next() {
		var phone user.PhoneNumber

		if err = result.Scan(
			&phone.ID,
			&phone.Code,
			&phone.Number,
		); err != nil {
			return []user.PhoneNumber{}, err
		}
		phones = append(phones, phone)
	}

	return phones, nil
}

func (storage Storage) DeletePhone(ctx context.Context, userID string, phoneID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM phone WHERE userID = ? AND id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = storage.MySQL.ExecContext(ctx,
		userID,
		phoneID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddAddress(ctx context.Context, userID string, address user.Address) (user.Address, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT COUNT(*) FROM address WHERE userID = ?",
	)
	if err != nil {
		return user.Address{}, err
	}
	defer stmt.Close()

	noAddressesQuery := stmt.QueryRowContext(ctx, userID)

	var noAddresses int
	if err = noAddressesQuery.Scan(&noAddresses); err != nil {
		return user.Address{}, err
	}

	if noAddresses >= 3 {
		return user.Address{}, errors.New("max number of addresses reached (3/3)")
	}

	stmt, err = storage.MySQL.PrepareContext(ctx,
		`INSERT INTO address 
		(country, province, city, street, postalcode, no, description, userID) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return user.Address{}, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		address.Country,
		address.Province,
		address.City,
		address.Street,
		address.PostalCode,
		address.No,
		address.Description,
		userID,
	); err != nil {
		return user.Address{}, err
	}

	return address, nil
}

func (storage Storage) GetAddress(ctx context.Context, userID string, addressID uint) (user.Address, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT country, province, city, street, postalCode, no, description FROM address 
		WHERE userID = ? AND id = ?`,
	)
	if err != nil {
		return user.Address{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		userID,
		addressID,
	)

	var address user.Address

	if err = result.Scan(
		&address.Country,
		&address.Province,
		&address.City,
		&address.Street,
		&address.PostalCode,
		&address.No,
		&address.Description,
	); err != nil {
		return user.Address{}, err
	}
	return address, nil
}

func (storage Storage) GetAddresses(ctx context.Context, userID string) ([]user.Address, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT id, country, province, city, street, postalCode, no, description FROM address 
		WHERE userID = ?`,
	)
	if err != nil {
		return []user.Address{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return []user.Address{}, err
	}

	var addresses []user.Address
	for result.Next() {
		var address user.Address

		if err = result.Scan(
			&address.ID,
			&address.Country,
			&address.Province,
			&address.City,
			&address.Street,
			&address.PostalCode,
			&address.No,
			&address.Description,
		); err != nil {
			return []user.Address{}, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}
func (storage Storage) DeleteAddress(ctx context.Context, userID string, addressID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM address WHERE userID = ? AND id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		userID,
		addressID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) DeleteUser(ctx context.Context, userID string) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM user WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (storage Storage) DoesUserExist(ctx context.Context, userID string) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, userID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
		return false, err
	}

	return doesExist, nil
}

func (storage Storage) DoesPhoneExist(ctx context.Context, phoneID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM phone WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, phoneID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
		return false, err
	}
	return doesExist, nil
}

func (storage Storage) DoesAddressExist(ctx context.Context, addressID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM address WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, addressID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
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
