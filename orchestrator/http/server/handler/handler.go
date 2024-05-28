package handler

import (
	"net/http"
	"encoding/json"
    
    "dac/orchestrator/models"
    "dac/orchestrator/internal/app"
)

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

func AgentTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		select {
		case task := <-models.Tasks:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"task": task,
			})
		default:
			http.Error(w, "No tasks available", http.StatusNotFound)
		}
	} else if r.Method == "POST" {
		var taskResult struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&taskResult); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	
		models.Mu.Lock()
		defer models.Mu.Unlock()
	
		expr, ok := models.Expressions[taskResult.ID]
		if !ok {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		expr.Result = taskResult.Result
		expr.Status = "completed"
		models.Expressions[taskResult.ID] = expr
	
		w.WriteHeader(http.StatusOK)
	}
    
}