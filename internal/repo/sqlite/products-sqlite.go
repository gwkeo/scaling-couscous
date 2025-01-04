package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID    int64
	Name  string
	Price float64
}

type ProductsRepoSqlite struct {
	db *sql.DB
}

func NewProductsRepo(ctx context.Context, path string) (*ProductsRepoSqlite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error while opening sqlite db: %w", err)
	}

	stmt, err := db.PrepareContext(ctx, `
		CREATE TABLE IF NOT EXISTS products (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name VARCHAR(255) UNIQUE NOT NULL,
	    price INT NOT NULL);`)

	if err != nil {
		return nil, fmt.Errorf("error while preparing statement: %w", err)
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while executing statement: %w", err)
	}

	return &ProductsRepoSqlite{db: db}, nil
}

func (r *ProductsRepoSqlite) CreateProduct(ctx context.Context, product *Product) (int64, error) {
	errLocation := "create product"
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO products (name, price) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error while preparing statement: %w", err)
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			return
		}
	}()

	res, err := stmt.ExecContext(ctx, product.Name, product.Price)
	if err != nil {
		return 0, fmt.Errorf("%s: error while executing statement: %w", errLocation, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: error while getting last insert ID: %w", errLocation, err)
	}

	return id, nil
}

func (r *ProductsRepoSqlite) Product(ctx context.Context, id int64) (*Product, error) {
	errLocation := "get product"

	product := &Product{}
	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM products WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("%s: error while preparing statement: %w", errLocation, err)
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			return
		}
	}()

	err = stmt.QueryRowContext(ctx, id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, errors.New("product not found")
		}
		return nil, fmt.Errorf("%s: error while executing statement: %w", errLocation, err)
	}

	return product, nil
}

func (r *ProductsRepoSqlite) AllProducts(ctx context.Context) ([]Product, error) {
	errLocation := "get all products"
	var products []Product
	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("%s: error while preparing statement: %w", errLocation, err)
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			return
		}
	}()

	err = stmt.QueryRowContext(ctx).Scan(&products)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return products, errors.New("products not found")
		}
		return nil, fmt.Errorf("%s: error while executing statement: %w", errLocation, err)
	}

	return products, nil
}

func (r *ProductsRepoSqlite) ProductByName(ctx context.Context, name string) (*Product, error) {
	errLocation := "get product by name"
	product := &Product{}
	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM products WHERE name=?")
	if err != nil {
		return nil, fmt.Errorf("%s: error while preparing statement: %w", errLocation, err)
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			return
		}
	}()

	err = stmt.QueryRowContext(ctx, name).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, errors.New("product not found")
		}
		return nil, fmt.Errorf("%s: error while executing statement: %w", errLocation, err)
	}

	return product, nil
}

func (r *ProductsRepoSqlite) UpdateProduct(ctx context.Context, product *Product) error {
	errLocation := "update product"
	stmt, err := r.db.PrepareContext(ctx, "UPDATE products SET name=?, price=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("%s: error while preparing statement: %w", errLocation, err)
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			return
		}
	}()
	_, err = stmt.ExecContext(ctx, product.Name, product.Price, product.ID)
	if err != nil {
		return fmt.Errorf("%s: error while executing statement: %w", errLocation, err)
	}

	return nil
}

func (r *ProductsRepoSqlite) DeleteProduct(ctx context.Context, id int64) error {
	errLocation := "delete product"
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM products WHERE id=?")
	if err != nil {
		return fmt.Errorf("%s: error while preparing statement: %w", errLocation, err)
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: error while executing statement: %w", errLocation, err)
	}
	return nil
}
