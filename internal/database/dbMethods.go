package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
	"wbLab0/internal/configuration"
	"wbLab0/internal/models"
)

type Client interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Close(context.Context) error
}

func NewClient(ctx context.Context, maxAttempts int, sc configuration.StConfig) (conn *pgx.Conn) {
	var err error
	connectionUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	err = attemptDatabaseConnection(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err = pgx.Connect(ctx, connectionUrl)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		fmt.Printf("Error connecting to database")
	}

	return conn
}

func attemptDatabaseConnection(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}

func AddMessageToDatabase(db Client, order models.Order) error {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	defer db.Close(ctx)

	insertItemQuery := "INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shred) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
	_, err = tx.Exec(ctx, insertItemQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShred)
	if err != nil {
		tx.Rollback(ctx)
		return errors.New(fmt.Sprintf("Orders insertion failed (%v)\n", err))
	}

	insertItemQuery = "INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	_, err = tx.Exec(ctx, insertItemQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback(ctx)
		return errors.New(fmt.Sprintf("Delivery insertion failed (%v)\n", err))
	}

	insertItemQuery = "INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
	_, err = tx.Exec(ctx, insertItemQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		tx.Rollback(ctx)
		return errors.New(fmt.Sprintf("Payment insertion failed (%v)\n", err))
	}

	insertItemQuery = "INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);"
	for i := range order.Items {
		_, err = tx.Exec(ctx, insertItemQuery,
			order.OrderUID, order.Items[i].ChrtID, order.Items[i].TrackNumber, order.Items[i].Price, order.Items[i].Rid,
			order.Items[i].Name, order.Items[i].Sale, order.Items[i].Size, order.Items[i].TotalPrice, order.Items[i].NmID,
			order.Items[i].Brand, order.Items[i].Status)
		if err != nil {
			tx.Rollback(ctx)
			return errors.New(fmt.Sprintf("Item %v insertion failed (%v)\n", i, err))
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("Commit failed: (%v)\n", err))
	} else {
		fmt.Printf("Insertion succeded!")
		return nil
	}
}

func SyncCacheAndDatabase(db Client) error {

	ctx := context.Background()

	tx, err := db.Begin(ctx)
	defer db.Close(ctx)

	var rowsInTable int

	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM orders;").Scan(&rowsInTable)
	if err != nil {
		return errors.New(fmt.Sprintf("QueryRow failed (%v)\n", err))
	}

	if len(models.Cache) != rowsInTable {

		rows, err := tx.Query(ctx, "select * from orders;")
		if err != nil {
			return errors.New(fmt.Sprintf("QueryRow (orders) failed (%v)\n", err))
		}
		for rows.Next() {
			var order models.Order
			err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry,
				&order.Locale, &order.InternalSignature, &order.CustomerId,
				&order.DeliveryService, &order.Shardkey, &order.SmId,
				&order.DateCreated, &order.OofShred)
			if err != nil {
				return errors.New(fmt.Sprintf("Error in scanning order row (%v)\n", err))
			}
			if _, found := models.Cache[order.OrderUID]; !found {
				models.Cache[order.OrderUID] = order
			}
		}

		rows, err = tx.Query(ctx, "SELECT * FROM delivery;")
		if err != nil {
			return errors.New(fmt.Sprintf("QueryRow (delivery) failed (%v)\n", err))
		}
		for rows.Next() {
			var delivery models.Delivery
			var uid string
			err := rows.Scan(&uid, &delivery.Name, &delivery.Phone,
				&delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region,
				&delivery.Email)
			if err != nil {
				return errors.New(fmt.Sprintf("Error in scanning delivery row (%v)\n", err))
			}
			if value, found := models.Cache[uid]; found {
				value.Delivery = delivery
				models.Cache[value.OrderUID] = value
			}
		}

		rows, err = tx.Query(ctx, "SELECT * FROM payment;")
		if err != nil {
			return errors.New(fmt.Sprintf("QueryRow (payment) failed (%v)\n", err))
		}
		for rows.Next() {
			var payment models.Payment
			var uid string
			err := rows.Scan(&uid, &payment.Transaction, &payment.RequestID,
				&payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDT,
				&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
			if err != nil {
				return errors.New(fmt.Sprintf("Error in scanning payment row (%v)\n", err))
			}
			if value, found := models.Cache[uid]; found {
				value.Payment = payment
				models.Cache[value.OrderUID] = value
			}
		}

		rows, err = tx.Query(ctx, "SELECT * FROM items;")
		if err != nil {
			return errors.New(fmt.Sprintf("QueryRow (items) failed (%v)\n", err))
		}
		for rows.Next() {
			var item models.Items
			var uid string
			err := rows.Scan(&uid, &item.ChrtID, &item.TrackNumber, &item.Price,
				&item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID,
				&item.Brand, &item.Status)
			if err != nil {
				return errors.New(fmt.Sprintf("Error in scanning item row (%v)\n", err))
			}
			if value, found := models.Cache[uid]; found {
				value.Items = append(value.Items, item)
				models.Cache[value.OrderUID] = value
			}
		}
	}

	log.Printf("Cache and database synchronized.\n")
	return nil
}
