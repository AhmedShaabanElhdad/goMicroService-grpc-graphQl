package order

import (
	"context"
	"database/sql"

	pq "github.com/lib/pq"
)

type Repository interface {
	Close()
	PostOrder(ctx context.Context, requst *Order) (*Order, error)
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetAccountOrders(ctx context.Context, accountId string) ([]Order, error)
}

type postgressRepository struct {
	db *sql.DB
}

func NewPostgressRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgressRepository{
		db: db,
	}, nil

}

func (repo *postgressRepository) Close() {
	repo.db.Close()
}

func (repo *postgressRepository) PostOrder(ctx context.Context, order *Order) (*Order, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	_, err = tx.ExecContext(ctx, "INSERT INTO orders(id,created_at, account_id, totoal_price) values($1,$2,$3,$4)",
		order.ID,
		order.AccountID,
		order.AccountID,
		order.Price,
	)

	if err != nil {
		return nil, err
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_product", "order_id", "product_id, quantity"))

	for _, product := range order.OrderProducts {
		if _, err := stmt.ExecContext(ctx, order.ID, product.ID, product.Quantity); err != nil {
			return nil, err
		}
	}

	defer stmt.Close()

	return order, nil

}

func (repo *postgressRepository) GetOrderById(ctx context.Context, id string) (*Order, error) {

	return nil, nil
}
func (repo *postgressRepository) GetAccountOrders(ctx context.Context, accountId string) ([]Order, error) {

	rows, err := repo.db.QueryContext(
		ctx,
		`
			SELECT 
				o.id,
				o.account_id,
				o.created_at,
				o.totoal_price::money::numeric::float8,
				op.product_id,
				op.quantity

			FROM  orders o JOIN order_products op 
			ON (o.id == op.§)
			WHERE o.id = $1,
			ORDER BY o.id
		`,
		accountId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orderMap := make(map[string]*Order)

	for rows.Next() {
		var (
			id, accountID string
			createdAt     sql.NullTime
			totlaPrice    sql.NullFloat64
			productID     string
			quantity      int
		)

		if err := rows.Scan(
			&id,
			&accountID,
			&createdAt,
			&totlaPrice,
			&productID,
			&quantity,
		); err != nil {
			return nil, err
		}
		order, exist := orderMap[id]
		if !exist {
			order = &Order{
				ID:            id,
				AccountID:     accountID,
				CreatedAt:     createdAt.Time,
				Price:         totlaPrice.Float64,
				OrderProducts: []OrderProduct{},
			}
			orderMap[id] = order
		}
		order.OrderProducts = append(order.OrderProducts, OrderProduct{
			ID:       productID,
			Quantity: quantity,
		})

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	// convert the map to array
	var orders []Order
	for _, order := range orderMap {
		orders = append(orders, *order)
	}
	return orders, nil
}
