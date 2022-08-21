package order

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type (
	ValidateCreateEmptyOrder func(ctx context.Context, req dto.CreateEmptyOrderRequest) error
	ValidateCheckOpenOrder   func(ctx context.Context, req dto.CheckOpenOrderRequest) error

	ValidateAddItem          func(ctx context.Context, req dto.AddItemRequest) error
	ValidateIncreaseQuantity func(ctx context.Context, req dto.IncreaseQuantityRequest) error
	ValidateDecreaseQuantity func(ctx context.Context, req dto.DecreaseQuantityRequest) error
	ValidateGetOrderItems    func(ctx context.Context, req dto.GetOrderItemsRequest) error
	ValidateRemoveItem       func(ctx context.Context, req dto.RemoveItemRequest) error

	ValidateCreatePromoCode func(ctx context.Context, req dto.CreatePromoCodeRequest) error
	ValidateDeletePromoCode func(ctx context.Context, req dto.DeletePromoCodeRequest) error

	ValidateSetOrderStatus   func(ctx context.Context, req dto.SetOrderStatusRequest) error
	ValidateSetOrderSTN      func(ctx context.Context, req dto.SetOrderSTNRequest) error
	ValidateSetOrderPromo    func(ctx context.Context, req dto.SetOrderPromoRequest) error
	ValidateRemoveOrderPromo func(ctx context.Context, req dto.RemoveOrderPromoRequest) error
	ValidateDeleteOrder      func(ctx context.Context, req dto.DeleteOrderRequest) error

	ValidateGetAllOrdersByStatus  func(ctx context.Context, req dto.GetAllOrdersByStatusRequest) error
	ValidateGetUserOrders         func(ctx context.Context, req dto.GetUserOrdersRequest) error
	ValidateGetUserOrdersByStatus func(ctx context.Context, req dto.GetUserOrdersByStatusRequest) error
	ValidateGetDateOrders         func(ctx context.Context, req dto.GetDateOrdersRequest) error
	ValidateGetDateOrdersByStatus func(ctx context.Context, req dto.GetDateOrdersByStatusRequest) error

	ValidateGetUserPromos   func(ctx context.Context, req dto.GetUserPromosRequest) error
	ValidateGetPromoByOrder func(ctx context.Context, req dto.GetPromoByOrderRequest) error

	ValidateSetOrderPhone   func(ctx context.Context, req dto.SetOrderPhoneRequest) error
	ValidateSetOrderAddress func(ctx context.Context, req dto.SetOrderAddressRequest) error

	ValidateGetOrderPaymentInfo func(ctx context.Context, req dto.GetOrderPaymentInfoRequest) error
)
