package main

import "database/sql"

type product struct {
	ID       int     `json:"id"`
	Name     string  `json: "name"`
	Quantity int     `json: "quantity"`
	Price    float64 `json: "price"`
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
