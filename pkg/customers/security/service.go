package security

import (
	"context"
	"errors"
	"log"

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

type Auth struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Service) Auth(login, password string) (ok bool) {

	cusPass := ""
	ctx := context.Background()

	err := s.pool.QueryRow(ctx, `
	 SELECT password FROM managers WHERE login = $1  
	   `, login).Scan(&cusPass)

	if err != nil {
		log.Print("Not ook", err)
		return false
	}

	if password != cusPass {
		log.Print("Not OOk", err)
		return false
	}

	log.Print("OK")
	return true

}
