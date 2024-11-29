package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	Ping() error
	AddAccount(ctx context.Context, account Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	FetchAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
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

func (repo *postgressRepository) Ping() error {
	return repo.db.Ping()
}

func (repo *postgressRepository) AddAccount(ctx context.Context, account Account) error {
	_, err := repo.db.ExecContext(ctx, "Insert INTO accounts(id, name) VALUES($1,$2)", account.ID, account.Name)
	return err
}

func (repo *postgressRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT id, Name FROM accounts WHERE id = $1", id)
	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}
	return a, nil

}

func (repo *postgressRepository) FetchAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {

	rows, err := repo.db.QueryContext(ctx, "SELECT id, Name FROM accounts ORDER BY id DESC $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := []Account{}

	for rows.Next() {
		a := &Account{}
		if err := rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil

}
