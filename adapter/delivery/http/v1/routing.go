package v1

import (
	"github.com/XBozorg/bookstore/adapter/auth"
	repository "github.com/XBozorg/bookstore/adapter/repository/mysql"
	"github.com/XBozorg/bookstore/config"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routing(repo repository.MySQLRepo) *echo.Echo {
	e := echo.New()

	userGroup := e.Group("/v1/user")

	userGroup.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			Claims:                  &auth.Claims{},
			SigningKey:              []byte(config.Conf.GetJWTConfig().Secret),
			TokenLookup:             "cookie:access-token",
			ErrorHandlerWithContext: auth.JWTErrorChecker,
			SigningMethod:           "HS256",
		},
	))
	userGroup.Use(auth.TokenRefresherMiddleware(dto.LoginUserResponse{}))

	userGroup.Use(UserAuth)

	userGroup.POST("", CreateUser(repo, validator.ValidateCreateUser))        // .../v1/user
	userGroup.GET("", GetUser(repo, validator.ValidateGetUser(repo)))         // .../v1/user
	e.GET("v1/users", GetUsers(repo))                                         // .../v1/users
	e.DELETE("v1/user", DeleteUser(repo, validator.ValidateDeleteUser(repo))) // .../v1/user

	e.POST("v1/user/login", LoginUser(repo, validator.ValidateLoginUser(repo)))              // .../v1/user/login
	e.GET("v1/user/login", LoginForm())                                                      // .../v1/user/login
	userGroup.PUT("/password", ChangePassword(repo, validator.ValidateChangePass(repo)))     // .../v1/user/password
	userGroup.PUT("/username", ChangeUsername(repo, validator.ValidateChangeUsername(repo))) // .../v1/user/username

	userGroup.POST("/phone", AddPhone(repo, validator.ValidateAddPhone(repo)))                  // .../v1/user/phone
	userGroup.GET("/phone/:phoneID", GetPhone(repo, validator.ValidateGetPhone(repo)))          // .../v1/user/phone/:phoneID
	userGroup.GET("/phone", GetPhones(repo, validator.ValidateGetPhones(repo)))                 // .../v1/user/phone
	userGroup.DELETE("/phone/:phoneID", DeletePhone(repo, validator.ValidateDeletePhone(repo))) // .../v1/user/phone/:phoneID

	userGroup.POST("/address", AddAddress(repo, validator.ValidateAddAddress(repo)))                    // .../v1/user/address
	userGroup.GET("/address/:addressID", GetAddress(repo, validator.ValidateGetAddress(repo)))          // .../v1/user/address/:addressID
	userGroup.GET("/address", GetAddresses(repo, validator.ValidateGetAddresses(repo)))                 // .../v1/user/address
	userGroup.DELETE("/address/:addressID", DeleteAddress(repo, validator.ValidateDeleteAddress(repo))) // .../v1/user/address/:addressID

	return e
}
