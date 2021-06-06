package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *Service) ById(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id=$1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
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
		return nil, ErrInternal
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
			return
		}
	}()
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrInternal
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
		return nil, ErrInternal
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
			return
		}
	}()
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrInternal
	}
	return items, nil
}
func (s *Service) Save(ctx context.Context, name string, phone string, id int64) error {
	if id == 0 {
		_, err := s.db.ExecContext(ctx, `
		INSERT INTO customers(name,phone) VALUES ($1,$2) ON CONFLICT(phone) DO UPDATE SET name=excluded.name;
		`, name, phone)
		if err != nil {
			log.Print(err)
			return ErrInternal
		}
		return nil
	}

	_, err := s.db.ExecContext(ctx, `
	UPDATE customers SET name=$1 phone=$2;
	`, name, phone)
	if err != nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}

func (s *Service) RemoveById(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
	DELETE FROM customers WHERE id=$1;
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}

func (s *Service) BlockById(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
	UPDATE customers SET active=0 WHERE id=$1;
	`, id)
	if err != nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}

func (s *Service) UnblockById(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
	UPDATE customers SET active=1 WHERE id=$1;
	`, id)
	if err != nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}
