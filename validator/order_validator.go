package validator

import (
	"context"
	"errors"
	"strings"
	"time"

	repository "github.com/XBozorg/bookstore/adapter/repository"
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

func ValidateCreateEmptyOrder(repo repository.Repo) order.ValidateCreateEmptyOrder {
	return func(ctx context.Context, req dto.CreateEmptyOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateCheckOpenOrder(repo repository.Repo) order.ValidateCheckOpenOrder {
	return func(ctx context.Context, req dto.CheckOpenOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateAddItem(repo repository.Repo) order.ValidateAddItem {
	return func(ctx context.Context, req dto.AddItemRequest) error {

		if errUserID := validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		); errUserID != nil {
			return errUserID
		}

		if errItem := validation.ValidateStruct(&req.Item,
			validation.Field(&req.Item.BookID, validation.Required, validation.By(doesBookExist(ctx, repo))),
			validation.Field(&req.Item.Type, validation.NotNil, validation.Min(uint(0)), validation.Max(uint(2))),
			validation.Field(&req.Item.Quantity, validation.Required, validation.Min(uint(0))),
		); errItem != nil {
			return errItem
		}

		return nil
	}
}

func ValidateIncreaseQuantity(repo repository.Repo) order.ValidateIncreaseQuantity {
	return func(ctx context.Context, req dto.IncreaseQuantityRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, repo))),
			validation.Field(&req.ItemID, validation.Required, validation.By(doesItemExist(ctx, repo))),
		)
	}
}

func ValidateDecreaseQuantity(repo repository.Repo) order.ValidateDecreaseQuantity {
	return func(ctx context.Context, req dto.DecreaseQuantityRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, repo))),
			validation.Field(&req.ItemID, validation.Required, validation.By(doesItemExist(ctx, repo))),
		)
	}
}

func ValidateGetOrderItems(repo repository.Repo) order.ValidateGetOrderItems {
	return func(ctx context.Context, req dto.GetOrderItemsRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, repo))),
		)
	}
}

func ValidateGetOrderPaymentInfo(repo repository.Repo) order.ValidateGetOrderPaymentInfo {
	return func(ctx context.Context, req dto.GetOrderPaymentInfoRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo)), validation.By(doesOrderOpen(ctx, repo))),
		)
	}
}

func ValidateRemoveItem(repo repository.Repo) order.ValidateRemoveItem {
	return func(ctx context.Context, req dto.RemoveItemRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, repo))),
			validation.Field(&req.ItemID, validation.Required, validation.By(doesItemExist(ctx, repo))),
		)
	}
}

func ValidateCreatePromoCode(repo repository.Repo) order.ValidateCreatePromoCode {
	return func(ctx context.Context, req dto.CreatePromoCodeRequest) error {

		if errUserID := validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		); errUserID != nil {
			return errUserID
		}

		if errPromo := validation.ValidateStruct(&req.Promo,
			validation.Field(&req.Promo.Code, validation.Required, is.Alphanumeric, validation.Length(3, 20)),
			validation.Field(&req.Promo.Expiration, validation.By(isValidDate(ctx, repo))),
			validation.Field(&req.Promo.Limit, validation.Required, validation.Min(uint(0))),
			validation.Field(&req.Promo.Percentage, validation.Required, validation.Min(uint(0)), validation.Max(uint(100))),
			validation.Field(&req.Promo.MaxPrice, validation.Min(uint(0))),
		); errPromo != nil {
			return errPromo
		}

		return nil
	}
}

func ValidateDeletePromoCode(repo repository.Repo) order.ValidateDeletePromoCode {
	return func(ctx context.Context, req dto.DeletePromoCodeRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.PromoID, validation.Required, validation.By(doesPromoExist(ctx, repo))),
		)
	}
}

func ValidateSetOrderStatus(repo repository.Repo) order.ValidateSetOrderStatus {
	return func(ctx context.Context, req dto.SetOrderStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo))),
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, repo))),
		)
	}
}

func ValidateSetOrderSTN(repo repository.Repo) order.ValidateSetOrderSTN {
	return func(ctx context.Context, req dto.SetOrderSTNRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo))),
			validation.Field(&req.STN, validation.Required, validation.Length(10, 50), validation.By(checkStatusForSTN(ctx, repo, req.OrderID))),
		)
	}
}

func ValidateSetOrderPromo(repo repository.Repo) order.ValidateSetOrderPromo {
	return func(ctx context.Context, req dto.SetOrderPromoRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, repo))),
			validation.Field(&req.PromoCode, validation.Required, is.Alphanumeric, validation.Length(3, 20), validation.By(doesPromoCodeExist(ctx, repo, req.UserID))),
		)
	}
}

func ValidateRemoveOrderPromo(repo repository.Repo) order.ValidateRemoveOrderPromo {
	return func(ctx context.Context, req dto.RemoveOrderPromoRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderOpen(ctx, repo))),
		)
	}
}

func ValidateDeleteOrder(repo repository.Repo) order.ValidateDeleteOrder {
	return func(ctx context.Context, req dto.DeleteOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo))),
		)
	}
}

func ValidateGetAllOrdersByStatus(repo repository.Repo) order.ValidateGetAllOrdersByStatus {
	return func(ctx context.Context, req dto.GetAllOrdersByStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, repo))),
		)
	}
}

func ValidateGetUserOrders(repo repository.Repo) order.ValidateGetUserOrders {
	return func(ctx context.Context, req dto.GetUserOrdersRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateGetUserOrdersByStatus(repo repository.Repo) order.ValidateGetUserOrdersByStatus {
	return func(ctx context.Context, req dto.GetUserOrdersByStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, repo))),
		)
	}
}

func ValidateGetDateOrders(repo repository.Repo) order.ValidateGetDateOrders {
	return func(ctx context.Context, req dto.GetDateOrdersRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Date, validation.Required, validation.By(isValidDate(ctx, repo))),
		)
	}
}

func ValidateGetDateOrdersByStatus(repo repository.Repo) order.ValidateGetDateOrdersByStatus {
	return func(ctx context.Context, req dto.GetDateOrdersByStatusRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Date, validation.Required, validation.By(isValidDate(ctx, repo))),
			validation.Field(&req.Status, validation.Required, validation.By(isValidStatus(ctx, repo))),
		)
	}
}

func ValidateGetUserPromos(repo repository.Repo) order.ValidateGetUserPromos {
	return func(ctx context.Context, req dto.GetUserPromosRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, repo))),
		)
	}
}

func ValidateGetPromoByOrder(repo repository.Repo) order.ValidateGetPromoByOrder {
	return func(ctx context.Context, req dto.GetPromoByOrderRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo))),
		)
	}
}

func ValidateSetOrderPhone(repo repository.Repo) order.ValidateSetOrderPhone {
	return func(ctx context.Context, req dto.SetOrderPhoneRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo))),
			validation.Field(&req.PhoneID, validation.Required, validation.By(doesPhoneExist(ctx, repo))),
		)
	}
}

func ValidateSetOrderAddress(repo repository.Repo) order.ValidateSetOrderAddress {
	return func(ctx context.Context, req dto.SetOrderAddressRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo))),
			validation.Field(&req.AddressID, validation.Required, validation.By(doesAddressExist(ctx, repo))),
		)
	}
}

func ValidateZarinpalPayment(repo repository.Repo) order.ValidateSetOrderAddress {
	return func(ctx context.Context, req dto.SetOrderAddressRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.OrderID, validation.Required, validation.By(doesOrderExist(ctx, repo))),
			validation.Field(&req.AddressID, validation.Required, validation.By(doesAddressExist(ctx, repo))),
		)
	}
}
