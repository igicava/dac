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

    expr.Status = "pending"

    models.Mu.Lock()
    err := models.Add(expr.ID, expr)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    models.Mu.Unlock()

    go app.ProcessExpression(expr)

    w.WriteHeader(http.StatusCreated)
}

// GetExpressions возвращает все выражения
func GetExpressions(w http.ResponseWriter, r *http.Request) {
    models.Mu.Lock()
    defer models.Mu.Unlock()

    exprs, err := models.SelectExpressions(context.TODO(), models.DB)
    if err != nil {
        log.Printf("Error handler.go 42 : %s", err)
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "expressions": exprs,
    })
}

// GetExpressionByID возвращает выражение по его идентификатору
func GetExpressionByID(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    models.Mu.Lock()
    expr, ok := models.SelectExpressionByID(context.TODO(), models.DB, id)
    models.Mu.Unlock()

    if ok != nil {
        log.Printf("Error handler.go 58: %s", ok)
        http.Error(w, "Expression not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "expression": expr,
    })
}