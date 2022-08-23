package validator

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/dto"
	eo "github.com/XBozorg/bookstore/entity/order"
	"github.com/XBozorg/bookstore/usecase/order"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func doesItemExist(ctx context.Context, repo order.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		itemID := value.(uint)

		ok, err := repo.DoesItemExist(ctx, itemID)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			return err
		}

		if !ok {
			return errors.New("item does not exist")
		}
		return nil
	}
}

func doesPromoExist(ctx context.Context, repo order.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		promoID := value.(uint)

		ok, err := repo.DoesPromoExist(ctx, promoID)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			return err
		}

		if !ok {
			return errors.New("promo does not exist")
		}
		return nil
	}
}

func doesPromoCodeExist(ctx context.Context, repo order.ValidatorRepo, userID string) validation.RuleFunc {
	return func(value interface{}) error {
		promoCode := value.(string)

		ok, err := repo.DoesPromoCodeExist(ctx, promoCode, userID)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			return err
		}

		if !ok {
			return errors.New("promo code does not exist")
		}
		return nil
	}
}

func doesOrderExist(ctx context.Context, repo order.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		orderID := value.(uint)

		ok, err := repo.DoesOrderExist(ctx, orderID)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			return err
		}

		if !ok {
			return errors.New("order does not exist")
		}
		return nil
	}
}

func doesOrderOpen(ctx context.Context, repo order.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		orderID := value.(uint)

		ok, err := repo.DoesOrderOpen(ctx, orderID)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			return err
		}

		if !ok {
			return errors.New("order does not open")
		}
		return nil
	}
}

func isValidDate(ctx context.Context, repo order.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		date := value.(string)

		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return errors.New("invalid date")
		}

		return nil
	}
}

func isValidStatus(ctx context.Context, repo order.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {

		if status := value.(uint); (status != eo.StatusCreated) && (status != eo.StatusPaid) && (status != eo.StatusVerified) && (status != eo.StatusShipped) {
			return errors.New("invalid order status")
		}

		return nil
	}
}

func checkStatusForSTN(ctx context.Context, repo order.Repository, orderID uint) validation.RuleFunc {
	return func(value interface{}) error {

		status, err := repo.GetOrderStatus(ctx, orderID)
		if err != nil {
			return err
		}

		if (status != eo.StatusVerified) && (status != eo.StatusShipped) {
			return errors.New("invalid order status for stn")
		}
		return nil
	}
}

func ValidateCreateEmptyOrder(storage repository.Storage) order.ValidateCreateEmptyOrder {
	return func(ctx context.Context, req dto.CreateEmptyOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
		)
	}
}

func ValidateCheckOpenOrder(storage repository.Storage) order.ValidateCheckOpenOrder {
	return func(ctx context.Context, req dto.CheckOpenOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
		)
	}
}

func ValidateAddItem(storage repository.Storage) order.ValidateAddItem {
	return func(ctx context.Context, req dto.AddItemRequest) error {

		if errUserID := validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
		); errUserID != nil {
			return errUserID
		}

		if errItem := validation.ValidateStruct(&req.Item,
			validation.Field(&req.Item.BookID, validation.Required, validation.By(doesBookExist(ctx, storage))),
			validation.Field(&req.Item.Type, validation.NotNil, validation.Min(uint(0)), validation.Max(uint(2))),
			validation.Field(&req.Item.Quantity, validation.Required, validation.Min(uint(0))),
		); errItem != nil {
			return errItem
		}

		return nil
	}
}

func ValidateIncreaseQuantity(storage repository.Storage) order.ValidateIncreaseQuantity {
	return func(ctx context.Context, req dto.IncreaseQuantityRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, storage))),
			validation.Field(&req.ItemID, validation.Required, validation.By(doesItemExist(ctx, storage))),
		)
	}
}

func ValidateDecreaseQuantity(storage repository.Storage) order.ValidateDecreaseQuantity {
	return func(ctx context.Context, req dto.DecreaseQuantityRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, storage))),
			validation.Field(&req.ItemID, validation.Required, validation.By(doesItemExist(ctx, storage))),
		)
	}
}

func ValidateGetOrderItems(storage repository.Storage) order.ValidateGetOrderItems {
	return func(ctx context.Context, req dto.GetOrderItemsRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, storage))),
		)
	}
}

func ValidateGetOrderPaymentInfo(storage repository.Storage) order.ValidateGetOrderPaymentInfo {
	return func(ctx context.Context, req dto.GetOrderPaymentInfoRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage)), validation.By(doesOrderOpen(ctx, storage))),
		)
	}
}

func ValidateRemoveItem(storage repository.Storage) order.ValidateRemoveItem {
	return func(ctx context.Context, req dto.RemoveItemRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, storage))),
			validation.Field(&req.ItemID, validation.Required, validation.By(doesItemExist(ctx, storage))),
		)
	}
}

func ValidateCreatePromoCode(storage repository.Storage) order.ValidateCreatePromoCode {
	return func(ctx context.Context, req dto.CreatePromoCodeRequest) error {

		if errUserID := validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
		); errUserID != nil {
			return errUserID
		}

		if errPromo := validation.ValidateStruct(&req.Promo,
			validation.Field(&req.Promo.Code, validation.Required, is.Alphanumeric, validation.Length(3, 20)),
			validation.Field(&req.Promo.Expiration, validation.By(isValidDate(ctx, storage))),
			validation.Field(&req.Promo.Limit, validation.Required, validation.Min(uint(0))),
			validation.Field(&req.Promo.Percentage, validation.Required, validation.Min(uint(0)), validation.Max(uint(100))),
			validation.Field(&req.Promo.MaxPrice, validation.Min(uint(0))),
		); errPromo != nil {
			return errPromo
		}

		return nil
	}
}

func ValidateDeletePromoCode(storage repository.Storage) order.ValidateDeletePromoCode {
	return func(ctx context.Context, req dto.DeletePromoCodeRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.PromoID, validation.Required, validation.By(doesPromoExist(ctx, storage))),
		)
	}
}

func ValidateSetOrderStatus(storage repository.Storage) order.ValidateSetOrderStatus {
	return func(ctx context.Context, req dto.SetOrderStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage))),
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, storage))),
		)
	}
}

func ValidateSetOrderSTN(storage repository.Storage) order.ValidateSetOrderSTN {
	return func(ctx context.Context, req dto.SetOrderSTNRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage))),
			validation.Field(&req.STN, validation.Required, validation.Length(10, 50), validation.By(checkStatusForSTN(ctx, storage, req.OrderID))),
		)
	}
}

func ValidateSetOrderPromo(storage repository.Storage) order.ValidateSetOrderPromo {
	return func(ctx context.Context, req dto.SetOrderPromoRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, storage))),
			validation.Field(&req.PromoCode, validation.Required, is.Alphanumeric, validation.Length(3, 20), validation.By(doesPromoCodeExist(ctx, storage, req.UserID))),
		)
	}
}

func ValidateRemoveOrderPromo(storage repository.Storage) order.ValidateRemoveOrderPromo {
	return func(ctx context.Context, req dto.RemoveOrderPromoRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, storage))),
		)
	}
}

func ValidateDeleteOrder(storage repository.Storage) order.ValidateDeleteOrder {
	return func(ctx context.Context, req dto.DeleteOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage))),
		)
	}
}

func ValidateGetAllOrdersByStatus(storage repository.Storage) order.ValidateGetAllOrdersByStatus {
	return func(ctx context.Context, req dto.GetAllOrdersByStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, storage))),
		)
	}
}

func ValidateGetUserOrders(storage repository.Storage) order.ValidateGetUserOrders {
	return func(ctx context.Context, req dto.GetUserOrdersRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
		)
	}
}

func ValidateGetUserOrdersByStatus(storage repository.Storage) order.ValidateGetUserOrdersByStatus {
	return func(ctx context.Context, req dto.GetUserOrdersByStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, storage))),
		)
	}
}

func ValidateGetDateOrders(storage repository.Storage) order.ValidateGetDateOrders {
	return func(ctx context.Context, req dto.GetDateOrdersRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Date, validation.Required, validation.By(isValidDate(ctx, storage))),
		)
	}
}

func ValidateGetDateOrdersByStatus(storage repository.Storage) order.ValidateGetDateOrdersByStatus {
	return func(ctx context.Context, req dto.GetDateOrdersByStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Date, validation.Required, validation.By(isValidDate(ctx, storage))),
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, storage))),
		)
	}
}

func ValidateGetUserPromos(storage repository.Storage) order.ValidateGetUserPromos {
	return func(ctx context.Context, req dto.GetUserPromosRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
		)
	}
}

func ValidateGetPromoByOrder(storage repository.Storage) order.ValidateGetPromoByOrder {
	return func(ctx context.Context, req dto.GetPromoByOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage))),
		)
	}
}

func ValidateSetOrderPhone(storage repository.Storage) order.ValidateSetOrderPhone {
	return func(ctx context.Context, req dto.SetOrderPhoneRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage))),
			validation.Field(&req.PhoneID, validation.Required, validation.By(doesPhoneExist(ctx, storage))),
		)
	}
}

func ValidateSetOrderAddress(storage repository.Storage) order.ValidateSetOrderAddress {
	return func(ctx context.Context, req dto.SetOrderAddressRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage))),
			validation.Field(&req.AddressID, validation.Required, validation.By(doesAddressExist(ctx, storage))),
		)
	}
}

func ValidateZarinpalPayment(storage repository.Storage) order.ValidateSetOrderAddress {
	return func(ctx context.Context, req dto.SetOrderAddressRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, storage))),
			validation.Field(&req.AddressID, validation.Required, validation.By(doesAddressExist(ctx, storage))),
		)
	}
}
