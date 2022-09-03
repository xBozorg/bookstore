package v1

import (
	"net/http"
	"strings"

	"github.com/XBozorg/bookstore/adapter/auth"
	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/admin"
	"github.com/labstack/echo/v4"
)

func GetAdmin(storage repository.Storage, validator admin.ValidateGetAdmin) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.GetAdminRequest{}
		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.AdminId = id

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := admin.New(storage).GetAdmin(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetAdmins(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAdminsRequest{}

		resp, err := admin.New(storage).GetAdmins(c.Request().Context(), req)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func LoginAdmin(storage repository.Storage, validator admin.ValidateLoginAdmin) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.LoginAdminRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := admin.New(storage).LoginAdmin(c.Request().Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			if err.Error() == "password does not match" {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := auth.GenerateTokens(c, storage,
			repository.Token{
				ID:   resp.Admin.ID,
				Role: "admin",
			},
		); err != nil {
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

func AdminLogOut(storage repository.Storage) echo.HandlerFunc {
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

func AdminLogOutAllDevices(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		err = storage.DeleteRefreshTokens(c.Request().Context(), "admin", id)
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
