package payment

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	repository "github.com/XBozorg/bookstore/adapter/repository/mysql"
	"github.com/XBozorg/bookstore/config"
	"github.com/XBozorg/bookstore/dto"
	orderEntity "github.com/XBozorg/bookstore/entity/order"
	"github.com/XBozorg/bookstore/usecase/order"
	"github.com/labstack/echo/v4"

	"github.com/xbozorg/zarinpal-api"
)

var zarinpalConfig = config.Conf.GetZarinpalConfig()

func ZarinpalPayment(repo repository.MySQLRepo, validator order.ValidateGetOrderPaymentInfo) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := dto.GetOrderPaymentInfoRequest{}
		oid, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid orderID")
		}
		req.OrderID = uint(oid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "order does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "order does not exist")
			}
			if strings.Contains(err.Error(), "order does not open") {
				return echo.NewHTTPError(http.StatusNotFound, "you don't have any open orders")
			}
			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound, "you don't have any order")
			}
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}

		paymentInfo, err := order.New(repo).GetOrderPaymentInfo(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		z := zarinpal.New(zarinpalConfig.MerchantID, zarinpalConfig.Sandbox)

		paymentResponseData, err := z.PaymentRequest(
			zarinpal.PaymentRequest{
				MerchantID:  z.MerchantID,
				Amount:      int(paymentInfo.Total),
				Description: "bookstore",
				CallbackURL: fmt.Sprintf("http://%s/v1/payment/zarinpal/check", c.Request().Host),
				Metadata:    map[string]string{"mobile": paymentInfo.Phone, "email": paymentInfo.Email},
			},
			zarinpal.ValidatePayment(),
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if _, err = order.New(repo).ZarinpalCreateOpenOrder(
			c.Request().Context(),
			dto.ZarinpalCreateOpenOrderRequest{
				OrderID:   req.OrderID,
				Authority: paymentResponseData.Authority,
			},
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.Redirect(
			http.StatusMovedPermanently,
			z.DefaultConfig.PaymentGatewayURL+paymentResponseData.Authority,
		)
	}
}

func ZarinpalPaymentVerification(repo repository.MySQLRepo) echo.HandlerFunc {
	return func(c echo.Context) error {

		authority := c.QueryParam("Authority")
		status := c.QueryParam("Status")

		if err := zarinpal.ValidateGateway(zarinpal.GatewayResponse{
			Status:    status,
			Authority: authority,
		}); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		getOrderResp, err := order.New(repo).ZarinpalGetOrderByAuthority(
			c.Request().Context(),
			dto.ZarinpalGetOrderByAuthorityRequest{Authority: authority},
		)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound, "order payment not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if getOrderResp.ZarinpalOrder.Code == 100 {
			return echo.NewHTTPError(http.StatusConflict, "order payment already verified")
		}

		total, err := order.New(repo).GetOrderTotal(
			c.Request().Context(),
			dto.GetOrderTotalRequest{
				OrderID: getOrderResp.ZarinpalOrder.OrderID,
			},
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		z := zarinpal.New(zarinpalConfig.MerchantID, zarinpalConfig.Sandbox)

		verificationResp, err := z.PaymentVerification(
			zarinpal.PaymentVerificationRequest{
				MerchantID: z.MerchantID,
				Amount:     int(total.Total),
				Authority:  authority,
			},
			zarinpal.ValidatePaymentVerification(),
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		switch verificationResp.Status {
		case 100:
			break
		case 101:
			return echo.NewHTTPError(http.StatusConflict, "order payment already verified")
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "invalid status code")
		}

		if _, err = order.New(repo).ZarinpalSetOrderPayment(
			c.Request().Context(),
			dto.ZarinpalSetOrderPaymentRequest{
				ZarinpalOrderID: getOrderResp.ZarinpalOrder.ID,
				Authority:       authority,
				RefID:           verificationResp.RefID,
				Code:            verificationResp.Status,
			},
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if _, err := order.New(repo).SetOrderStatus(
			c.Request().Context(),
			dto.SetOrderStatusRequest{
				OrderID: getOrderResp.ZarinpalOrder.OrderID,
				Status:  orderEntity.StatusPaid,
			},
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if _, err := order.New(repo).SetOrderReceiptDate(
			c.Request().Context(),
			dto.SetOrderReceiptDateRequest{
				OrderID: getOrderResp.ZarinpalOrder.OrderID,
			},
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, verificationResp)
	}
}
