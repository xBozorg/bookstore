package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/XBozorg/bookstore/entity/book"
	"github.com/XBozorg/bookstore/entity/order"
)

type itemPrice struct {
	BookID           uint
	Type             uint
	Quantity         uint
	DigitalPrice     uint
	DigitalDiscount  uint
	PhysicalPrice    uint
	PhysicalDiscount uint
}

func (m MySQLRepo) DoesOrderOpen(ctx context.Context, orderID uint) (bool, error) {

	result := m.db.QueryRowContext(ctx, "SELECT 1 FROM orders WHERE id = ? AND status = ?", orderID, order.StatusCreated)

	var open bool
	err := result.Scan(&open)
	if err != nil {
		return false, err
	}

	return open, err
}

func (m MySQLRepo) DoesOrderExist(ctx context.Context, orderID uint) (bool, error) {

	result := m.db.QueryRowContext(ctx, "SELECT 1 FROM orders WHERE id = ?", orderID)
	var exists bool
	err := result.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}

func (m MySQLRepo) DoesPromoExist(ctx context.Context, promoID uint) (bool, error) {

	result := m.db.QueryRowContext(ctx, "SELECT 1 FROM promo WHERE id = ?", promoID)
	var exists bool
	err := result.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}

func (m MySQLRepo) DoesPromoCodeExist(ctx context.Context, promoCode, userID string) (bool, error) {

	result := m.db.QueryRowContext(ctx,
		`SELECT 1 FROM promo 
		WHERE id IN (SELECT promo_id FROM promo_user WHERE user_id = ?) AND promo.code = ?`, userID, promoCode)

	var exist bool
	err := result.Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (m MySQLRepo) DoesItemExist(ctx context.Context, itemID uint) (bool, error) {

	result := m.db.QueryRowContext(ctx, "SELECT 1 FROM item WHERE id = ?", itemID)
	var exists bool
	err := result.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}

func (m MySQLRepo) CreateEmptyOrder(ctx context.Context, userID string) (uint, error) {

	result, err := m.db.ExecContext(ctx, `INSERT INTO orders 
    (creation_date , status , total , user_id) 
    VALUES (?,?,?,?)`, time.Now().Format("2006-01-02 15:04:05"), order.StatusCreated, 0, userID)

	if err != nil {
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(orderID), nil
}

func (m MySQLRepo) CheckOpenOrder(ctx context.Context, userID string) (uint, error) {

	result := m.db.QueryRowContext(ctx, "SELECT id FROM orders WHERE user_id=? AND status=?", userID, order.StatusCreated)

	var id uint
	err := result.Scan(&id)

	if err != nil {

		if strings.Contains(err.Error(), "no rows") {
			id, errC := m.CreateEmptyOrder(ctx, userID)
			if errC != nil {
				return 0, errC
			}
			return id, nil
		}

		return 0, err
	}

	return id, nil
}

func (m MySQLRepo) AddItem(ctx context.Context, item order.Item, userID string) error {

	err := m.CheckQuantity(ctx, item.Quantity, item.BookID)
	if err != nil {
		return err
	}

	availability, err := m.CheckAvailability(ctx, item.BookID)
	if err != nil {
		return err
	}

	orderID, err := m.CheckOpenOrder(ctx, userID)
	if err != nil {
		return err
	}

	switch {
	case item.Type == order.Bundle && availability == book.BundleAvailable:

		_, err = m.db.ExecContext(ctx,
			`INSERT INTO item 
    				(book_id , type , quantity , order_id) 
    				VALUES (?,?,?,?) , (?,?,?,?)`,
			item.BookID, order.Digital, 1, orderID, // Digital
			item.BookID, order.Physical, item.Quantity, orderID, // Physical
		)
		if err != nil {
			return err
		}

		err := m.DecreasePhysicalStock(ctx, item.Quantity, item.BookID)
		if err != nil {
			return err
		}

	case item.Type == order.Physical && availability == book.PhysicalAvailable:

		_, err = m.db.ExecContext(ctx,
			`INSERT INTO item 
    				(book_id , type , quantity , order_id) 
    				VALUES (?,?,?,?)`,
			item.BookID, order.Physical, item.Quantity, orderID,
		)
		if err != nil {
			return err
		}

		err := m.DecreasePhysicalStock(ctx, item.Quantity, item.BookID)
		if err != nil {
			return err
		}

	case item.Type == order.Digital && availability == book.DigitalAvailable:
		_, err = m.db.ExecContext(ctx,
			`INSERT INTO item 
				(book_id , type , quantity , order_id) 
				VALUES (?,?,?,?)`,
			item.BookID, order.Digital, 1, orderID,
		)

		if err != nil {
			return err
		}

	case availability == 0:
		return errors.New("item unavailable")
	case item.Type > 2:
		return errors.New("invalid item type")
	case availability > 3:
		return errors.New("invalid item availability")
	default:
		return errors.New("type / availability does not match")
	}

	err = m.SetOrderTotal(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) GetOrderItems(ctx context.Context, orderID uint) ([]order.Item, error) {

	result, err := m.db.QueryContext(ctx, "SELECT id,book_id,type,quantity FROM item WHERE order_id = ?", orderID)
	if err != nil {
		return []order.Item{}, err
	}

	items := []order.Item{}
	for result.Next() {
		var i order.Item
		err := result.Scan(&i.ID, &i.BookID, &i.Type, &i.Quantity)
		if err != nil {
			return []order.Item{}, err
		}

		items = append(items, i)
	}

	return items, nil
}

func (m MySQLRepo) CheckQuantity(ctx context.Context, quantity, bookID uint) error {
	var stock uint
	result := m.db.QueryRowContext(ctx, "SELECT physical_stock FROM book WHERE id = ?", bookID)
	err := result.Scan(&stock)
	if err != nil {
		return err
	}

	if quantity > stock {
		return errors.New("requested item quantity is bigger than the available stock")
	}
	return nil
}

func (m MySQLRepo) CheckAvailability(ctx context.Context, bookID uint) (uint, error) {

	var availability uint
	result := m.db.QueryRowContext(ctx, "SELECT availability FROM book WHERE id = ?", bookID)
	err := result.Scan(&availability)
	if err != nil {
		return 0, err
	}

	return availability, nil
}

func (m MySQLRepo) SetOrderPhone(ctx context.Context, orderID, phoneID uint) error {

	_, err := m.db.ExecContext(ctx, "UPDATE orders SET phone_id = ? WHERE id = ?", phoneID, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) SetOrderAddress(ctx context.Context, orderID, addressID uint) error {

	_, err := m.db.ExecContext(ctx, "UPDATE orders SET address_id = ? WHERE id = ?", addressID, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) IncreaseQuantity(ctx context.Context, itemID, orderID uint) error {

	_, err := m.db.ExecContext(ctx,
		`UPDATE item SET 
			item.quantity = item.quantity + 1 
		WHERE item.id = ? 
		AND
		item.type != ?
		AND
		item.quantity < (SELECT book.physical_stock FROM book WHERE book.id = item.book_id)`, itemID, order.Digital)

	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx,
		`UPDATE book SET physical_stock = physical_stock - 1 
		WHERE item.id = ? 
		AND 
		book.id = item.book_id`, itemID)

	if err != nil {
		return err
	}

	err = m.SetOrderTotal(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) DecreaseQuantity(ctx context.Context, itemID, orderID uint) error {

	_, err := m.db.ExecContext(ctx,
		`UPDATE item SET 
			quantity = quantity - 1 
		WHERE id = ? 
		AND
		type != ?
		AND 
		quantity > 0`, itemID, order.Digital)

	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx,
		`UPDATE book SET physical_stock = physical_stock + 1 
		WHERE item.id = ? 
		AND 
		book.id = item.book_id`, itemID)

	if err != nil {
		return err
	}

	err = m.SetOrderTotal(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) DecreasePhysicalStock(ctx context.Context, quantity, bookID uint) error {

	_, err := m.db.ExecContext(ctx, "UPDATE book SET physical_stock = physical_stock - ? WHERE id = ?", quantity, bookID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) RemoveItem(ctx context.Context, itemID, orderID uint) error {

	result := m.db.QueryRowContext(ctx, "SELECT type FROM item WHERE id = ?", itemID)
	var itemType uint
	err := result.Scan(&itemType)
	if err != nil {
		return err
	}

	if itemType == order.Physical || itemType == order.Bundle {
		_, err = m.db.ExecContext(ctx,
			`UPDATE book SET physical_stock = physical_stock + (SELECT quantity FROM item WHERE id = ?)
			WHERE 
			book.id = (SELECT book_id FROM item WHERE id = ?)`, itemID, itemID)

		if err != nil {
			return err
		}
	}

	_, err = m.db.ExecContext(ctx, "DELETE FROM item WHERE id = ?", itemID)
	if err != nil {
		return err
	}

	err = m.SetOrderTotal(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) CreatePromoCode(ctx context.Context, promo order.Promo, userID string) error {

	if promo.Percentage == 0 {
		return errors.New("percentage cannot be 0")
	}
	if promo.Limit == 0 {
		return errors.New("limit cannot be 0")
	}

	result, err := m.db.ExecContext(ctx,
		`INSERT INTO promo 
		(code , expiration , promo.limit , percentage , max_price)
		VALUES (?,?,?,?,?)`,
		promo.Code, promo.Expiration, promo.Limit, promo.Percentage, promo.MaxPrice,
	)
	if err != nil {
		return err
	}
	promoID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx, "INSERT INTO promo_user (promo_id , user_id) VALUES (?,?)", promoID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) DeletePromoCode(ctx context.Context, promoID uint) error {

	_, err := m.db.ExecContext(ctx, "DELETE FROM promo WHERE id = ?", promoID)
	if err != nil {
		return err
	}
	/*
		_, err = m.db.ExecContext(ctx, "DELETE FROM promo_user WHERE promo_id = ? AND user_id = ?", promoID, userID)
		if err != nil {
			return err
		}
	*/

	return nil
}

func (m MySQLRepo) SetOrderStatus(ctx context.Context, status, orderID uint) error {

	var isShipmentOrder bool
	result := m.db.QueryRowContext(ctx, `SELECT 1 FROM item WHERE type != 0 AND order_id = ?`, orderID)
	err := result.Scan(&isShipmentOrder)
	if err != nil {
		return err
	}

	if status != order.StatusCreated && isShipmentOrder {

		_, err := m.db.ExecContext(ctx,
			`UPDATE orders SET status = ? 
			WHERE id = ? 
			AND 
			(SELECT 1 FROM orders WHERE phone_id IS NOT NULL AND address_id IS NOT NULL)`,
			status, orderID)

		if err != nil {
			return err
		}

	} else {

		_, err := m.db.ExecContext(ctx, "UPDATE orders SET status = ? WHERE id = ?", status, orderID)

		if err != nil {
			return err
		}

	}

	return nil
}

func (m MySQLRepo) GetOrderStatus(ctx context.Context, orderID uint) (uint, error) {

	result := m.db.QueryRowContext(ctx, "SELECT status FROM orders WHERE id = ?", orderID)

	var status uint
	err := result.Scan(&status)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (m MySQLRepo) SetOrderSTN(ctx context.Context, stn string, orderID uint) error {

	_, err := m.db.ExecContext(ctx, "UPDATE orders SET stn = ? WHERE id = ?", stn, orderID)

	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) SetOrderTotal(ctx context.Context, orderID uint) error {

	result, err := m.db.QueryContext(ctx, "SELECT book_id , type , quantity FROM item WHERE order_id = ?", orderID)
	if err != nil {
		return err
	}

	items := []itemPrice{}
	for result.Next() {
		var i itemPrice
		err := result.Scan(&i.BookID, &i.Type, &i.Quantity)
		if err != nil {
			return err
		}

		priceResult := m.db.QueryRowContext(ctx, `SELECT digital_price,digital_discount,
		physical_price,physical_discount FROM book WHERE id = ?`, i.BookID)

		err = priceResult.Scan(&i.DigitalPrice, &i.DigitalDiscount, &i.PhysicalPrice, &i.PhysicalDiscount)
		if err != nil {
			return err
		}
		items = append(items, i)
	}

	total, err := m.CalculateTotal(ctx, items)
	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx, "UPDATE orders SET total = ? WHERE id = ?", total, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) CalculateTotal(ctx context.Context, i []itemPrice) (uint, error) {

	var total uint = 0
	for _, item := range i {
		switch item.Type {

		case order.Digital:
			total += (item.DigitalPrice * (100 - item.DigitalDiscount)) / 100

		case order.Physical:
			total += (item.PhysicalPrice * (100 - item.PhysicalDiscount)) / 100 * item.Quantity

		default:
			return 0, errors.New("invalid order type")
		}
	}
	return total, nil
}

func (m MySQLRepo) SetOrderPromo(ctx context.Context, orderID uint, promoCode, userID string) error {

	result := m.db.QueryRowContext(ctx,
		`SELECT * FROM promo 
		WHERE code = ? 
		AND 
		(SELECT 1 FROM promo_user WHERE promo_id = promo.id AND user_id = ?)`,
		promoCode, userID)

	var promo order.Promo

	err := result.Scan(&promo.ID, &promo.Code, &promo.Expiration, &promo.Limit, &promo.Percentage, &promo.MaxPrice)
	if err != nil {
		return err
	}

	if promo.Percentage == 0 {
		return errors.New("percentage cannot be 0")
	}

	exp, err := time.Parse("2006-01-02 15:04:05", promo.Expiration)
	if err != nil {
		return err
	}

	if expired := time.Now().After(exp); expired {
		return errors.New("expired promo code")
	}

	if promo.Limit == 0 {
		return errors.New("promo limit reached")
	}

	err = m.UpdateOrderWithPromo(ctx, promo, orderID)
	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx, "UPDATE promo SET promo.limit = promo.limit - 1 WHERE id = ?", promo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) UpdateOrderWithPromo(ctx context.Context, promo order.Promo, orderID uint) error {

	result := m.db.QueryRowContext(ctx, "SELECT total FROM orders WHERE id = ?", orderID)

	var total uint
	err := result.Scan(&total)
	if err != nil {
		return err
	}

	offer := (total * promo.Percentage) / 100

	if promo.Percentage == 100 {
		total = 0
	} else if promo.MaxPrice != 0 && promo.MaxPrice < offer {
		total -= promo.MaxPrice
	} else {
		total -= offer
	}

	_, err = m.db.ExecContext(ctx, "UPDATE orders SET total = ?,promo_id = ? WHERE id = ?", total, promo.ID, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) RemoveOrderPromo(ctx context.Context, orderID uint) error {

	_, err := m.db.ExecContext(ctx, `UPDATE promo SET promo.limit = promo.limit + 1 WHERE id = (SELECT promo_id FROM orders WHERE orders.id = ?)`, orderID)
	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx, "UPDATE orders SET promo_id = NULL WHERE id = ?", orderID)
	if err != nil {
		return err
	}

	err = m.SetOrderTotal(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) DeleteOrder(ctx context.Context, orderID uint) error {

	_, err := m.db.ExecContext(ctx, "DELETE FROM orders WHERE orderID = ?")
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) GetAllOrders(ctx context.Context) ([]order.Order, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM orders")
	if err != nil {
		return []order.Order{}, err
	}

	var rd sql.NullTime
	var stn sql.NullString
	var pid sql.NullInt64
	var phid sql.NullInt64
	var aid sql.NullInt64

	orders := []order.Order{}
	for result.Next() {
		var o order.Order
		err := result.Scan(&o.ID, &o.CreationDate, &rd, &o.Status, &o.Total, &stn, &o.UserID, &pid, &phid, &aid)
		if err != nil {
			return []order.Order{}, err
		}

		o.ReceiptionDate = rd.Time.Format("2006-01-02 15:04:05")
		o.STN = stn.String
		o.Promo.ID = uint(pid.Int64)
		o.PhoneID = uint(phid.Int64)
		o.AddressID = uint(aid.Int64)

		orders = append(orders, o)
	}

	return orders, nil
}

func (m MySQLRepo) GetUserOrders(ctx context.Context, userID string) ([]order.Order, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM orders WHERE user_id = ?", userID)
	if err != nil {
		return []order.Order{}, err
	}

	var rd sql.NullTime
	var stn sql.NullString
	var pid sql.NullInt64
	var phid sql.NullInt64
	var aid sql.NullInt64

	orders := []order.Order{}
	for result.Next() {
		var o order.Order
		err := result.Scan(&o.ID, &o.CreationDate, &rd, &o.Status, &o.Total, &stn, &o.UserID, &pid, &phid, &aid)
		if err != nil {
			return []order.Order{}, err
		}

		o.ReceiptionDate = rd.Time.Format("2006-01-02 15:04:05")
		o.STN = stn.String
		o.Promo.ID = uint(pid.Int64)
		o.PhoneID = uint(phid.Int64)
		o.AddressID = uint(aid.Int64)

		orders = append(orders, o)
	}

	return orders, nil
}

func (m MySQLRepo) GetDateOrders(ctx context.Context, date string) ([]order.Order, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM orders WHERE DATE(creation_date) = ?", date)
	if err != nil {
		return []order.Order{}, err
	}

	var rd sql.NullTime
	var stn sql.NullString
	var pid sql.NullInt64
	var phid sql.NullInt64
	var aid sql.NullInt64

	orders := []order.Order{}
	for result.Next() {
		var o order.Order
		err := result.Scan(&o.ID, &o.CreationDate, &rd, &o.Status, &o.Total, &stn, &o.UserID, &pid, &phid, &aid)
		if err != nil {
			return []order.Order{}, err
		}

		o.ReceiptionDate = rd.Time.Format("2006-01-02 15:04:05")
		o.STN = stn.String
		o.Promo.ID = uint(pid.Int64)
		o.PhoneID = uint(phid.Int64)
		o.AddressID = uint(aid.Int64)

		orders = append(orders, o)
	}

	return orders, nil
}

func (m MySQLRepo) GetDateOrdersByStatus(ctx context.Context, date string, status uint) ([]order.Order, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM orders WHERE DATE(creation_date) = ? AND status = ?", date, status)
	if err != nil {
		return []order.Order{}, err
	}

	var rd sql.NullTime
	var stn sql.NullString
	var pid sql.NullInt64
	var phid sql.NullInt64
	var aid sql.NullInt64

	orders := []order.Order{}
	for result.Next() {
		var o order.Order
		err := result.Scan(&o.ID, &o.CreationDate, &rd, &o.Status, &o.Total, &stn, &o.UserID, &pid, &phid, &aid)
		if err != nil {
			return []order.Order{}, err
		}

		o.ReceiptionDate = rd.Time.Format("2006-01-02 15:04:05")
		o.STN = stn.String
		o.Promo.ID = uint(pid.Int64)
		o.PhoneID = uint(phid.Int64)
		o.AddressID = uint(aid.Int64)

		orders = append(orders, o)
	}

	return orders, nil
}

func (m MySQLRepo) GetAllOrdersByStatus(ctx context.Context, status uint) ([]order.Order, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM orders WHERE status = ?", status)
	if err != nil {
		return []order.Order{}, err
	}

	var rd sql.NullTime
	var stn sql.NullString
	var pid sql.NullInt64
	var phid sql.NullInt64
	var aid sql.NullInt64

	orders := []order.Order{}
	for result.Next() {
		var o order.Order
		err := result.Scan(&o.ID, &o.CreationDate, &rd, &o.Status, &o.Total, &stn, &o.UserID, &pid, &phid, &aid)
		if err != nil {
			return []order.Order{}, err
		}

		o.ReceiptionDate = rd.Time.Format("2006-01-02 15:04:05")
		o.STN = stn.String
		o.Promo.ID = uint(pid.Int64)
		o.PhoneID = uint(phid.Int64)
		o.AddressID = uint(aid.Int64)

		orders = append(orders, o)
	}

	return orders, nil
}

func (m MySQLRepo) GetUserOrdersByStatus(ctx context.Context, userID string, status uint) ([]order.Order, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM orders WHERE user_id = ? AND status = ?", userID, status)
	if err != nil {
		return []order.Order{}, err
	}

	var rd sql.NullTime
	var stn sql.NullString
	var pid sql.NullInt64
	var phid sql.NullInt64
	var aid sql.NullInt64

	orders := []order.Order{}
	for result.Next() {
		var o order.Order
		err := result.Scan(&o.ID, &o.CreationDate, &rd, &o.Status, &o.Total, &stn, &o.UserID, &pid, &phid, &aid)
		if err != nil {
			return []order.Order{}, err
		}

		o.ReceiptionDate = rd.Time.Format("2006-01-02 15:04:05")
		o.STN = stn.String
		o.Promo.ID = uint(pid.Int64)
		o.PhoneID = uint(phid.Int64)
		o.AddressID = uint(aid.Int64)

		orders = append(orders, o)
	}

	return orders, nil
}

func (m MySQLRepo) GetAllPromos(ctx context.Context) ([]order.Promo, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM promo")
	if err != nil {
		return []order.Promo{}, err
	}

	promos := []order.Promo{}
	for result.Next() {
		var p order.Promo
		err := result.Scan(&p.ID, &p.Code, &p.Expiration, &p.Limit, &p.Percentage, &p.MaxPrice)
		if err != nil {
			return []order.Promo{}, err
		}

		promos = append(promos, p)
	}

	return promos, nil
}

func (m MySQLRepo) GetUserPromos(ctx context.Context, userID string) ([]order.Promo, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM promo WHERE id IN (SELECT promo_id FROM promo_user WHERE user_id = ?)", userID)
	if err != nil {
		return []order.Promo{}, err
	}

	promos := []order.Promo{}
	for result.Next() {
		var p order.Promo
		err := result.Scan(&p.ID, &p.Code, &p.Expiration, &p.Limit, &p.Percentage, &p.MaxPrice)
		if err != nil {
			return []order.Promo{}, err
		}

		promos = append(promos, p)
	}

	return promos, nil
}

func (m MySQLRepo) GetPromoByOrder(ctx context.Context, orderID uint) (order.Promo, error) {

	result := m.db.QueryRowContext(ctx, "SELECT * FROM promo WHERE id = (SELECT promo_id FROM orders WHERE orders.id = ?)", orderID)

	var p order.Promo
	var maxPrice sql.NullInt64

	err := result.Scan(&p.ID, &p.Code, &p.Expiration, &p.Limit, &p.Percentage, &maxPrice)
	if err != nil {
		return order.Promo{}, err
	}
	p.MaxPrice = uint(maxPrice.Int64)

	return p, nil
}
