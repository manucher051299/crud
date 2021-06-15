package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

//ErrNotFound возврашается, когда покупатель не найден
var ErrNotFound = errors.New("Item not found")

//ErrInternal возвращается, когда произашла внутренняя ошибка
var ErrInternal = errors.New("internal error")

//Service представляет собой сервис по управлению баннерами
type Service struct {
	pool *pgxpool.Pool
}

//NewService создаёт сервис
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

//Переменная для Id
//var bannerId int64

//////////////////////
func (s *Service) ByID(ctx context.Context, ID int64) (*Customer, error) {

	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id=$1
	`, ID).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

/////////////////////////////
func (s *Service) All(ctx context.Context) ([]*Customer, error) {

	items := make([]*Customer, 0)

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created FROM customers
	`)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
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

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE active
	`)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
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

func (s *Service) Save(ctx context.Context, customer *Customer) (*Customer, error) {

	item := &Customer{
		ID:    customer.ID,
		Name:  customer.Name,
		Phone: customer.Phone,
	}

	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
	INSERT INTO customers(name,phone) VALUES ($1,$2) ON CONFLICT(phone) DO UPDATE SET name=excluded.name RETURNING id,name,phone,active,created;
	`, item.Name, item.Phone).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		return item, nil
	}

	err := s.pool.QueryRow(ctx, `
UPDATE customers SET name=$1 phone=$2 RETURNING id,name,phone,active,created;
`, item.Name, item.Phone).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

/////////////////////////
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
	DELETE FROM customers WHERE id=$1 RETURNING id,name,phone,active,created;
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

func (s *Service) BlockById(ctx context.Context, id int64) (*Customer, error) {

	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
UPDATE customers SET active=false WHERE id = $1 RETURNING id,name,phone,active,created	
`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}
	return item, nil
}

func (s *Service) UnblockById(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
	UPDATE customers SET active=true WHERE id = $1 RETURNING id,name,phone,active,created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}
	return item, nil
}
