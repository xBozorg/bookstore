package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/XBozorg/bookstore/adapter/auth"
	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/order"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func AddItem(storage repository.Storage, validator order.ValidateAddItem) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddItemRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "user does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "user does not exist")
			}

			if err.Error() == "book does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "book does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).AddItem(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "item already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func IncreaseQuantity(storage repository.Storage, validator order.ValidateIncreaseQuantity) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.IncreaseQuantityRequest{}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		iid, err := strconv.ParseUint(c.Param("itemID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.ItemID = uint(iid)

		if err := validator(c.Request().Context(), req); err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "order does not open" {
				return echo.NewHTTPError(http.StatusForbidden, "you don't have any open orders")
			}

			if err.Error() == "item does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "item does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).IncreaseQuantity(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DecreaseQuantity(storage repository.Storage, validator order.ValidateDecreaseQuantity) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DecreaseQuantityRequest{}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		iid, err := strconv.ParseUint(c.Param("itemID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.ItemID = uint(iid)

		if err := validator(c.Request().Context(), req); err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "order does not open" {
				return echo.NewHTTPError(http.StatusForbidden, "you don't have any open orders")
			}

			if err.Error() == "item does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "item does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).DecreaseQuantity(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func RemoveItem(storage repository.Storage, validator order.ValidateRemoveItem) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.RemoveItemRequest{}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		iid, err := strconv.ParseUint(c.Param("itemID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.ItemID = uint(iid)

		if err := validator(c.Request().Context(), req); err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "order does not open" {
				return echo.NewHTTPError(http.StatusForbidden, "you don't have any open orders")
			}

			if err.Error() == "item does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "item does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).RemoveItem(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetOrderItems(storage repository.Storage, validator order.ValidateGetOrderItems) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetOrderItemsRequest{}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "order does not open" {
				return echo.NewHTTPError(http.StatusForbidden, "you don't have any open orders")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetOrderItems(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func CreatePromoCode(storage repository.Storage, validator order.ValidateCreatePromoCode) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.CreatePromoCodeRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "user does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "user does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).CreatePromoCode(c.Request().Context(), req)
		if err != nil {

			if err.Error() == "percentage cannot be 0" {
				return errors.New("percentage cannot be 0")
			}
			if err.Error() == "limit cannot be 0" {
				return errors.New("limit cannot be 0")
			}

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "promo code already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DeletePromoCode(storage repository.Storage, validator order.ValidateDeletePromoCode) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeletePromoCodeRequest{}

		pid, err := strconv.ParseUint(c.Param("promoID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.PromoID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "promo does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "promo does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).DeletePromoCode(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func SetOrderStatus(storage repository.Storage, validator order.ValidateSetOrderStatus) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.SetOrderStatusRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "order does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "order does not exist")
			}

			if err.Error() == "invalid order status" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid order status")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).SetOrderStatus(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func SetOrderSTN(storage repository.Storage, validator order.ValidateSetOrderSTN) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.SetOrderSTNRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "order does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "order does not exist")
			}

			if err.Error() == "invalid order status for stn" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid order status for stn")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).SetOrderSTN(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func SetOrderPromo(storage repository.Storage, validator order.ValidateSetOrderPromo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.SetOrderPromoRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "user does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "user does not exist")
			}

			if err.Error() == "order does not open" {
				return echo.NewHTTPError(http.StatusForbidden, "you don't have any open orders")
			}

			if err.Error() == "promo code does not exist" {
				return echo.NewHTTPError(http.StatusBadRequest, "promo code does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).SetOrderPromo(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func RemoveOrderPromo(storage repository.Storage, validator order.ValidateRemoveOrderPromo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.RemoveOrderPromoRequest{}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "order does not open" {
				return echo.NewHTTPError(http.StatusForbidden, "you don't have any open orders")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).RemoveOrderPromo(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteOrder(storage repository.Storage, validator order.ValidateDeleteOrder) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeleteOrderRequest{}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			if err.Error() == "order does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "order does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).DeleteOrder(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetAllOrders(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAllOrdersRequest{}

		resp, err := order.New(storage).GetAllOrders(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetAllOrdersByStatus(storage repository.Storage, validator order.ValidateGetAllOrdersByStatus) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAllOrdersByStatusRequest{}

		s, err := strconv.ParseUint(c.Param("code"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.Status = uint(s)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "invalid order status" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid order status")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetAllOrdersByStatus(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetUserOrders(storage repository.Storage, validator order.ValidateGetUserOrders) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetUserOrdersRequest{}

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "user does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "user does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetUserOrders(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetUserOrdersByStatus(storage repository.Storage, validator order.ValidateGetUserOrdersByStatus) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetUserOrdersByStatusRequest{}

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		s, err := strconv.ParseUint(c.Param("code"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.Status = uint(s)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "user does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "user does not exist")
			}

			if err.Error() == "invalid order status" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid order status")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetUserOrdersByStatus(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetDateOrders(storage repository.Storage, validator order.ValidateGetDateOrders) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetDateOrdersRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "invalid date" {
				return echo.NewHTTPError(http.StatusNotFound, "invalid date")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetDateOrders(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetDateOrdersByStatus(storage repository.Storage, validator order.ValidateGetDateOrdersByStatus) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetDateOrdersByStatusRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		s, err := strconv.ParseUint(c.Param("code"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.Status = uint(s)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "invalid date" {
				return echo.NewHTTPError(http.StatusNotFound, "invalid date")
			}

			if err.Error() == "invalid order status" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid order status")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetDateOrdersByStatus(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetAllPromos(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAllPromosRequest{}

		resp, err := order.New(storage).GetAllPromos(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetPromoByOrder(storage repository.Storage, validator order.ValidateGetPromoByOrder) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPromoByOrderRequest{}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "order does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "order does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetPromoByOrder(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func GetUserPromos(storage repository.Storage, validator order.ValidateGetUserPromos) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetUserPromosRequest{}

		id, err := auth.GetID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		req.UserID = id

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "user does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "user does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).GetUserPromos(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func SetOrderPhone(storage repository.Storage, validator order.ValidateSetOrderPhone) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.SetOrderPhoneRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "order does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "order does not exist")
			}

			if err.Error() == "phone does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "phone does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).SetOrderPhone(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func SetOrderAddress(storage repository.Storage, validator order.ValidateSetOrderAddress) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.SetOrderAddressRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "order does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "order does not exist")
			}

			if err.Error() == "address does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "address does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := order.New(storage).SetOrderAddress(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}
