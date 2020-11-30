package customers

import (
	"log"
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("Item not found")
var ErrInternal = errors.New("Internal error")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Customer struct {
	ID		int64 		`json:"id"`
	Name	string 		`json:"name"`
	Phone	string		`json:"phone"`
	Active	bool		`json:"active"`
	Created	time.Time	`json:"created"`
}
func (s *Service) UnBlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx,`
		UPDATE customers SET active = true WHERE id = $1 RETURNING id, name, phone, active, created
	`,id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}
func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx,`
		UPDATE customers SET active = false WHERE id = $1 RETURNING id, name, phone, active, created
	`,id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx,`
		DELETE FROM customers WHERE id = $1 RETURNING id, name, phone, active, created;
	`,id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}
func (s *Service) Save(ctx context.Context, customer *Customer) (*Customer, error) {
	item := &Customer{}
	if customer.ID == 0 {
		err := s.db.QueryRowContext(ctx,`
			INSERT INTO customers (name, phone) VALUES ($1, $2) RETURNING id, name, phone, active, created
		`,customer.Name, customer.Phone).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		return item, nil
	}
	err := s.db.QueryRowContext(ctx,`
		UPDATE customers SET name = $1, phone = $2 WHERE id = $3 RETURNING id, name, phone, active, created
	`,customer.Name, customer.Phone, customer.ID).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx,`
		SELECT id, name, phone, active, created FROM customers WHERE id = $1	
	`,id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}
func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, phone, active, created FROM customers	
	`)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return items, nil
}
func (s *Service) AllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE active	
	`)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return items, nil
}