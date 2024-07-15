package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"dac/orchestrator/internal/app"
	"dac/orchestrator/models"
)

// AddExpression добавляет выражение
func AddExpression(w http.ResponseWriter, r *http.Request) {
    var expr models.Expression
    if err := json.NewDecoder(r.Body).Decode(&expr); err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }
    _, ok := models.Results[expr.ID]
    if ok {
        w.Write([]byte("Your ID is already occupied by another user. Come up with another one or wait for it to be released\n"))
        return
    }
    expr.Status = "pending"

    models.Mu.Lock()
    name, err := models.ValidToken(expr.Token)
    if err != nil {
        log.Println(err)
        return
    }
    err = models.Add(expr.ID, expr, name)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }

    models.Mu.Unlock()

    go app.ProcessExpression(expr, name)

    w.WriteHeader(http.StatusCreated)
}

// GetExpressions возвращает все выражения
func GetExpressions(w http.ResponseWriter, r *http.Request) {
    models.Mu.Lock()
    defer models.Mu.Unlock()
    
    var req models.Expression
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    name, err := models.ValidToken(req.Token)
    if err != nil {
        log.Println(err)
        return
    }

    exprs, err := models.SelectExpressions(context.TODO(), models.DB, name)
    if err != nil {
        log.Printf("Error handler.go : %s", err)
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "expressions": exprs,
    })
}

// GetExpressionByID возвращает выражение по его идентификатору
func GetExpressionByID(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    models.Mu.Lock()
    var req models.Expression
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    name, err := models.ValidToken(req.Token)
    if err != nil {
        log.Println(err)
        return
    }
    expr, ok := models.SelectExpressionByID(context.TODO(), models.DB, id, name)
    models.Mu.Unlock()

    if ok != nil {
        log.Printf("Error handler.go : %s", ok)
        http.Error(w, "Expression not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "expression": expr,
    })
}

// Регистрация нового пользователя
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
    var user models.UserModel
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        log.Printf("Error in handler.go : %s", err)
    }
    var q = `
	INSERT INTO users (name, password) values ($1, $2)
	`
	tx, err := models.DB.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Printf("Error in handler.go : %s", err)
	}
	result, err := tx.ExecContext(context.TODO(), q, user.Name, user.Password)
	if err != nil {
		log.Printf("Error in handler.go : %s", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		log.Printf("Error in handler.go : %s", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error in handler.go : %s", err)
	}

	w.WriteHeader(http.StatusOK)
}

// Вход для пользователей
func LoginUser(w http.ResponseWriter, r *http.Request) {
    var user models.UserModel
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        log.Printf("Error in handler.go : %s", err)
    }
    
    u, err := models.SelectUserByName(context.TODO(), models.DB, user.Name)
    if err != nil {
        log.Printf("Error handler.go : %s", err)
    }
    
    if u.Name == "" || u.Password == "" {
        log.Println("User not found")
        return
    } 
    
    if user.Password != u.Password {
        log.Println("Uncorrected password")
        return
    }

    token := models.NewToken(user.Name)
    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
    })
}