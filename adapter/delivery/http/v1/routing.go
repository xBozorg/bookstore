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
	adminGroup := e.Group("/v1/admin")

	userGroup.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			Claims:                  &auth.Claims{},
			SigningKey:              []byte(config.Conf.GetJWTConfig().Secret),
			TokenLookup:             "cookie:access-token",
			ErrorHandlerWithContext: auth.UserJWTErrorChecker,
			SigningMethod:           "HS256",
		},
	))
	adminGroup.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			Claims:                  &auth.Claims{},
			SigningKey:              []byte(config.Conf.GetJWTConfig().Secret),
			TokenLookup:             "cookie:access-token",
			ErrorHandlerWithContext: auth.AdminJWTErrorChecker,
			SigningMethod:           "HS256",
		},
	))

	userGroup.Use(auth.TokenRefresherMiddleware(dto.LoginUserResponse{}.User.ID, "user"))
	adminGroup.Use(auth.TokenRefresherMiddleware(dto.LoginAdminResponse{}.Admin.ID, "admin"))

	userGroup.Use(UserAuth)
	adminGroup.Use(AdminAuth)

	e.POST("v1/user", CreateUser(repo, validator.ValidateCreateUser))              // <Create User>    .../v1/user
	e.POST("v1/admin/login", LoginAdmin(repo, validator.ValidateLoginAdmin(repo))) // <LoginAdmin>     .../v1/admin/login
	e.GET("v1/admin/login", AdminLoginForm())                                      // <AdminLoginForm> .../v1/admin/login
	e.POST("v1/user/login", LoginUser(repo, validator.ValidateLoginUser(repo)))    // <LoginUser>      .../v1/user/login
	e.GET("v1/user/login", UserLoginForm())                                        // <UserLoginForm>  .../v1/user/login

	userGroup.GET("", GetUser(repo, validator.ValidateGetUser(repo)))                                   // <GetUser>        .../v1/user
	userGroup.DELETE("", DeleteUser(repo, validator.ValidateDeleteUser(repo)))                          // <DeleteUser>     .../v1/user
	userGroup.PUT("/password", ChangePassword(repo, validator.ValidateChangePass(repo)))                // <ChangePassword> .../v1/user/password
	userGroup.PUT("/username", ChangeUsername(repo, validator.ValidateChangeUsername(repo)))            // <ChangeUsername> .../v1/user/username
	userGroup.POST("/phone", AddPhone(repo, validator.ValidateAddPhone(repo)))                          // <AddPhone>       .../v1/user/phone
	userGroup.GET("/phone/:phoneID", GetPhone(repo, validator.ValidateGetPhone(repo)))                  // <GetPhone>       .../v1/user/phone/:phoneID
	userGroup.GET("/phone", GetPhones(repo, validator.ValidateGetPhones(repo)))                         // <GetPhones>      .../v1/user/phone
	userGroup.DELETE("/phone/:phoneID", DeletePhone(repo, validator.ValidateDeletePhone(repo)))         // <DeletePhone>    .../v1/user/phone/:phoneID
	userGroup.POST("/address", AddAddress(repo, validator.ValidateAddAddress(repo)))                    // <AddAddress>     .../v1/user/address
	userGroup.GET("/address/:addressID", GetAddress(repo, validator.ValidateGetAddress(repo)))          // <GetAddress>     .../v1/user/address/:addressID
	userGroup.GET("/address", GetAddresses(repo, validator.ValidateGetAddresses(repo)))                 // <GetAddresses>   .../v1/user/address
	userGroup.DELETE("/address/:addressID", DeleteAddress(repo, validator.ValidateDeleteAddress(repo))) // <DeleteAddress>  .../v1/user/address/:addressID

	adminGroup.GET("/users", GetUsers(repo))                             // <GetUsers>  .../v1/admin/users
	adminGroup.GET("", GetAdmin(repo, validator.ValidateGetAdmin(repo))) // <GetAdmin>  .../v1/admin
	adminGroup.GET("s", GetAdmins(repo))                                 // <GetAdmins> .../v1/admins

	return e
}
