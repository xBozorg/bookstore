package dto

import (
	"github.com/XBozorg/bookstore/entity/order"
)

type CheckOpenOrderRequest struct {
	UserID string `json:"userID"`
}
type CheckOpenOrderResponse struct {
	OrderID uint `json:"orderID"`
}

type AddItemRequest struct {
	UserID string     `json:"userID"`
	Item   order.Item `json:"item"`
}
type AddItemResponse struct{}

type IncreaseQuantityRequest struct {
	OrderID uint `json:"orderID"`
	ItemID  uint `json:"itemID"`
}
type IncreaseQuantityResponse struct{}

type DecreaseQuantityRequest struct {
	OrderID uint `json:"orderID"`
	ItemID  uint `json:"itemID"`
}
type DecreaseQuantityResponse struct{}

type GetOrderItemsRequest struct {
	OrderID uint `json:"orderID"`
}
type GetOrderItemsResponse struct {
	Items []order.Item `json:"items"`
}

type RemoveItemRequest struct {
	OrderID uint `json:"orderID"`
	ItemID  uint `json:"itemID"`
}
type RemoveItemResponse struct{}

type CreatePromoCodeRequest struct {
	UserID string      `json:"userID"`
	Promo  order.Promo `json:"promo"`
}
type CreatePromoCodeResponse struct{}

type DeletePromoCodeRequest struct {
	PromoID uint `json:"promoID"`
}
type DeletePromoCodeResponse struct{}

type SetOrderStatusRequest struct {
	OrderID uint `json:"orderID"`
	Status  uint `json:"status"`
}
type SetOrderStatusResponse struct{}

type GetOrderStatusRequest struct {
	OrderID uint `json:"orderID"`
}
type GetOrderStatusResponse struct {
	Status uint `json:"status"`
}

type SetOrderSTNRequest struct {
	OrderID uint   `json:"orderID"`
	STN     string `json:"stn"`
}
type SetOrderSTNResponse struct{}

type SetOrderPromoRequest struct {
	UserID    string `json:"userID"`
	OrderID   uint   `json:"orderID"`
	PromoCode string `json:"promoCode"`
}
type SetOrderPromoResponse struct{}

type SetOrderReceiptDateRequest struct {
	OrderID uint `json:"orderID"`
}
type SetOrderReceiptDateResponse struct{}

type RemoveOrderPromoRequest struct {
	OrderID uint `json:"orderID"`
}
type RemoveOrderPromoResponse struct{}

type DeleteOrderRequest struct {
	OrderID uint `json:"orderID"`
}
type DeleteOrderResponse struct{}

type DoesOrderExistRequest struct {
	OrderID uint `json:"orderID"`
}
type DoesOrderExistResponse struct{}

type GetAllOrdersRequest struct{}
type GetAllOrdersResponse struct {
	Orders []order.Order `json:"orders"`
}

type GetAllOrdersByStatusRequest struct {
	Status uint `json:"status"`
}
type GetAllOrdersByStatusResponse struct {
	Orders []order.Order `json:"orders"`
}

type GetUserOrdersRequest struct {
	UserID string `json:"userID"`
}
type GetUserOrdersResponse struct {
	Orders []order.Order `json:"orders"`
}

type GetUserOrdersByStatusRequest struct {
	UserID string `json:"userID"`
	Status uint   `json:"status"`
}
type GetUserOrdersByStatusResponse struct {
	Orders []order.Order `json:"orders"`
}

type GetDateOrdersRequest struct {
	Date string `json:"date"`
}
type GetDateOrdersResponse struct {
	Orders []order.Order `json:"orders"`
}

type GetDateOrdersByStatusRequest struct {
	Date   string `json:"date"`
	Status uint   `json:"status"`
}
type GetDateOrdersByStatusResponse struct {
	Orders []order.Order `json:"orders"`
}

type GetAllPromosRequest struct{}
type GetAllPromosResponse struct {
	Promos []order.Promo `json:"promos"`
}

type GetPromoByOrderRequest struct {
	OrderID uint `json:"orderID"`
}
type GetPromoByOrderResponse struct {
	Promo order.Promo `json:"promo"`
}

type GetUserPromosRequest struct {
	UserID string `json:"userID"`
}
type GetUserPromosResponse struct {
	Promos []order.Promo `json:"promos"`
}

type SetOrderPhoneRequest struct {
	OrderID uint `json:"orderID"`
	PhoneID uint `json:"phoneID"`
}
type SetOrderPhoneResponse struct{}

type SetOrderAddressRequest struct {
	OrderID   uint `json:"orderID"`
	AddressID uint `json:"addressID"`
}
type SetOrderAddressResponse struct{}

type GetOrderPaymentInfoRequest struct {
	OrderID uint `json:"orderID"`
}
type GetOrderPaymentInfoResponse struct {
	Total uint   `json:"total"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type GetOrderTotalRequest struct {
	OrderID uint `json:"orderID"`
}
type GetOrderTotalResponse struct {
	Total uint `json:"total"`
}

type ZarinpalPaymentRequset struct {
	OrderID uint `json:"orderID"`
}

type ZarinpalCreateOpenOrderRequest struct {
	OrderID   uint   `json:"orderID"`
	Authority string `json:"authority"`
}
type ZarinpalCreateOpenOrderResponse struct{}

type ZarinpalGetOrderByAuthorityRequest struct {
	Authority string `json:"authority"`
}
type ZarinpalGetOrderByAuthorityResponse struct {
	ZarinpalOrder order.ZarinpalOrder `json:"zarinpalOrder"`
}

type ZarinpalSetOrderPaymentRequest struct {
	ZarinpalOrderID uint   `json:"zarinpalOrderID"`
	Authority       string `json:"authority"`
	RefID           int    `json:"refID"`
	Code            int    `json:"code"`
}
type ZarinpalSetOrderPaymentResponse struct {
}
