package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/XBozorg/bookstore/adapter/auth"
	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/user"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func CreateUser(storage repository.Storage, validator user.ValidateCreateUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		createUserReq := dto.CreateUserRequest{}
		if err := c.Bind(&createUserReq); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(createUserReq); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		createUserResp, err := user.New(storage).CreateUser(c.Request().Context(), createUserReq)

		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {

				if strings.Contains(driverErr.Message, "user.email") {
					return echo.NewHTTPError(http.StatusConflict, "email already exists")
				}
				if strings.Contains(driverErr.Message, "user.username") {
					return echo.NewHTTPError(http.StatusConflict, "username already exists")
				}

			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		loginUserResp := dto.LoginUserResponse(createUserResp)

		if err := auth.GenerateTokens(c, storage,
			repository.Token{
				ID:   loginUserResp.User.ID,
				Role: "user",
			},
		); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, loginUserResp)
	}
}

func LoginUser(storage repository.Storage, validator user.ValidateLoginUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.LoginUserRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).LoginUser(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound, "user not found")
			}
			if err.Error() == "password does not match" {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := auth.GenerateTokens(c, storage,
			repository.Token{
				ID:   resp.User.ID,
				Role: "user",
			},
		); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func UserLoginForm() echo.HandlerFunc { // simple handler for redirect unauthenticated users
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "login page")
	}
}

func GetUser(storage repository.Storage, validator user.ValidateGetUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.GetUserRequest{}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).GetUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetUsers(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetUsersRequest{}

		resp, err := user.New(storage).GetUsers(c.Request().Context(), req)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteUser(storage repository.Storage, validator user.ValidateDeleteUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.DeleteUserRequest{}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err = user.New(storage).DeleteUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func ChangePassword(storage repository.Storage, validator user.ValidateChangePass) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.ChangePassRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err = user.New(storage).ChangePassword(c.Request().Context(), req)
		if err != nil {
			if err.Error() == "password does not match" {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		return c.NoContent(http.StatusOK)
	}
}

func ChangeUsername(storage repository.Storage, validator user.ValidateChangeUsername) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.ChangeUsernameRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).ChangeUsername(c.Request().Context(), req)
		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "username already exists")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func AddPhone(storage repository.Storage, validator user.ValidateAddPhone) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddPhoneRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).AddPhone(c.Request().Context(), req)
		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "Phonenumber already exists")
			}
			if strings.Contains(err.Error(), "max") {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetPhone(storage repository.Storage, validator user.ValidateGetPhone) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPhoneRequest{}
		pid, err := strconv.ParseUint(c.Param("phoneID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id
		req.PhoneID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "phonenumber does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).GetPhone(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetPhones(storage repository.Storage, validator user.ValidateGetPhones) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPhonesRequest{}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).GetPhones(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DeletePhone(storage repository.Storage, validator user.ValidateDeletePhone) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeletePhoneRequest{}
		pid, err := strconv.ParseUint(c.Param("phoneID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id
		req.PhoneID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).DeletePhone(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func AddAddress(storage repository.Storage, validator user.ValidateAddAddress) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddAddressRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).AddAddress(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "max") {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAddress(storage repository.Storage, validator user.ValidateGetAddress) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.GetAddressRequest{}

		aid, err := strconv.ParseUint(c.Param("addressID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id
		req.AddressID = uint(aid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "address does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).GetAddress(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAddresses(storage repository.Storage, validator user.ValidateGetAddresses) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAddressesRequest{}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).GetAddresses(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteAddress(storage repository.Storage, validator user.ValidateDeleteAddress) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeleteAddressRequest{}
		aid, err := strconv.ParseUint(c.Param("addressID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id
		req.AddressID = uint(aid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(storage).DeleteAddress(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func UserLogOut(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {

		tk, err := auth.GetSignOutInfo(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		err = storage.DeleteRefreshToken(c.Request().Context(), tk)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		err = auth.DeleteAccessCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.Redirect(http.StatusMovedPermanently, "/v1")
	}
}

func UserLogOutAllDevices(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		err = storage.DeleteRefreshTokens(c.Request().Context(), "user", id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		err = auth.DeleteAccessCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.Redirect(http.StatusMovedPermanently, "/v1")
	}
}
