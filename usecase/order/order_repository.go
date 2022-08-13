package order

import (
	"context"

	"github.com/XBozorg/bookstore/entity/order"
)

type Repository interface {
	AddItem(ctx context.Context, item order.Item, userID string) error
	IncreaseQuantity(ctx context.Context, itemID, orderID uint) error
	DecreaseQuantity(ctx context.Context, itemID, orderID uint) error
	GetOrderItems(ctx context.Context, orderID uint) ([]order.Item, error)
	RemoveItem(ctx context.Context, itemID, orderID uint) error
	CheckQuantity(ctx context.Context, quantity, bookID uint) error

	CreatePromoCode(ctx context.Context, promo order.Promo, userID string) error
	DeletePromoCode(ctx context.Context, promoID uint) error

	CreateEmptyOrder(ctx context.Context, userID string) (uint, error)
	CheckOpenOrder(ctx context.Context, userID string) (uint, error)
	SetOrderStatus(ctx context.Context, status, orderID uint) error
	GetOrderStatus(ctx context.Context, orderID uint) (uint, error)
	SetOrderSTN(ctx context.Context, stn string, orderID uint) error
	SetOrderTotal(ctx context.Context, orderID uint) error
	SetOrderPromo(ctx context.Context, orderID uint, promoCode, userID string) error
	RemoveOrderPromo(ctx context.Context, orderID uint) error
	DeleteOrder(ctx context.Context, orderID uint) error

	GetAllOrders(ctx context.Context) ([]order.Order, error)
	GetAllOrdersByStatus(ctx context.Context, status uint) ([]order.Order, error)
	GetUserOrders(ctx context.Context, userID string) ([]order.Order, error)
	GetUserOrdersByStatus(ctx context.Context, userID string, status uint) ([]order.Order, error)
	GetDateOrders(ctx context.Context, date string) ([]order.Order, error)
	GetDateOrdersByStatus(ctx context.Context, date string, status uint) ([]order.Order, error)

	GetAllPromos(ctx context.Context) ([]order.Promo, error)
	GetPromoByOrder(ctx context.Context, orderID uint) (order.Promo, error)
	GetUserPromos(ctx context.Context, userID string) ([]order.Promo, error)
}

type ValidatorRepo interface {
	DoesItemExist(ctx context.Context, itemID uint) (bool, error)
	DoesPromoExist(ctx context.Context, promoID uint) (bool, error)
	DoesPromoCodeExist(ctx context.Context, promoCode, userID string) (bool, error)
	DoesOrderExist(ctx context.Context, orderID uint) (bool, error)
	DoesOrderOpen(ctx context.Context, orderID uint) (bool, error)
}
