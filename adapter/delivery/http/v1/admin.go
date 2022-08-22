package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/XBozorg/bookstore/adapter/auth"
	repository "github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/config"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/admin"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		adminTokenCookie, err := c.Cookie("ID")
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
		}

		accessTokenCookie, err := c.Cookie("access-token")
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
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
			return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
		}

		if _, ok := token.Claims.(*auth.Claims); !ok || !token.Valid || claims.Role != "admin" || claims.Subject != adminTokenCookie.Value {
			// use as error string for debugging -> fmt.Sprintf("Claims is valid:%v , Token is valid:%v , Role is User:%v , userId of userToken and accessToken is the same:%v ", ok, token.Valid, claims.Role == "user",claims.Subject != userTokenCookie.Value)
			return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
		}

		return next(c)
	}
}

func GetAdmin(repo repository.Repo, validator admin.ValidateGetAdmin) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.GetAdminRequest{}
		adminCookie, _ := c.Cookie("ID")
		req.AdminId = adminCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := admin.New(repo).GetAdmin(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetAdmins(repo repository.Repo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAdminsRequest{}

		resp, err := admin.New(repo).GetAdmins(c.Request().Context(), req)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func LoginAdmin(repo repository.Repo, validator admin.ValidateLoginAdmin) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.LoginAdminRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := admin.New(repo).LoginAdmin(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			if err.Error() == "password does not match" {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = auth.GenerateTokensAndSetCookies(c, resp.Admin.ID, "admin")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func AdminLoginForm() echo.HandlerFunc { // simple handler for redirect unauthenticated admins
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "admin login page")
	}
}
