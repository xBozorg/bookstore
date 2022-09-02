package v1

import (
	"net/http"

	"github.com/XBozorg/bookstore/adapter/auth"
	"github.com/XBozorg/bookstore/adapter/payment"
	"github.com/XBozorg/bookstore/adapter/repository"

	"github.com/XBozorg/bookstore/config"
	"github.com/XBozorg/bookstore/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Home() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Bookstore Home Page")
	}
}

func Routing(storage repository.Storage) *echo.Echo {
	e := echo.New()

	userGroup := e.Group("/v1/user")
	adminGroup := e.Group("/v1/admin")

	userGroup.Use(auth.UserTokenRefresher(storage))
	adminGroup.Use(auth.AdminTokenRefresher(storage))

	userGroup.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			Claims:                  &auth.Claims{},
			SigningKey:              []byte(config.Conf.GetJWTConfig().Secret),
			TokenLookup:             "cookie:access-token,cookie:refresh-token",
			ErrorHandlerWithContext: auth.UserJWTErrorChecker,
			SigningMethod:           "HS256",
		},
	))
	adminGroup.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			Claims:                  &auth.Claims{},
			SigningKey:              []byte(config.Conf.GetJWTConfig().Secret),
			TokenLookup:             "cookie:access-token,cookie:refresh-token",
			ErrorHandlerWithContext: auth.AdminJWTErrorChecker,
			SigningMethod:           "HS256",
		},
	))

	userGroup.Use(UserAuth)
	adminGroup.Use(AdminAuth)

	e.GET("v1", Home())

	e.POST("v1/user", CreateUser(storage, validator.ValidateCreateUser))                                              // <Create User>       .../v1/user
	e.POST("v1/admin/login", LoginAdmin(storage, validator.ValidateLoginAdmin(storage)))                              // <LoginAdmin>        .../v1/admin/login
	e.GET("v1/admin/login", AdminLoginForm())                                                                         // <AdminLoginForm>    .../v1/admin/login
	e.POST("v1/user/login", LoginUser(storage, validator.ValidateLoginUser(storage)))                                 // <LoginUser>         .../v1/user/login
	e.GET("v1/user/login", UserLoginForm())                                                                           // <UserLoginForm>     .../v1/user/login
	e.GET("v1/author/:authorID", GetAuthor(storage, validator.ValidateGetAuthor(storage)))                            // <GetAuthor>         .../v1/author/:authorID
	e.GET("v1/author", GetAuthors(storage))                                                                           // <GetAuthors>        .../v1/author
	e.GET("v1/publisher/:publisherID", GetPublisher(storage, validator.ValidateGetPublisher(storage)))                // <GetPublisher>      .../v1/publisher/:publisherID
	e.GET("v1/publisher", GetPublishers(storage))                                                                     // <GetPublishers>     .../v1/publisher
	e.GET("v1/topic/:topicID", GetTopic(storage, validator.ValidateGetTopic(storage)))                                // <GetTopic>          .../v1/topic/:topicID
	e.GET("v1/topic", GetTopics(storage))                                                                             // <GetTopics>         .../v1/topic
	e.GET("v1/lang/:langID", GetLanguage(storage, validator.ValidateGetLanguage(storage)))                            // <GetLanguage>       .../v1/lang/:langID
	e.GET("v1/lang", GetLanguages(storage))                                                                           // <GetLanguages>      .../v1/lang
	e.GET("v1/book/:bookID", GetBook(storage, validator.ValidateGetBook(storage)))                                    // <GetBook>           .../v1/book/:bookID
	e.GET("v1/book", GetAllBooks(storage))                                                                            // <GetAllBooks>       .../v1/book
	e.GET("v1/book/author/:authorID", GetAuthorBooks(storage, validator.ValidateGetAuthorBooks(storage)))             // <GetAuthorBooks>    .../v1/book/author/:authorID
	e.GET("v1/book/publisher/:publisherID", GetPublisherBooks(storage, validator.ValidateGetPublisherBooks(storage))) // <GetPublisherBooks> .../v1/book/publisher/:publisherID
	e.GET("v1/book/topic/:topicID", GetTopicBooks(storage, validator.ValidateGetTopicBooks(storage)))                 // <GetTopicBooks>     .../v1/book/topic/:topicID
	e.GET("v1/book/lang/:langID", GetLangBooks(storage, validator.ValidateGetLangBooks(storage)))                     // <GetLangBooks>      .../v1/book/lang/:langID

	userGroup.GET("", GetUser(storage, validator.ValidateGetUser(storage)))                                                     // <GetUser>               .../v1/user
	userGroup.DELETE("", DeleteUser(storage, validator.ValidateDeleteUser(storage)))                                            // <DeleteUser>            .../v1/user
	userGroup.PATCH("/password", ChangePassword(storage, validator.ValidateChangePass(storage)))                                // <ChangePassword>        .../v1/user/password
	userGroup.PATCH("/username", ChangeUsername(storage, validator.ValidateChangeUsername(storage)))                            // <ChangeUsername>        .../v1/user/username
	userGroup.POST("/phone", AddPhone(storage, validator.ValidateAddPhone(storage)))                                            // <AddPhone>              .../v1/user/phone
	userGroup.GET("/phone/:phoneID", GetPhone(storage, validator.ValidateGetPhone(storage)))                                    // <GetPhone>              .../v1/user/phone/:phoneID
	userGroup.GET("/phone", GetPhones(storage, validator.ValidateGetPhones(storage)))                                           // <GetPhones>             .../v1/user/phone
	userGroup.DELETE("/phone/:phoneID", DeletePhone(storage, validator.ValidateDeletePhone(storage)))                           // <DeletePhone>           .../v1/user/phone/:phoneID
	userGroup.POST("/address", AddAddress(storage, validator.ValidateAddAddress(storage)))                                      // <AddAddress>            .../v1/user/address
	userGroup.GET("/address/:addressID", GetAddress(storage, validator.ValidateGetAddress(storage)))                            // <GetAddress>            .../v1/user/address/:addressID
	userGroup.GET("/address", GetAddresses(storage, validator.ValidateGetAddresses(storage)))                                   // <GetAddresses>          .../v1/user/address
	userGroup.DELETE("/address/:addressID", DeleteAddress(storage, validator.ValidateDeleteAddress(storage)))                   // <DeleteAddress>         .../v1/user/address/:addressID
	userGroup.POST("/order/item", AddItem(storage, validator.ValidateAddItem(storage)))                                         // <AddItem>               .../v1/user/order/item
	userGroup.PATCH("/order/:orderID/item/:itemID/inc", IncreaseQuantity(storage, validator.ValidateIncreaseQuantity(storage))) // <IncreaseQuantity>      .../v1/user/order/:orderID/item/:itemID/inc
	userGroup.PATCH("/order/:orderID/item/:itemID/dec", DecreaseQuantity(storage, validator.ValidateDecreaseQuantity(storage))) // <DecreaseQuantity>      .../v1/user/order/:orderID/item/:itemID/dec
	userGroup.DELETE("/order/:orderID/item/:itemID", RemoveItem(storage, validator.ValidateRemoveItem(storage)))                // <RemoveItem>            .../v1/user/order/:orderID/item/:itemID
	userGroup.GET("/order/:orderID/item", GetOrderItems(storage, validator.ValidateGetOrderItems(storage)))                     // <GetOrderItems>         .../v1/user/order/:orderID/item
	userGroup.PATCH("/order/:orderID/promo", SetOrderPromo(storage, validator.ValidateSetOrderPromo(storage)))                  // <SetOrderPromo>         .../v1/user/order/:orderID/promo
	userGroup.DELETE("/order/:orderID/promo", RemoveOrderPromo(storage, validator.ValidateRemoveOrderPromo(storage)))           // <RemoveOrderPromo>      .../v1/user/order/:orderID/promo
	userGroup.GET("/order", GetUserOrders(storage, validator.ValidateGetUserOrders(storage)))                                   // <GetUserOrders>         .../v1/user/order
	userGroup.GET("/order/status/:code", GetUserOrdersByStatus(storage, validator.ValidateGetUserOrdersByStatus(storage)))      // <GetUserOrdersByStatus> .../v1/user/order/status/:code
	userGroup.GET("/promo", GetUserPromos(storage, validator.ValidateGetUserPromos(storage)))                                   // <GetUserPromos>         .../v1/user/promo
	userGroup.GET("/dashboard/digital", GetUserDigitalBooks(storage, validator.ValidateGetUserDigitalBooks(storage)))           // <GetUserDigitalBooks>   .../v1/user/dashboard/digital
	userGroup.GET("/dashboard/download/:bookID", DownloadBook(storage, validator.ValidateDownloadBook(storage)))                // <DownloadBook>          .../v1/user/dashboard/download/:bookID
	userGroup.PATCH("/order/:orderID/phone", SetOrderPhone(storage, validator.ValidateSetOrderPhone(storage)))                  // <SetOrderPhone>         .../v1/user/order/:orderID/phone
	userGroup.PATCH("/order/:orderID/address", SetOrderAddress(storage, validator.ValidateSetOrderAddress(storage)))            // <SetOrderAddress>       .../v1/user/order/:orderID/address
	userGroup.DELETE("/logout", UserLogOut(storage))                                                                            // <UserLogOut>            .../v1/logout
	userGroup.DELETE("/logout/all", UserLogOutAllDevices(storage))                                                              // <UserLogOutAllDevices>  .../v1/logout/all

	userGroup.POST("/order/:orderID/payment/zarinpal", payment.ZarinpalPayment(storage, validator.ValidateGetOrderPaymentInfo(storage))) // <ZarinpalPayment>             .../v1/user/order/:orderID/payment/zarinpal
	e.GET("v1/payment/zarinpal/check", payment.ZarinpalPaymentVerification(storage))                                                     // <ZarinpalPaymentVerification> .../v1/payment/zarinpal/check

	adminGroup.GET("/users", GetUsers(storage))                                                                                  // <GetUsers>              .../v1/admin/users
	adminGroup.GET("", GetAdmin(storage, validator.ValidateGetAdmin(storage)))                                                   // <GetAdmin>              .../v1/admin
	adminGroup.GET("s", GetAdmins(storage))                                                                                      // <GetAdmins>             .../v1/admins
	adminGroup.POST("/author", AddAuthor(storage, validator.ValidateAddAuthor(storage)))                                         // <AddAuthor>             .../v1/admin/author
	adminGroup.DELETE("/author/:authorID", DeleteAuthor(storage, validator.ValidateDeleteAuthor(storage)))                       // <DeleteAuthor>          .../v1/admin/author/:authorID
	adminGroup.POST("/publisher", AddPublisher(storage, validator.ValidateAddPublisher(storage)))                                // <AddPublisher>          .../v1/admin/publisher
	adminGroup.DELETE("/publisher/:publisherID", DeletePublisher(storage, validator.ValidateDeletePublisher(storage)))           // <DeltePublisher>        .../v1/admin/publisher/:publisherID
	adminGroup.POST("/topic", AddTopic(storage, validator.ValidateAddTopic(storage)))                                            // <AddTopic>              .../v1/admin/topic
	adminGroup.DELETE("/topic/:topicID", DeleteTopic(storage, validator.ValidateDeleteTopic(storage)))                           // <DeleteTopic>           .../v1/admin/topic/:topicID
	adminGroup.POST("/lang", AddLanguage(storage, validator.ValidateAddLanguage(storage)))                                       // <AddLanguage>           .../v1/admin/lang
	adminGroup.DELETE("/lang/:langID", DeleteLanguage(storage, validator.ValidateDeleteLanguage(storage)))                       // <DeleteLanguage>        .../v1/admin/lang/:langID
	adminGroup.POST("/book", AddBook(storage, validator.ValidateAddBook(storage)))                                               // <AddBook>               .../v1/admin/book
	adminGroup.PATCH("/discount/:bookID", SetBookDiscount(storage, validator.ValidateSetBookDiscount(storage)))                  // <SetBookDiscount>       .../v1/admin/discount/:bookID
	adminGroup.PUT("/book/:bookID", EditBook(storage, validator.ValidateEditBook(storage)))                                      // <EditBook>              .../v1/admin/book/:bookID
	adminGroup.DELETE("/book/:bookID", DeleteBook(storage, validator.ValidateDeleteBook(storage)))                               // <DeleteBook>            .../v1/admin/book/:bookID
	adminGroup.POST("/promo", CreatePromoCode(storage, validator.ValidateCreatePromoCode(storage)))                              // <CreatePromoCode>       .../v1/admin/promo
	adminGroup.DELETE("/promo/:promoID", DeletePromoCode(storage, validator.ValidateDeletePromoCode(storage)))                   // <DeletePromoCode>       .../v1/admin/promo/:promoID
	adminGroup.PATCH("/order/:orderID/status", SetOrderStatus(storage, validator.ValidateSetOrderStatus(storage)))               // <SetOrderStatus>        .../v1/admin/order/:orderID/status
	adminGroup.PATCH("/order/:orderID/stn", SetOrderSTN(storage, validator.ValidateSetOrderSTN(storage)))                        // <SetOrderSTN>           .../v1/admin/order/:orderID/stn
	adminGroup.DELETE("/order/:orderID", DeleteOrder(storage, validator.ValidateDeleteOrder(storage)))                           // <DeleteOrder>           .../v1/admin/order/:orderID
	adminGroup.GET("/order", GetAllOrders(storage))                                                                              // <GetAllOrders>          .../v1/admin/order
	adminGroup.GET("/order/status/:code", GetAllOrdersByStatus(storage, validator.ValidateGetAllOrdersByStatus(storage)))        // <GetAllOrdersByStatus>  .../v1/admin/order/:status
	adminGroup.GET("/order/date", GetDateOrders(storage, validator.ValidateGetDateOrders(storage)))                              // <GetDateOrders>         .../v1/admin/order/date
	adminGroup.GET("/order/date/status/:code", GetDateOrdersByStatus(storage, validator.ValidateGetDateOrdersByStatus(storage))) // <GetDateOrdersByStatus> .../v1/admin/order/date/status/:code
	adminGroup.GET("/promo", GetAllPromos(storage))                                                                              // <GetAllPromos>          .../v1/admin/promo
	adminGroup.GET("/promo/order/:orderID", GetPromoByOrder(storage, validator.ValidateGetPromoByOrder(storage)))                // <GetPromoByOrder>       .../v1/admin/promo/order/:orderID
	adminGroup.DELETE("/logout", AdminLogOut(storage))                                                                           // <AdminLogOut>            .../v1/admin/logout
	adminGroup.DELETE("/logout/all", AdminLogOutAllDevices(storage))                                                             // <AdminLogOutAllDevices>  .../v1/admin/logout/all

	return e
}
