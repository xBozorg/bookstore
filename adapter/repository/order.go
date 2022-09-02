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

func (storage Storage) DoesOrderOpen(ctx context.Context, orderID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT 1 FROM orders WHERE id = ? AND status = ?",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		orderID,
		order.StatusCreated,
	)

	var open bool
	if err = result.Scan(&open); err != nil {
		return false, err
	}

	return open, err
}

func (storage Storage) DoesOrderExist(ctx context.Context, orderID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT 1 FROM orders WHERE id = ?",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, orderID)

	var exists bool
	if err = result.Scan(&exists); err != nil {
		return false, err
	}

	return exists, err
}

func (storage Storage) DoesPromoExist(ctx context.Context, promoID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT 1 FROM promo WHERE id = ?",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, promoID)

	var exists bool
	if err = result.Scan(&exists); err != nil {
		return false, err
	}

	return exists, err
}

func (storage Storage) DoesPromoCodeExist(ctx context.Context, promoCode, userID string) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT 1 FROM promo 
		WHERE id IN (SELECT promo_id FROM promo_user WHERE user_id = ?) AND promo.code = ?`,
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		userID,
		promoCode,
	)

	var exist bool
	if err = result.Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (storage Storage) DoesItemExist(ctx context.Context, itemID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT 1 FROM item WHERE id = ?",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, itemID)

	var exists bool
	if err = result.Scan(&exists); err != nil {
		return false, err
	}

	return exists, err
}

func (storage Storage) CreateEmptyOrder(ctx context.Context, tx *sql.Tx, userID string) (uint, error) {

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO orders 
		(creation_date , status , total , user_id) VALUES (?,?,?,?)`,
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		order.StatusCreated,
		0,
		userID,
	)
	if err != nil {
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(orderID), nil
}

func (storage Storage) CheckOpenOrder(ctx context.Context, tx *sql.Tx, userID string) (uint, error) {

	stmt, err := tx.PrepareContext(ctx,
		"SELECT id FROM orders WHERE user_id=? AND status=?",
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		userID,
		order.StatusCreated,
	)

	var id uint
	if err = result.Scan(&id); err != nil {

		if strings.Contains(err.Error(), "no rows") {
			id, errC := storage.CreateEmptyOrder(ctx, tx, userID)
			if errC != nil {
				return 0, errC
			}
			return id, nil
		}

		return 0, err
	}

	return id, nil
}

func (storage Storage) AddDigitalItem(ctx context.Context, tx *sql.Tx, bookID, orderID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO item (book_id , type , quantity , order_id) 
		VALUES (?,?,?,?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		// Digital
		bookID,
		order.Digital,
		1,
		orderID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddPhysicalItem(ctx context.Context, tx *sql.Tx, item order.Item, orderID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO item (book_id , type , quantity , order_id) 
			VALUES (?,?,?,?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		// Physical
		item.BookID,
		order.Physical,
		item.Quantity,
		orderID,
	); err != nil {
		return err
	}

	stmt, err = tx.PrepareContext(ctx,
		`UPDATE book SET 
			physical_stock = physical_stock - ? 
			WHERE id = ?`,
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, item.Quantity, item.BookID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddBundleItem(ctx context.Context, tx *sql.Tx, item order.Item, orderID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO item (book_id , type , quantity , order_id) 
		VALUES (?,?,?,?) , (?,?,?,?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		// Digital
		item.BookID,
		order.Digital,
		1,
		orderID,
		// Physical
		item.BookID,
		order.Physical,
		item.Quantity,
		orderID,
	); err != nil {
		return err
	}
	stmt, err = tx.PrepareContext(ctx,
		`UPDATE book SET 
		physical_stock = physical_stock - ? 
		WHERE id = ?`,
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, item.Quantity, item.BookID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddItem(ctx context.Context, item order.Item, userID string) error {

	tx, err := storage.MySQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = storage.CheckQuantity(ctx, tx, item.Quantity, item.BookID)
	if err != nil {
		return err
	}

	availability, err := storage.CheckAvailability(ctx, tx, item.BookID)
	if err != nil {
		return err
	}

	orderID, err := storage.CheckOpenOrder(ctx, tx, userID)
	if err != nil {
		return err
	}

	switch {
	case item.Type == order.Bundle && availability == book.BundleAvailable:
		if err = storage.AddBundleItem(ctx, tx, item, orderID); err != nil {
			return err
		}

	case item.Type == order.Physical && (availability == book.PhysicalAvailable || availability == book.BundleAvailable):
		if err = storage.AddPhysicalItem(ctx, tx, item, orderID); err != nil {
			return err
		}

	case item.Type == order.Digital && (availability == book.DigitalAvailable || availability == book.BundleAvailable):
		if err = storage.AddDigitalItem(ctx, tx, item.BookID, orderID); err != nil {
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

	price, err := storage.CalculateItemPrice(ctx, tx, item)
	if err != nil {
		return err
	}

	if err = storage.AddTotal(ctx, tx, price, orderID); err != nil {
		return err
	}

	return tx.Commit()
}

func (storage Storage) GetOrderItems(ctx context.Context, orderID uint) ([]order.Item, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id , book_id , type , quantity FROM item WHERE order_id = ?",
	)
	if err != nil {
		return []order.Item{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, orderID)
	if err != nil {
		return []order.Item{}, err
	}

	items := []order.Item{}
	for result.Next() {
		var i order.Item

		if err = result.Scan(
			&i.ID,
			&i.BookID,
			&i.Type,
			&i.Quantity,
		); err != nil {
			return []order.Item{}, err
		}

		items = append(items, i)
	}

	return items, nil
}

func (storage Storage) CheckQuantity(ctx context.Context, tx *sql.Tx, quantity, bookID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		"SELECT physical_stock FROM book WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, bookID)

	var stock uint
	if err = result.Scan(&stock); err != nil {
		return err
	}

	if quantity > stock {
		return errors.New("requested item quantity is bigger than the available stock")
	}
	return nil
}

func (storage Storage) CheckAvailability(ctx context.Context, tx *sql.Tx, bookID uint) (uint, error) {

	stmt, err := tx.PrepareContext(ctx,
		"SELECT availability FROM book WHERE id = ?",
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, bookID)

	var availability uint
	if err = result.Scan(&availability); err != nil {
		return 0, err
	}

	return availability, nil
}

func (storage Storage) SetOrderPhone(ctx context.Context, orderID, phoneID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`UPDATE orders SET phone_id = ? WHERE id = ?
		AND (SELECT 1 FROM phone WHERE id = ? AND phone.userID = orders.user_id)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		phoneID,
		orderID,
		phoneID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) SetOrderAddress(ctx context.Context, orderID, addressID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`UPDATE orders SET address_id = ? WHERE id = ?
		AND (SELECT 1 FROM address WHERE id = ? AND address.userID = orders.user_id)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		addressID,
		orderID,
		addressID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) IncreaseQuantity(ctx context.Context, itemID, orderID uint) error {

	tx, err := storage.MySQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`SELECT book_id , type , quantity FROM item WHERE id = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var item order.Item
	if err := stmt.QueryRowContext(ctx, itemID).Scan(&item.BookID, &item.Type, &item.Quantity); err != nil {
		return err
	}
	item.ID = itemID

	if item.Type == order.Digital {
		return errors.New("cannot increase digital item")
	}

	stmt, err = tx.PrepareContext(ctx,
		`UPDATE item SET 
		item.quantity = item.quantity + 1 
		WHERE item.id = ? 
		AND
		item.quantity < ( SELECT book.physical_stock FROM book WHERE book.id = item.book_id )`,
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, itemID); err != nil {
		return err
	}

	stmt, err = tx.PrepareContext(ctx,
		`UPDATE book SET physical_stock = physical_stock - 1 
		WHERE id = ?`,
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, item.BookID); err != nil {
		return err
	}

	price, err := storage.CalculateItemPrice(ctx, tx, item)
	if err != nil {
		return err
	}

	if err = storage.AddTotal(ctx, tx, price, orderID); err != nil {
		return err
	}

	return tx.Commit()
}

func (storage Storage) DecreaseQuantity(ctx context.Context, itemID, orderID uint) error {

	tx, err := storage.MySQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`SELECT book_id , type , quantity FROM item WHERE id = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var item order.Item
	if err := stmt.QueryRowContext(ctx, itemID).Scan(&item.BookID, &item.Type, &item.Quantity); err != nil {
		return err
	}
	item.ID = itemID

	stmt, err = tx.PrepareContext(ctx,
		`UPDATE item SET 
		quantity = quantity - 1 
		WHERE id = ? 
		AND 
		quantity > 0`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, itemID); err != nil {
		return err
	}

	stmt, err = tx.PrepareContext(ctx,
		`UPDATE book SET 
		physical_stock = physical_stock + 1 
		WHERE id = ?`,
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, item.BookID); err != nil {
		return err
	}

	price, err := storage.CalculateItemPrice(ctx, tx, item)
	if err != nil {
		return err
	}

	if err = storage.SubTotal(ctx, tx, price, orderID); err != nil {
		return err
	}

	return tx.Commit()
}

func (storage Storage) RemoveItem(ctx context.Context, itemID, orderID uint) error {

	tx, err := storage.MySQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`SELECT book_id , type , quantity FROM item WHERE id = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var item order.Item
	if err := stmt.QueryRowContext(ctx, itemID).Scan(&item.BookID, &item.Type, &item.Quantity); err != nil {
		return err
	}
	item.ID = itemID

	if item.Type == order.Physical || item.Type == order.Bundle {
		stmt, err = tx.PrepareContext(ctx,
			`UPDATE book SET 
			physical_stock = physical_stock + (SELECT quantity FROM item WHERE id = ?)
			WHERE 
			book.id = (SELECT book_id FROM item WHERE id = ?)`,
		)
		if err != nil {
			return err
		}
		if _, err = stmt.ExecContext(ctx, itemID, itemID); err != nil {
			return err
		}
	}

	stmt, err = tx.PrepareContext(ctx,
		"DELETE FROM item WHERE id = ?",
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, itemID); err != nil {
		return err
	}

	price, err := storage.CalculateItemPrice(ctx, tx, item)
	if err != nil {
		return err
	}

	if err = storage.SubTotal(ctx, tx, price, orderID); err != nil {
		return err
	}

	return tx.Commit()
}

func (storage Storage) CreatePromoCode(ctx context.Context, promo order.Promo, userID string) error {

	if promo.Percentage == 0 {
		return errors.New("percentage cannot be 0")
	}
	if promo.Limit == 0 {
		return errors.New("limit cannot be 0")
	}

	tx, err := storage.MySQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO promo 
		(code , expiration , promo.limit , percentage , max_price)
		VALUES (?,?,?,?,?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		promo.Code,
		promo.Expiration,
		promo.Limit,
		promo.Percentage,
		promo.MaxPrice,
	)
	if err != nil {
		return err
	}

	promoID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stmt, err = tx.PrepareContext(ctx,
		"INSERT INTO promo_user (promo_id , user_id) VALUES (?,?)",
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, promoID, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (storage Storage) DeletePromoCode(ctx context.Context, promoID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM promo WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		promoID,
	); err != nil {
		return err
	}
	/*
		_, err = storage.MySQL.ExecContext(
			ctx,
			"DELETE FROM promo_user WHERE promo_id = ? AND user_id = ?",
			promoID,
			userID,
		)

		if err != nil {
			return err
		}
	*/

	return nil
}

func (storage Storage) SetOrderStatus(ctx context.Context, status, orderID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT 1 FROM item WHERE type != 0 AND order_id = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, orderID)

	var isShipmentOrder bool
	if err = result.Scan(&isShipmentOrder); err != nil {
		return err
	}

	if status != order.StatusCreated && isShipmentOrder {

		stmt, err = storage.MySQL.PrepareContext(ctx,
			`UPDATE orders AS o , (SELECT phone_id , address_id FROM orders WHERE id = ?) AS PA
			SET status = ? 
			WHERE PA.phone_id IS NOT NULL AND PA.address_id IS NOT NULL
			AND id = ?
			`,
		)
		if err != nil {
			return err
		}

		if _, err = stmt.ExecContext(ctx,
			orderID,
			status,
			orderID,
		); err != nil {
			return err
		}

	} else {

		stmt, err = storage.MySQL.PrepareContext(ctx,
			"UPDATE orders SET status = ? WHERE id = ?",
		)
		if err != nil {
			return err
		}

		if _, err = stmt.ExecContext(ctx,
			status,
			orderID,
		); err != nil {
			return err
		}

	}

	return nil
}

func (storage Storage) GetOrderStatus(ctx context.Context, orderID uint) (uint, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT status FROM orders WHERE id = ?",
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, orderID)

	var status uint
	if err = result.Scan(&status); err != nil {
		return 0, err
	}

	return status, nil
}

func (storage Storage) SetOrderSTN(ctx context.Context, stn string, orderID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"UPDATE orders SET stn = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		stn,
		orderID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) SetOrderTotal(ctx context.Context, tx *sql.Tx, total, orderID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE orders SET total = ? WHERE id = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, total, orderID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) CalculateItemPrice(ctx context.Context, tx *sql.Tx, item order.Item) (uint, error) {

	stmt, err := tx.PrepareContext(ctx,
		`SELECT digital_price , digital_discount , physical_price , physical_discount 
		FROM book WHERE id = ?`,
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, item.BookID)

	var i itemPrice
	if err = result.Scan(
		&i.DigitalPrice,
		&i.DigitalDiscount,
		&i.PhysicalPrice,
		&i.PhysicalDiscount,
	); err != nil {
		return 0, err
	}

	var price uint

	switch item.Type {
	case order.Digital:
		price = (i.DigitalPrice * (100 - i.DigitalDiscount)) / 100
	case order.Physical:
		price = (i.PhysicalPrice * (100 - i.PhysicalDiscount)) / 100 * item.Quantity
	case order.Bundle:
		price = ((i.DigitalPrice*(100-i.DigitalDiscount))/100 +
			(i.PhysicalPrice*(100-i.PhysicalDiscount))/100*item.Quantity) * 80 / 100
	default:
		return 0, errors.New("invalid order type")
	}

	return price, nil
}

func (storage Storage) AddTotal(ctx context.Context, tx *sql.Tx, value, orderID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE orders SET total = total + ? WHERE id = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, value, orderID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) SubTotal(ctx context.Context, tx *sql.Tx, value, orderID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE orders SET total = total - ? WHERE id = ? AND total >= ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, value, orderID, value); err != nil {
		return err
	}

	return nil
}

func (storage Storage) CalculateOrderTotal(ctx context.Context, tx *sql.Tx, orderID uint) error {

	stmt, err := tx.PrepareContext(ctx,
		"SELECT book_id , type , quantity FROM item WHERE order_id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, orderID)
	if err != nil {
		return err
	}

	items := []itemPrice{}
	for result.Next() {
		var i itemPrice

		if err = result.Scan(
			&i.BookID,
			&i.Type,
			&i.Quantity,
		); err != nil {
			return err
		}

		stmt, err := tx.PrepareContext(ctx,
			`SELECT digital_price , digital_discount , physical_price , physical_discount 
			FROM book WHERE id = ?`,
		)
		if err != nil {
			return err
		}
		defer stmt.Close()

		priceResult := stmt.QueryRowContext(ctx, i.BookID)
		if err = priceResult.Scan(
			&i.DigitalPrice,
			&i.DigitalDiscount,
			&i.PhysicalPrice,
			&i.PhysicalDiscount,
		); err != nil {
			return err
		}

		items = append(items, i)
	}

	var total uint = 0
	for _, item := range items {
		switch item.Type {
		case order.Digital:
			total += (item.DigitalPrice * (100 - item.DigitalDiscount)) / 100
		case order.Physical:
			total += (item.PhysicalPrice * (100 - item.PhysicalDiscount)) / 100 * item.Quantity
		default:
			return errors.New("invalid order type")
		}
	}

	if err = storage.SetOrderTotal(ctx, tx, total, orderID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) SetOrderPromo(ctx context.Context, orderID uint, promoCode, userID string) error {

	tx, err := storage.MySQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`SELECT * FROM promo 
		WHERE code = ? 
		AND 
		(SELECT 1 FROM promo_user WHERE promo_id = promo.id AND user_id = ?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		promoCode,
		userID,
	)

	var promo order.Promo

	if err = result.Scan(
		&promo.ID,
		&promo.Code,
		&promo.Expiration,
		&promo.Limit,
		&promo.Percentage,
		&promo.MaxPrice,
	); err != nil {
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

	tx, err = storage.UpdateOrderWithPromo(ctx, tx, promo, orderID)
	if err != nil {
		return err
	}

	stmt, err = tx.PrepareContext(ctx,
		"UPDATE promo SET promo.limit = promo.limit - 1 WHERE id = ?",
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, promo.ID); err != nil {
		return err
	}

	return tx.Commit()
}

func (storage Storage) UpdateOrderWithPromo(ctx context.Context, tx *sql.Tx, promo order.Promo, orderID uint) (*sql.Tx, error) {

	stmt, err := tx.PrepareContext(ctx,
		"SELECT total FROM orders WHERE id = ?",
	)
	if err != nil {
		return tx, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, orderID)

	var total uint
	if err = result.Scan(&total); err != nil {
		return tx, err
	}

	offer := (total * promo.Percentage) / 100

	if promo.Percentage == 100 {
		total = 0
	} else if promo.MaxPrice != 0 && promo.MaxPrice < offer {
		total -= promo.MaxPrice
	} else {
		total -= offer
	}

	stmt, err = tx.PrepareContext(ctx,
		"UPDATE orders SET total = ?,promo_id = ? WHERE id = ?",
	)
	if err != nil {
		return tx, err
	}

	if _, err = stmt.ExecContext(ctx,
		total,
		promo.ID,
		orderID,
	); err != nil {
		return tx, err
	}

	return tx, nil
}

func (storage Storage) RemoveOrderPromo(ctx context.Context, orderID uint) error {

	tx, err := storage.MySQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE promo SET 
		promo.limit = promo.limit + 1 
		WHERE id = (SELECT promo_id FROM orders WHERE orders.id = ?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, orderID); err != nil {
		return err
	}

	stmt, err = tx.PrepareContext(ctx,
		"UPDATE orders SET promo_id = NULL WHERE id = ?",
	)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, orderID); err != nil {
		return err
	}

	if err = storage.CalculateOrderTotal(ctx, tx, orderID); err != nil {
		return err
	}

	return tx.Commit()
}

func (storage Storage) DeleteOrder(ctx context.Context, orderID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM orders WHERE orderID = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, orderID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) GetAllOrders(ctx context.Context) ([]order.Order, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM orders",
	)
	if err != nil {
		return []order.Order{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
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

		if err = result.Scan(
			&o.ID,
			&o.CreationDate,
			&rd,
			&o.Status,
			&o.Total,
			&stn,
			&o.UserID,
			&pid,
			&phid,
			&aid,
		); err != nil {
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

func (storage Storage) GetUserOrders(ctx context.Context, userID string) ([]order.Order, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM orders WHERE user_id = ?",
	)
	if err != nil {
		return []order.Order{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, userID)
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

		if err = result.Scan(
			&o.ID,
			&o.CreationDate,
			&rd,
			&o.Status,
			&o.Total,
			&stn,
			&o.UserID,
			&pid,
			&phid,
			&aid,
		); err != nil {
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

func (storage Storage) GetDateOrders(ctx context.Context, date string) ([]order.Order, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM orders WHERE DATE(creation_date) = ?",
	)
	if err != nil {
		return []order.Order{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, date)
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

		if err = result.Scan(
			&o.ID,
			&o.CreationDate,
			&rd,
			&o.Status,
			&o.Total,
			&stn,
			&o.UserID,
			&pid,
			&phid,
			&aid,
		); err != nil {
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

func (storage Storage) GetDateOrdersByStatus(ctx context.Context, date string, status uint) ([]order.Order, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM orders WHERE DATE(creation_date) = ? AND status = ?",
	)
	if err != nil {
		return []order.Order{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx,
		date,
		status,
	)
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

		if err = result.Scan(
			&o.ID,
			&o.CreationDate,
			&rd,
			&o.Status,
			&o.Total,
			&stn,
			&o.UserID,
			&pid,
			&phid,
			&aid,
		); err != nil {
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

func (storage Storage) GetAllOrdersByStatus(ctx context.Context, status uint) ([]order.Order, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM orders WHERE status = ?",
	)
	if err != nil {
		return []order.Order{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, status)
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

		if err = result.Scan(
			&o.ID,
			&o.CreationDate,
			&rd,
			&o.Status,
			&o.Total,
			&stn,
			&o.UserID,
			&pid,
			&phid,
			&aid,
		); err != nil {
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

func (storage Storage) GetUserOrdersByStatus(ctx context.Context, userID string, status uint) ([]order.Order, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM orders WHERE user_id = ? AND status = ?",
	)
	if err != nil {
		return []order.Order{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx,
		userID,
		status,
	)
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

		if err = result.Scan(
			&o.ID,
			&o.CreationDate,
			&rd,
			&o.Status,
			&o.Total,
			&stn,
			&o.UserID,
			&pid,
			&phid,
			&aid,
		); err != nil {
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

func (storage Storage) GetAllPromos(ctx context.Context) ([]order.Promo, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM promo",
	)
	if err != nil {
		return []order.Promo{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []order.Promo{}, err
	}

	promos := []order.Promo{}
	for result.Next() {
		var p order.Promo

		if err = result.Scan(
			&p.ID,
			&p.Code,
			&p.Expiration,
			&p.Limit,
			&p.Percentage,
			&p.MaxPrice,
		); err != nil {
			return []order.Promo{}, err
		}

		promos = append(promos, p)
	}

	return promos, nil
}

func (storage Storage) GetUserPromos(ctx context.Context, userID string) ([]order.Promo, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM promo WHERE id IN (SELECT promo_id FROM promo_user WHERE user_id = ?)",
	)
	if err != nil {
		return []order.Promo{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, userID)

	if err != nil {
		return []order.Promo{}, err
	}

	promos := []order.Promo{}
	for result.Next() {
		var p order.Promo

		if err = result.Scan(
			&p.ID,
			&p.Code,
			&p.Expiration,
			&p.Limit,
			&p.Percentage,
			&p.MaxPrice,
		); err != nil {
			return []order.Promo{}, err
		}

		promos = append(promos, p)
	}

	return promos, nil
}

func (storage Storage) GetPromoByOrder(ctx context.Context, orderID uint) (order.Promo, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM promo WHERE id = (SELECT promo_id FROM orders WHERE orders.id = ?)",
	)
	if err != nil {
		return order.Promo{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, orderID)

	var p order.Promo
	var maxPrice sql.NullInt64

	if err = result.Scan(
		&p.ID,
		&p.Code,
		&p.Expiration,
		&p.Limit,
		&p.Percentage,
		&maxPrice,
	); err != nil {
		return order.Promo{}, err
	}
	p.MaxPrice = uint(maxPrice.Int64)

	return p, nil
}

func (storage Storage) GetOrderPaymentInfo(ctx context.Context, orderID uint) (order.OrderPaymentInfo, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT total , user_id , phone_id FROM orders WHERE id = ? AND status = ?`,
	)
	if err != nil {
		return order.OrderPaymentInfo{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		orderID,
		order.StatusCreated,
	)

	info := order.OrderPaymentInfo{}

	var uid string
	var pid sql.NullInt64

	if err = result.Scan(
		&info.Total,
		&uid,
		&pid,
	); err != nil {
		return order.OrderPaymentInfo{}, err
	}

	stmt, err = storage.MySQL.PrepareContext(ctx,
		`SELECT email FROM user WHERE id = ?`,
	)
	if err != nil {
		return order.OrderPaymentInfo{}, err
	}
	defer stmt.Close()

	result = stmt.QueryRowContext(ctx, uid)
	if err = result.Scan(&info.Email); err != nil {
		return order.OrderPaymentInfo{}, err
	}

	stmt, err = storage.MySQL.PrepareContext(ctx,
		`SELECT phonenumber FROM phone WHERE id = ?`,
	)
	if err != nil {
		return order.OrderPaymentInfo{}, err
	}
	defer stmt.Close()

	result = stmt.QueryRowContext(ctx, uint(pid.Int64))
	if err = result.Scan(&info.Phone); err != nil {
		return order.OrderPaymentInfo{}, err
	}

	return info, nil
}

func (storage Storage) GetOrderTotal(ctx context.Context, orderID uint) (uint, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT total FROM orders WHERE id = ?`,
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, orderID)

	var total uint
	if err = result.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (storage Storage) ZarinpalCreateOpenOrder(ctx context.Context, orderID uint, authority string) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`INSERT INTO zarinpal (order_id , authority , code) VALUES (?,?,0)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(
		ctx,
		orderID,
		authority,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) ZarinpalDoesAuthorityExist(ctx context.Context, authority string) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT 1 FROM zarinpal WHERE authority = ?`,
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, authority)

	var exist bool
	if err = result.Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (storage Storage) ZarinpalGetOrderByAuthority(ctx context.Context, authority string) (order.ZarinpalOrder, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT * FROM zarinpal WHERE authority = ?`,
	)
	if err != nil {
		return order.ZarinpalOrder{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, authority)

	var refID sql.NullInt64
	var zOrder order.ZarinpalOrder

	if err = result.Scan(
		&zOrder.ID,
		&zOrder.OrderID,
		&zOrder.Authority,
		&refID,
		&zOrder.Code,
	); err != nil {
		return order.ZarinpalOrder{}, err
	}
	if refID.Valid {
		zOrder.RefID = int(refID.Int64)
	}

	return zOrder, nil
}

func (storage Storage) ZarinpalSetOrderPayment(ctx context.Context, zarinpalOrderID uint, authority string, refID, code int) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`UPDATE zarinpal SET ref_id = ? , code = ? WHERE id = ? AND authority = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		refID,
		code,
		zarinpalOrderID,
		authority,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) SetOrderReceiptDate(ctx context.Context, orderID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`UPDATE orders SET receipt_date = ? WHERE id = ?`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		orderID,
	); err != nil {
		return err
	}

	return nil
}
