# tools_postgres
# psql

Пакет `psql` предоставляет удобный интерфейс для работы с PostgreSQL, включая управление пулом соединений, выполнение миграций и настройку параметров подключения.

## Установка

Для использования пакета добавьте его в ваш проект:

```bash
go get github.com/zatrasz75/tools_postgres
```

## Использование

```bash
func main() {
  package main

import (
	"log"
	"time"

	"github.com/zatrasz75/tools_postgres/psql"
)
  
	// Строка подключения к PostgreSQL
	connStr := "postgres://user:password@localhost:5432/dbname"
	
		// Настройки подключения к PostgreSQL
	PoolMax := 10                   // Максимальный размер пула соединений
	ConnAttempts := 5               // Количество попыток подключения
	ConnTimeout := 10 * time.Second // Время ожидания подключения

	// Создание нового подключения с настройками
	pg, err := psql.New(connStr, psql.OptionSet(PoolMax, ConnAttempts, ConnTimeout))
	if err != nil {
		l.Fatal("ошибка запуска - postgres.New:", err)
	}
	defer pg.Close()
	
	err = pg.Migrate()
	if err != nil {
		log.Fatal("Ошибка выполнения миграций: %v", err)
	}
	log.Println("Подключение к PostgreSQL успешно установлено!")
}
```

## Пример структуры миграций
Миграции должны находиться в папке migrations. Пример структуры файлов:

```bash
migrations/
├── 0001_initial.sql
├── 0002_add_users_table.sql
└── 0003_add_indexes.sql
```

### sql
```Bash
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE;
    
    -- +migrate Down
    DROP TABLE IF EXISTS users;
);
```