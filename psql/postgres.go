package psql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rubenv/sql-migrate"
	"log"
	"time"
)

// Postgres Хранилище данных
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

func New(connStr string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{}

	// Пользовательские параметры
	for _, opt := range opts {
		opt(pg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), pg.connTimeout)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			// Проверяем, что подключение действительно было установлено
			err = pg.Pool.Ping(ctx)
			if err == nil {
				// Подключение успешно, выходим из цикла
				break
			}

			pg.connAttempts--
		}
		log.Println("Postgres пытается подключиться, попыток осталось:", pg.connAttempts+1)

		time.Sleep(pg.connTimeout)
	}
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close Закрыть
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// Migrate Миграция таблиц
func (p *Postgres) Migrate() error {
	// Прочитать миграции из папки:
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	// Преобразование pgxpool.Pool в *sql.DB
	db := stdlib.OpenDBFromPool(p.Pool)

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal("errors", err)
	}
	log.Println("Применена миграция!", n)

	return nil
}
