package data

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"time"
)

// Product defines the structure for an API product
type ProductInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	SKU         string    `json:"sku"`
	CreatedOn   time.Time `json:"-"`
	UpdatedOn   time.Time `json:"-"`
	DeletedOn   time.Time `json:"-"`
}

// Define a ProductModel type which wraps a sql.DB connection pool
type ProductModel struct {
	DB *sql.DB
}

func (p *ProductInfo) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// this will insert a new snippet into database
func (m *ProductModel) Insert(p *ProductInfo) (uuid.UUID, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	//stmt := "SELECT MAX(ID) FROM product"
	//var latestID int
	//err := m.DB.QueryRow(stmt).Scan(&latestID)
	//if err != nil {
	//	return 0, err
	//}
	//
	//newID := latestID + 1
	p.ID = uuid.New()

	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for the placeholder parameters. This
	// method returns a sql.Result type, which contains some basic
	// information about what happened when the statement was executed.
	stmt := "INSERT INTO product (ID, name, description, price, sku) VALUES (@ID, @Name, @Description, @Price, @SKU)"
	_, err := m.DB.Exec(stmt, sql.Named("ID", p.ID), sql.Named("Name", p.Name), sql.Named("Description", p.Description), sql.Named("Price", p.Price), sql.Named("SKU", p.SKU))
	if err != nil {
		return uuid.Nil, err
	}
	//// Use the LastInsertId() method on the result to get the ID of our
	//// newly inserted record in the snippets table.
	//_, err = result.LastInsertId()
	//if err != nil {
	//	return 0, err
	//}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return p.ID, nil
}

// this will return the 10 most recently created products
func (m *ProductModel) Latest() ([]*ProductInfo, error) {
	query := `
		SELECT TOP 10
			ID,
			Name,
			Description,
			Price,
			SKU
-- 			Created_On,
-- 			Updated_On,
-- 			Deleted_On
		FROM
			product
-- 		ORDER BY
-- 			Created_On DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	products := []*ProductInfo{}

	for rows.Next() {
		product := &ProductInfo{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.SKU,
			//&product.CreatedOn,
			//&product.UpdatedOn,
			//&product.DeletedOn,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return products, nil
}
