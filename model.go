package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	// Query to fetch the data from DB
	query := "SELECT id, name, quantity, price from products"
	rows, err := db.Query(query)

	// Check for errors in Querying
	if err != nil {
		return nil, err
	}

	// Loop on the rows and copy the columns data in the
	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, quantity, price FROM products where id=%v", p.ID)
	rows := db.QueryRow(query) // use QueryRow not Query when expecting atmost 1 row.
	err := rows.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}

	return nil
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("insert into products(name, quantity, price) values('%v', %v, %v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)

	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	result, _ := db.Exec(query)

	// Check the existance of the product in our DB.
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row exists")
	}

	return err
}
