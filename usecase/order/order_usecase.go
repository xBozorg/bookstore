package order

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type UseCase interface {
	AddItem(ctx context.Context, req dto.AddItemRequest) (dto.AddItemResponse, error)
	IncreaseQuantity(ctx context.Context, req dto.IncreaseQuantityRequest) (dto.IncreaseQuantityResponse, error)
	DecreaseQuantity(ctx context.Context, req dto.DecreaseQuantityRequest) (dto.DecreaseQuantityResponse, error)
	GetOrderItems(ctx context.Context, req dto.GetOrderItemsRequest) (dto.GetOrderItemsResponse, error)
	RemoveItem(ctx context.Context, req dto.RemoveItemRequest) (dto.RemoveItemResponse, error)

	SetOrderPhone(ctx context.Context, req dto.SetOrderPhoneRequest) (dto.SetOrderPhoneResponse, error)
	SetOrderAddress(ctx context.Context, req dto.SetOrderAddressRequest) (dto.SetOrderAddressResponse, error)

	CreatePromoCode(ctx context.Context, req dto.CreatePromoCodeRequest) (dto.CreatePromoCodeResponse, error)
	DeletePromoCode(ctx context.Context, req dto.CreatePromoCodeRequest) (dto.CreatePromoCodeResponse, error)

	SetOrderStatus(ctx context.Context, req dto.SetOrderStatusRequest) (dto.SetOrderStatusResponse, error)
	GetOrderStatus(ctx context.Context, req dto.GetOrderStatusRequest) (dto.GetOrderStatusResponse, error)
	SetOrderSTN(ctx context.Context, req dto.SetOrderSTNRequest) (dto.SetOrderSTNResponse, error)
	SetOrderPromo(ctx context.Context, req dto.SetOrderPromoRequest) (dto.SetOrderPromoResponse, error)
	SetOrderReceiptDate(ctx context.Context, req dto.SetOrderReceiptDateRequest) (dto.SetOrderReceiptDateResponse, error)
	RemoveOrderPromo(ctx context.Context, req dto.RemoveOrderPromoRequest) (dto.RemoveOrderPromoResponse, error)
	DeleteOrder(ctx context.Context, req dto.DeleteOrderRequest) (dto.DeleteOrderResponse, error)

	GetAllOrders(ctx context.Context, req dto.GetAllOrdersRequest) (dto.GetAllOrdersResponse, error)
	GetAllOrdersByStatus(ctx context.Context, req dto.GetAllOrdersByStatusRequest) (dto.GetAllOrdersByStatusResponse, error)
	GetUserOrders(ctx context.Context, req dto.GetUserOrdersRequest) (dto.GetUserOrdersResponse, error)
	GetUserOrdersByStatus(ctx context.Context, req dto.GetUserOrdersByStatusRequest) (dto.GetUserOrdersByStatusResponse, error)
	GetDateOrders(ctx context.Context, req dto.GetDateOrdersRequest) (dto.GetDateOrdersResponse, error)
	GetDateOrdersByStatus(ctx context.Context, req dto.GetDateOrdersByStatusRequest) (dto.GetDateOrdersByStatusResponse, error)

	GetAllPromos(ctx context.Context, req dto.GetAllPromosRequest) (dto.GetAllPromosResponse, error)
	GetPromoByOrder(ctx context.Context, req dto.GetPromoByOrderRequest) (dto.GetPromoByOrderResponse, error)
	GetUserPromos(ctx context.Context, req dto.GetUserPromosRequest) (dto.GetUserPromosResponse, error)

	GetOrderPaymentInfo(ctx context.Context, req dto.GetOrderPaymentInfoRequest) (dto.GetOrderPaymentInfoResponse, error)
	GetOrderTotal(ctx context.Context, req dto.GetOrderTotalRequest) (dto.GetOrderTotalResponse, error)

	ZarinpalCreateOpenOrder(ctx context.Context, req dto.ZarinpalCreateOpenOrderRequest) (dto.ZarinpalCreateOpenOrderResponse, error)
	ZarinpalGetOrderByAuthority(ctx context.Context, req dto.ZarinpalGetOrderByAuthorityRequest) (dto.ZarinpalGetOrderByAuthorityResponse, error)
	ZarinpalSetOrderPayment(ctx context.Context, req dto.ZarinpalSetOrderPaymentRequest) (dto.ZarinpalSetOrderPaymentResponse, error)
}

type UseCaseRepo struct {
	repo Repository
}

func New(r Repository) UseCaseRepo {
	return UseCaseRepo{repo: r}
}

func (u UseCaseRepo) AddItem(ctx context.Context, req dto.AddItemRequest) (dto.AddItemResponse, error) {

	err := u.repo.AddItem(ctx, req.Item, req.UserID)
	if err != nil {
		return dto.AddItemResponse{}, err
	}

	return dto.AddItemResponse{}, nil
}

func (u UseCaseRepo) IncreaseQuantity(ctx context.Context, req dto.IncreaseQuantityRequest) (dto.IncreaseQuantityResponse, error) {

	err := u.repo.IncreaseQuantity(ctx, req.ItemID, req.OrderID)
	if err != nil {
		return dto.IncreaseQuantityResponse{}, err
	}

	return dto.IncreaseQuantityResponse{}, nil
}

func (u UseCaseRepo) DecreaseQuantity(ctx context.Context, req dto.DecreaseQuantityRequest) (dto.DecreaseQuantityResponse, error) {

	err := u.repo.DecreaseQuantity(ctx, req.ItemID, req.OrderID)
	if err != nil {
		return dto.DecreaseQuantityResponse{}, err
	}

	return dto.DecreaseQuantityResponse{}, nil
}

func (u UseCaseRepo) GetOrderItems(ctx context.Context, req dto.GetOrderItemsRequest) (dto.GetOrderItemsResponse, error) {

	items, err := u.repo.GetOrderItems(ctx, req.OrderID)
	if err != nil {
		return dto.GetOrderItemsResponse{}, err
	}

	return dto.GetOrderItemsResponse{Items: items}, nil
}

func (u UseCaseRepo) RemoveItem(ctx context.Context, req dto.RemoveItemRequest) (dto.RemoveItemResponse, error) {

	err := u.repo.RemoveItem(ctx, req.ItemID, req.OrderID)
	if err != nil {
		return dto.RemoveItemResponse{}, err
	}

	return dto.RemoveItemResponse{}, nil
}

func (u UseCaseRepo) CreatePromoCode(ctx context.Context, req dto.CreatePromoCodeRequest) (dto.CreatePromoCodeResponse, error) {

	err := u.repo.CreatePromoCode(ctx, req.Promo, req.UserID)
	if err != nil {
		return dto.CreatePromoCodeResponse{}, err
	}

	return dto.CreatePromoCodeResponse{}, nil
}

func (u UseCaseRepo) DeletePromoCode(ctx context.Context, req dto.DeletePromoCodeRequest) (dto.DeletePromoCodeResponse, error) {

	err := u.repo.DeletePromoCode(ctx, req.PromoID)
	if err != nil {
		return dto.DeletePromoCodeResponse{}, err
	}

	return dto.DeletePromoCodeResponse{}, nil
}

func (u UseCaseRepo) SetOrderStatus(ctx context.Context, req dto.SetOrderStatusRequest) (dto.SetOrderStatusResponse, error) {

	err := u.repo.SetOrderStatus(ctx, req.Status, req.OrderID)
	if err != nil {
		return dto.SetOrderStatusResponse{}, err
	}

	return dto.SetOrderStatusResponse{}, nil
}

func (u UseCaseRepo) GetOrderStatus(ctx context.Context, req dto.GetOrderStatusRequest) (dto.GetOrderStatusResponse, error) {

	status, err := u.repo.GetOrderStatus(ctx, req.OrderID)
	if err != nil {
		return dto.GetOrderStatusResponse{}, err
	}

	return dto.GetOrderStatusResponse{Status: status}, nil
}

func (u UseCaseRepo) SetOrderSTN(ctx context.Context, req dto.SetOrderSTNRequest) (dto.SetOrderSTNResponse, error) {

	err := u.repo.SetOrderSTN(ctx, req.STN, req.OrderID)
	if err != nil {
		return dto.SetOrderSTNResponse{}, err
	}

	return dto.SetOrderSTNResponse{}, nil
}

func (u UseCaseRepo) SetOrderPromo(ctx context.Context, req dto.SetOrderPromoRequest) (dto.SetOrderPromoResponse, error) {

	err := u.repo.SetOrderPromo(ctx, req.OrderID, req.PromoCode, req.UserID)
	if err != nil {
		return dto.SetOrderPromoResponse{}, err
	}

	return dto.SetOrderPromoResponse{}, nil
}

func (u UseCaseRepo) SetOrderReceiptDate(ctx context.Context, req dto.SetOrderReceiptDateRequest) (dto.SetOrderReceiptDateResponse, error) {

	err := u.repo.SetOrderReceiptDate(ctx, req.OrderID)
	if err != nil {
		return dto.SetOrderReceiptDateResponse{}, err
	}

	return dto.SetOrderReceiptDateResponse{}, nil
}

func (u UseCaseRepo) RemoveOrderPromo(ctx context.Context, req dto.RemoveOrderPromoRequest) (dto.RemoveOrderPromoResponse, error) {

	err := u.repo.RemoveOrderPromo(ctx, req.OrderID)
	if err != nil {
		return dto.RemoveOrderPromoResponse{}, err
	}

	return dto.RemoveOrderPromoResponse{}, nil
}

func (u UseCaseRepo) DeleteOrder(ctx context.Context, req dto.DeleteOrderRequest) (dto.DeleteOrderResponse, error) {

	err := u.repo.DeleteOrder(ctx, req.OrderID)
	if err != nil {
		return dto.DeleteOrderResponse{}, err
	}

	return dto.DeleteOrderResponse{}, nil
}

func (u UseCaseRepo) GetAllOrders(ctx context.Context, req dto.GetAllOrdersRequest) (dto.GetAllOrdersResponse, error) {

	orders, err := u.repo.GetAllOrders(ctx)
	if err != nil {
		return dto.GetAllOrdersResponse{}, err
	}

	return dto.GetAllOrdersResponse{Orders: orders}, nil
}

func (u UseCaseRepo) GetAllOrdersByStatus(ctx context.Context, req dto.GetAllOrdersByStatusRequest) (dto.GetAllOrdersByStatusResponse, error) {

	orders, err := u.repo.GetAllOrdersByStatus(ctx, req.Status)
	if err != nil {
		return dto.GetAllOrdersByStatusResponse{}, err
	}

	return dto.GetAllOrdersByStatusResponse{Orders: orders}, nil
}

func (u UseCaseRepo) GetUserOrders(ctx context.Context, req dto.GetUserOrdersRequest) (dto.GetUserOrdersResponse, error) {

	orders, err := u.repo.GetUserOrders(ctx, req.UserID)
	if err != nil {
		return dto.GetUserOrdersResponse{}, err
	}

	return dto.GetUserOrdersResponse{Orders: orders}, nil
}

func (u UseCaseRepo) GetUserOrdersByStatus(ctx context.Context, req dto.GetUserOrdersByStatusRequest) (dto.GetUserOrdersByStatusResponse, error) {

	orders, err := u.repo.GetUserOrdersByStatus(ctx, req.UserID, req.Status)
	if err != nil {
		return dto.GetUserOrdersByStatusResponse{}, err
	}

	return dto.GetUserOrdersByStatusResponse{Orders: orders}, nil
}

func (u UseCaseRepo) GetDateOrders(ctx context.Context, req dto.GetDateOrdersRequest) (dto.GetDateOrdersResponse, error) {

	orders, err := u.repo.GetDateOrders(ctx, req.Date)
	if err != nil {
		return dto.GetDateOrdersResponse{}, err
	}

	return dto.GetDateOrdersResponse{Orders: orders}, nil
}

func (u UseCaseRepo) GetDateOrdersByStatus(ctx context.Context, req dto.GetDateOrdersByStatusRequest) (dto.GetDateOrdersByStatusResponse, error) {

	orders, err := u.repo.GetDateOrdersByStatus(ctx, req.Date, req.Status)
	if err != nil {
		return dto.GetDateOrdersByStatusResponse{}, err
	}

	return dto.GetDateOrdersByStatusResponse{Orders: orders}, nil
}

func (u UseCaseRepo) GetAllPromos(ctx context.Context, req dto.GetAllPromosRequest) (dto.GetAllPromosResponse, error) {

	promos, err := u.repo.GetAllPromos(ctx)
	if err != nil {
		return dto.GetAllPromosResponse{}, err
	}

	return dto.GetAllPromosResponse{Promos: promos}, nil
}

func (u UseCaseRepo) GetUserPromos(ctx context.Context, req dto.GetUserPromosRequest) (dto.GetUserPromosResponse, error) {

	promos, err := u.repo.GetUserPromos(ctx, req.UserID)
	if err != nil {
		return dto.GetUserPromosResponse{}, err
	}

	return dto.GetUserPromosResponse{Promos: promos}, nil
}

func (u UseCaseRepo) GetPromoByOrder(ctx context.Context, req dto.GetPromoByOrderRequest) (dto.GetPromoByOrderResponse, error) {

	promo, err := u.repo.GetPromoByOrder(ctx, req.OrderID)
	if err != nil {
		return dto.GetPromoByOrderResponse{}, err
	}

	return dto.GetPromoByOrderResponse{Promo: promo}, nil
}

func (u UseCaseRepo) SetOrderPhone(ctx context.Context, req dto.SetOrderPhoneRequest) (dto.SetOrderPhoneResponse, error) {

	err := u.repo.SetOrderPhone(ctx, req.OrderID, req.PhoneID)
	if err != nil {
		return dto.SetOrderPhoneResponse{}, err
	}

	return dto.SetOrderPhoneResponse{}, nil
}

func (u UseCaseRepo) SetOrderAddress(ctx context.Context, req dto.SetOrderAddressRequest) (dto.SetOrderAddressResponse, error) {

	err := u.repo.SetOrderAddress(ctx, req.OrderID, req.AddressID)
	if err != nil {
		return dto.SetOrderAddressResponse{}, err
	}

	return dto.SetOrderAddressResponse{}, nil
}

func (u UseCaseRepo) GetOrderPaymentInfo(ctx context.Context, req dto.GetOrderPaymentInfoRequest) (dto.GetOrderPaymentInfoResponse, error) {

	info, err := u.repo.GetOrderPaymentInfo(ctx, req.OrderID)
	if err != nil {
		return dto.GetOrderPaymentInfoResponse{}, err
	}

	return dto.GetOrderPaymentInfoResponse{
		Total: info.Total,
		Email: info.Email,
		Phone: info.Phone,
	}, nil
}

func (u UseCaseRepo) GetOrderTotal(ctx context.Context, req dto.GetOrderTotalRequest) (dto.GetOrderTotalResponse, error) {

	total, err := u.repo.GetOrderTotal(ctx, req.OrderID)
	if err != nil {
		return dto.GetOrderTotalResponse{}, err
	}

	return dto.GetOrderTotalResponse{Total: total}, nil
}

func (u UseCaseRepo) ZarinpalCreateOpenOrder(ctx context.Context, req dto.ZarinpalCreateOpenOrderRequest) (dto.ZarinpalCreateOpenOrderResponse, error) {

	err := u.repo.ZarinpalCreateOpenOrder(ctx, req.OrderID, req.Authority)
	if err != nil {
		return dto.ZarinpalCreateOpenOrderResponse{}, err
	}

	return dto.ZarinpalCreateOpenOrderResponse{}, nil
}

func (u UseCaseRepo) ZarinpalGetOrderByAuthority(ctx context.Context, req dto.ZarinpalGetOrderByAuthorityRequest) (dto.ZarinpalGetOrderByAuthorityResponse, error) {

	zarinpalOrder, err := u.repo.ZarinpalGetOrderByAuthority(ctx, req.Authority)
	if err != nil {
		return dto.ZarinpalGetOrderByAuthorityResponse{}, err
	}

	return dto.ZarinpalGetOrderByAuthorityResponse{ZarinpalOrder: zarinpalOrder}, nil
}

func (u UseCaseRepo) ZarinpalSetOrderPayment(ctx context.Context, req dto.ZarinpalSetOrderPaymentRequest) (dto.ZarinpalSetOrderPaymentResponse, error) {

	err := u.repo.ZarinpalSetOrderPayment(ctx, req.ZarinpalOrderID, req.Authority, req.RefID, req.Code)
	if err != nil {
		return dto.ZarinpalSetOrderPaymentResponse{}, err
	}

	return dto.ZarinpalSetOrderPaymentResponse{}, nil
}
