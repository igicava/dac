package handler

import (
	"net/http"
	"encoding/json"
    
    "dac/orchestrator/models"
    "dac/orchestrator/internal/app"
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
    models.Expressions[expr.ID] = expr
    models.Mu.Unlock()

    go app.ProcessExpression(expr)

    w.WriteHeader(http.StatusCreated)
}

// GetExpressions возвращает все выражения
func GetExpressions(w http.ResponseWriter, r *http.Request) {
    models.Mu.Lock()
    defer models.Mu.Unlock()

    var exprs []models.Expression
    for _, expr := range models.Expressions {
        exprs = append(exprs, expr)
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "expressions": exprs,
    })
}

// GetExpressionByID возвращает выражение по его идентификатору
func GetExpressionByID(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    models.Mu.Lock()
    expr, ok := models.Expressions[id]
    models.Mu.Unlock()

    if !ok {
        http.Error(w, "Expression not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "expression": expr,
    })
}