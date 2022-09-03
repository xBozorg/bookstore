package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/config"
	uuid "github.com/satori/go.uuid"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
)

func GenerateTokens(c echo.Context, storage repository.Storage, tk repository.Token) error {

	accessToken, expA, err := generateAccessToken(tk)
	if err != nil {
		return err
	}

	refreshToken, expR, err := generateAndSaveRefreshToken(c, storage, tk)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, expA, c)
	setTokenCookie(refreshTokenCookieName, refreshToken, expR, c)

	return nil
}

func generateAccessToken(tk repository.Token) (string, time.Time, error) {

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &CustomClaims{
		Role: tk.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   tk.ID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Conf.GetJWTConfig().Secret))
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func GetID(c echo.Context) (string, error) {

	accessCookie, err := c.Cookie("access-token")
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(
		accessCookie.Value,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.GetJWTConfig().Secret), nil
		},
	)
	if err != nil {
		return "", err
	}

	accessClaims := token.Claims.(jwt.MapClaims)
	id := accessClaims["sub"].(string)

	return id, nil
}

func GetSignOutInfo(c echo.Context) (repository.Token, error) {

	refreshCookie, err := c.Cookie("refresh-token")
	if err != nil {
		return repository.Token{}, err
	}

	token, err := jwt.Parse(
		refreshCookie.Value,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.GetJWTConfig().Secret), nil
		},
	)
	if err != nil {
		return repository.Token{}, err
	}

	refreshClaims := token.Claims.(jwt.MapClaims)
	tk := repository.Token{
		ID:  refreshClaims["sub"].(string),
		JTI: refreshClaims["jti"].(string),
	}

	return tk, nil
}

func DeleteAccessCookie(c echo.Context) error {

	_, err := c.Cookie("access-token")
	if err != nil {
		return err
	}

	ac := &http.Cookie{
		Name:     "access-token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	c.SetCookie(ac)

	return nil
}

func generateAndSaveRefreshToken(c echo.Context, storage repository.Storage, tk repository.Token) (string, time.Time, error) {

	expirationTime := time.Now().Add(10 * 24 * time.Hour)

	tk.JTI = uuid.NewV4().String()

	claims := &CustomClaims{
		Role: tk.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   tk.ID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			ID:        tk.JTI,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Conf.GetJWTConfig().Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	tk.RefreshToken = tokenString
	tk.RefreshExp = expirationTime
	if err := storage.SaveRefreshToken(c.Request().Context(), tk); err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true
	//cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}

func UserJWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
}
func AdminJWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
}

func JWTCookieChecker(cookie *http.Cookie) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(
		cookie.Value,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.GetJWTConfig().Secret), nil
		},
	)
	if err != nil {
		return &CustomClaims{}, err
	}

	claim, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return &CustomClaims{}, errors.New("jwt.Claims to CustomClaims error")
	}

	if !(time.Until(time.Unix(int64(claim.ExpiresAt.Unix()), 0)) > 0) || token == nil || !token.Valid {
		return &CustomClaims{}, errors.New("invalid jwt cookie")
	}

	return claim, nil
}

func UserTokenRefresher(repo repository.Storage) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			accessCookie, errAccess := c.Cookie("access-token")
			refreshCookie, errRefresh := c.Cookie("refresh-token")

			switch {
			case errAccess != nil && errRefresh != nil && !strings.Contains(c.Request().URL.String(), "user/login"):
				return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")

			case errAccess == nil:
				_, err := JWTCookieChecker(accessCookie)
				if err != nil && !strings.Contains(c.Request().URL.String(), "user/login") {
					return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
				}
				if strings.Contains(c.Request().URL.String(), "user/login") {
					return c.Redirect(http.StatusMovedPermanently, "/v1")
				}
				return next(c)

			case errAccess != nil && errRefresh == nil:
				rtClaim, err := JWTCookieChecker(refreshCookie)
				if err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
				}

				tk := repository.Token{
					ID:           rtClaim.Subject,
					Role:         rtClaim.Role,
					JTI:          rtClaim.RegisteredClaims.ID,
					RefreshToken: refreshCookie.Value,
				}

				if exist, err := repo.DoesRefreshTokenExist(c.Request().Context(), tk); !exist || err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
				}

				if err := repo.DeleteRefreshToken(c.Request().Context(), tk); err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
				}

				if err := GenerateTokens(c, repo, tk); err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
				}

				if strings.Contains(c.Request().URL.String(), "user/login") {
					return c.Redirect(http.StatusMovedPermanently, "/v1")
				}

				return c.Redirect(http.StatusMovedPermanently, c.Request().URL.String())
			}

			return next(c)
		}
	}
}

func AdminTokenRefresher(repo repository.Storage) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			accessCookie, errAccess := c.Cookie("access-token")
			refreshCookie, errRefresh := c.Cookie("refresh-token")

			switch {
			case errAccess != nil && errRefresh != nil && !strings.Contains(c.Request().URL.String(), "admin/login"):
				return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")

			case errAccess == nil:
				_, err := JWTCookieChecker(accessCookie)
				if err != nil && !strings.Contains(c.Request().URL.String(), "adminlogin") {
					return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
				}
				if strings.Contains(c.Request().URL.String(), "admin/login") {
					return c.Redirect(http.StatusMovedPermanently, "/v1")
				}
				return next(c)

			case errAccess != nil && errRefresh == nil:
				rtClaim, err := JWTCookieChecker(refreshCookie)
				if err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
				}

				tk := repository.Token{
					ID:           rtClaim.Subject,
					Role:         rtClaim.Role,
					JTI:          rtClaim.RegisteredClaims.ID,
					RefreshToken: refreshCookie.Value,
				}

				if exist, err := repo.DoesRefreshTokenExist(c.Request().Context(), tk); !exist || err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
				}

				if err := repo.DeleteRefreshToken(c.Request().Context(), tk); err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
				}

				if err := GenerateTokens(c, repo, tk); err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
				}

				if strings.Contains(c.Request().URL.String(), "admin/login") {
					return c.Redirect(http.StatusMovedPermanently, "/v1")
				}

				return c.Redirect(http.StatusMovedPermanently, c.Request().URL.String())
			}

			return next(c)
		}
	}
}
