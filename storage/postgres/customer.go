package postgres

import (
	"database/sql"
	"fmt"
	"rent-car/models"
	"rent-car/pkg"

	"github.com/google/uuid"
)

type customerRepo struct {
	db *sql.DB
}

func NewCustomer(db *sql.DB) customerRepo {
	return customerRepo{
		db: db,
	}
}

func (c *customerRepo) Create(customer models.Customer) (string, error) {
	id := uuid.New()
	query := `INSERT INTO customers(
		id,
		first_name,
		last_name,
		gmail,
		phone,
		is_Blocked
	) Values($1,$2,$3,$4,$5,$6)`

	_, err := c.db.Exec(query, id.String(),
		customer.FirstName,
		customer.LastName,
		customer.Gmail,
		customer.Phone,
		customer.Is_Blocked)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (c *customerRepo) Update(customer models.Customer) (string, error) {

	query := `UPDATE customers set
	
	firstName=$2,
	lastName=$3,
	gmail=$4,
	phone=$5,
	is_Blocked=$6,
	updated_at=CURRENT_TIMESTAMP,
	orders=$1 WHERE id=$7 AND deleted_at=0
	`
	_, err := c.db.Exec(query,
		customer.Orders,
		customer.FirstName,
		customer.LastName,
		customer.Gmail,
		customer.Phone,
		customer.Is_Blocked,
		customer.Id)

	if err != nil {
		return "", err
	}

	return customer.Id, nil

}

func (c customerRepo) GetAll(req models.GetAllCustomerRequest) (models.GetAllCustomersResponse, error) {
	var (
		resp   = models.GetAllCustomersResponse{}
		filter = ""
	)

	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` and name ILIKE '%%%v%%' `, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := c.db.Query(`select 
	count(id) OVER(),
		id,
		first_name,
		last_name,
		gmail,
		phone,
		is_Blocked,
		created_at::date,
		updated_at 
		FROM customers where deleted_at=0` + filter + ``)

	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			customer  = models.Customer{}
			updatedAt sql.NullString
		)

		if err := rows.Scan(
			&resp.Count,
			&customer.Id,
			&customer.FirstName,
			&customer.LastName,
			&customer.Gmail,
			&customer.Phone,
			&customer.Is_Blocked,
			&customer.CreatedAt, &updatedAt); err != nil {
			return resp, err
		}

		customer.UpdatedAt = pkg.NullStringToString(updatedAt)
		resp.Customers = append(resp.Customers, customer)
	}
	return resp, nil

}

func (c *customerRepo) Delete(id string) error {
	query := `UPDATE customers set 
				deleted_at = date_part('epoch', CURRENT_TIMESTAMP)::int
				where id = $1 and deleted_at=0`

	_, err := c.db.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}
