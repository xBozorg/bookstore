package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/XBozorg/bookstore/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
)

func GenerateTokensAndSetCookies(c echo.Context, id, role string) error {

	accessToken, expA, err := generateAccessToken(config.Conf.GetJWTConfig().Secret, id, role)
	if err != nil {
		return err
	}

	refreshToken, expR, err := generateRefreshToken(config.Conf.GetJWTConfig().RefreshSecret, id, role)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, expA, c)
	setTokenCookie(refreshTokenCookieName, refreshToken, expR, c)
	setIDCookie(id, expA, c)

	return nil
}

func generateAccessToken(secret, id, role string) (string, time.Time, error) {

	expirationTime := time.Now().Add(15 * time.Minute)

	return generateToken(id, role, expirationTime, []byte(secret))
}

func generateRefreshToken(refreshSecret, id, role string) (string, time.Time, error) {

	expirationTime := time.Now().Add(10 * 24 * time.Hour)

	return generateToken(id, role, expirationTime, []byte(refreshSecret))
}

func generateToken(id, role string, expirationTime time.Time, secret []byte) (string, time.Time, error) {

	claims := &Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
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

func setIDCookie(id string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "ID"
	cookie.Value = id
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true
	//cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}

func CheckUserID(c echo.Context, conf config.JwtConfig, userID string) (bool, error) {
	accessTokenCookie, err := c.Cookie("access-token")
	if err != nil {
		return false, err
	}

	tokenString := accessTokenCookie.Value

	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(conf.Secret), nil
		})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["sub"] == userID {
			return true, nil
		}

	} else {
		return false, err
	}

	return false, nil
}

func UserJWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
}
func AdminJWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/v1/admin/login")
}

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		_, errAccess := c.Cookie("access-token")
		_, errID := c.Cookie("ID")
		refreshCookie, errRefresh := c.Cookie("refresh-token")

		// if access-token or ID not exists but refresh-token exists
		if (errAccess != nil || errID != nil) && errRefresh == nil {

			token, err := jwt.Parse(
				refreshCookie.Value,
				func(token *jwt.Token) (interface{}, error) {
					return []byte(config.Conf.GetJWTConfig().RefreshSecret), nil
				},
			)

			if err != nil {
				return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
			}

			claims := token.Claims.(jwt.MapClaims)
			role := claims["role"].(string)
			id := claims["sub"].(string)
			exp := claims["exp"].(float64)

			if (time.Until(time.Unix(int64(exp), 0)) > 0) && token != nil && token.Valid {
				err = GenerateTokensAndSetCookies(c, id, role)
				if err != nil {
					return c.Redirect(http.StatusMovedPermanently, "/v1/user/login")
				}
			}

			return c.Redirect(http.StatusMovedPermanently, c.Request().URL.String())
		}

		return next(c)
	}
}
