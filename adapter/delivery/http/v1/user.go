package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/XBozorg/bookstore/adapter/auth"
	repository "github.com/XBozorg/bookstore/adapter/repository/mysql"
	"github.com/XBozorg/bookstore/config"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func UserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userTokenCookie, err := c.Cookie("ID")
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
		}
		accessTokenCookie, err := c.Cookie("access-token")
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
		}

		accessToken := accessTokenCookie.Value

		claims := auth.Claims{}
		token, err := jwt.ParseWithClaims(accessToken, &claims,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.Conf.GetJWTConfig().Secret), nil
			})

		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
		}

		if _, ok := token.Claims.(*auth.Claims); !ok || !token.Valid || claims.Role != "user" || claims.Subject != userTokenCookie.Value {
			// use as error string for debugging -> fmt.Sprintf("Claims is valid:%v , Token is valid:%v , Role is User:%v , userId of userToken and accessToken is the same:%v ", ok, token.Valid, claims.Role == "user",claims.Subject != userTokenCookie.Value)
			return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
		}

		return next(c)
	}
}

func CreateUser(repo repository.MySQLRepo, validator user.ValidateCreateUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		createUserReq := dto.CreateUserRequest{}
		if err := c.Bind(&createUserReq); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(createUserReq); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		createUserResp, err := user.New(repo).CreateUser(c.Request().Context(), createUserReq)

		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				if strings.Contains(err.Error(), "user_email") {
					return echo.NewHTTPError(http.StatusConflict, "email already exists")
				}
				if strings.Contains(err.Error(), "user_name") {
					return echo.NewHTTPError(http.StatusConflict, "username already exists")
				}
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		loginUserResp := dto.LoginUserResponse(createUserResp)

		err = auth.GenerateTokensAndSetCookies(c, loginUserResp.User.ID, "user")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, loginUserResp)
	}
}

func LoginUser(repo repository.MySQLRepo, validator user.ValidateLoginUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.LoginUserRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).LoginUser(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound, "user not found")
			}
			if err.Error() == "password does not match" {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = auth.GenerateTokensAndSetCookies(c, resp.User.ID, "user")
		if err != nil {
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

func GetUser(repo repository.MySQLRepo, validator user.ValidateGetUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.GetUserRequest{}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).GetUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetUsers(repo repository.MySQLRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetUsersRequest{}

		resp, err := user.New(repo).GetUsers(c.Request().Context(), req)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteUser(repo repository.MySQLRepo, validator user.ValidateDeleteUser) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.DeleteUserRequest{}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err := user.New(repo).DeleteUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func ChangePassword(repo repository.MySQLRepo, validator user.ValidateChangePass) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.ChangePassRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err := user.New(repo).ChangePassword(c.Request().Context(), req)
		if err != nil {
			if err.Error() == "password doesn't match" {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		return c.NoContent(http.StatusOK)
	}
}

func ChangeUsername(repo repository.MySQLRepo, validator user.ValidateChangeUsername) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.ChangeUsernameRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).ChangeUsername(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				return echo.NewHTTPError(http.StatusConflict, "username already exists")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func AddPhone(repo repository.MySQLRepo, validator user.ValidateAddPhone) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddPhoneRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).AddPhone(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "max") {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}
func GetPhone(repo repository.MySQLRepo, validator user.ValidateGetPhone) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPhoneRequest{}
		pid, err := strconv.ParseUint(c.Param("phoneID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value
		req.PhoneID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "phone does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).GetPhone(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}
func GetPhones(repo repository.MySQLRepo, validator user.ValidateGetPhones) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPhonesRequest{}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).GetPhones(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}
func DeletePhone(repo repository.MySQLRepo, validator user.ValidateDeletePhone) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeletePhoneRequest{}
		pid, err := strconv.ParseUint(c.Param("phoneID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value
		req.PhoneID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).DeletePhone(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func AddAddress(repo repository.MySQLRepo, validator user.ValidateAddAddress) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddAddressRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).AddAddress(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "max") {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func GetAddress(repo repository.MySQLRepo, validator user.ValidateGetAddress) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.GetAddressRequest{}

		aid, err := strconv.ParseUint(c.Param("addressID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value
		req.AddressID = uint(aid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "address does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).GetAddress(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func GetAddresses(repo repository.MySQLRepo, validator user.ValidateGetAddresses) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAddressesRequest{}
		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).GetAddresses(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func DeleteAddress(repo repository.MySQLRepo, validator user.ValidateDeleteAddress) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeleteAddressRequest{}
		aid, err := strconv.ParseUint(c.Param("addressID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value
		req.AddressID = uint(aid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := user.New(repo).DeleteAddress(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}
