package models

import (
	"context"
	"database/sql"
	"sync"

    _ "github.com/mattn/go-sqlite3"
)

var Tasks = make(chan Task)                   // Канал с тасками
var Results = make(map[string]chan Result)    // Канал с результатами
var Mu sync.Mutex                             // Мьютекс
var DB *sql.DB                                // Для доступа к БД

// Структура самого выражения
type Expression struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

// Структура таска
type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type ExpressionModel struct {
	ID           int
	ExpressionID string
	Expression   string
	UserID       int
	Status       string
	Result       float64
}

type UserModel struct {
	ID   int
	Name string
}

// Структура результатов
type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

// Новый канал для выражения
func NewChan(id string) {
	Mu.Lock()
	Results[id] = make(chan Result)
	Mu.Unlock()
}

// Открываем БД
func OpenDB() {
	ctx := context.TODO()
	var err error

	DB, err = sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}

	err = DB.PingContext(ctx)
	if err != nil {
		panic(err)
	}
}

// Создание таблиц
func СreateTables(ctx context.Context) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT
	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        expression_id TEXT NOT NULL, 
		expression TEXT NOT NULL,
		user_id INTEGER NOT NULL,
        status TEXT NOT NULL,
        result REAL,
	
		FOREIGN KEY (user_id)  REFERENCES expressions (id)
	);`
	)

	if _, err := DB.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	if _, err := DB.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}

// Добавление выражения
func Add(expr_id string, expr Expression) error {
	var q = `
	INSERT INTO expressions (expression_id, expression, user_id, status, result) values ($1, $2, $3, $4, $5)
	`
	tx, err := DB.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	result, err := tx.ExecContext(context.TODO(), q, expr_id, expr.Expression, "1", "pending", 0)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Обновить результат
func UpdateResult(id string, result float64) error {
	var q = "UPDATE expressions SET result = $1 WHERE expression_id = $2"
	tx, err := DB.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(context.TODO(), q, result, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Обновить статус
func UpdateStatus(id string) error {
	var q = "UPDATE expressions SET status = $1 WHERE expression_id = $2"
	tx, err := DB.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(context.TODO(), q, "completed", id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Выражение по ID
func SelectExpressionByID(ctx context.Context, db *sql.DB, id string) (Expression, error) {
	u := Expression{}
	var q = "SELECT expression_id, expression, status, result FROM expressions WHERE expression_id = $1"
	err := db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Expression, &u.Status, &u.Result)
	if err != nil {
		return u, err
	}

	return u, nil
}

// Все выражения
func SelectExpressions(ctx context.Context, db *sql.DB) ([]Expression, error) {
	var expressions []Expression
	var q = "SELECT expression_id, expression, status, result FROM expressions"

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.ID, &e.Expression, &e.Status, &e.Result)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}