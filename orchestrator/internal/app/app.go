package app

import (
	"math"
	"strconv"
	"time"
	"unicode"

	"dac/orchestrator/models"
)

func ProcessExpression(expr models.Expression) {
	tokens := tokenize(expr.Expression)
	postfix := infixToPostfix(tokens)
	computePostfix(postfix, expr.ID)
}

func computePostfix(tokens []string, exprID string) {
	var stack []float64
	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			opTime := getOperationTime(token)
			task := models.Task{
				ID:            exprID,
				Arg1:          a,
				Arg2:          b,
				Operation:     token,
				OperationTime: opTime,
			}
			models.Tasks <- task
			result := waitForResult(task)
			stack = append(stack, result)
		}
	}
	if len(stack) == 1 {
		models.Mu.Lock()
		expr := models.Expressions[exprID]
		expr.Result = stack[0]
		expr.Status = "completed"
		models.Expressions[exprID] = expr
		models.Mu.Unlock()
	}
}

func waitForResult(task models.Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		return task.Arg1 / task.Arg2
	}
	return math.NaN()
}

func getOperationTime(op string) int {
	switch op {
	case "+":
		return 5000
	case "-":
		return 5000
	case "*":
		return 7000
	case "/":
		return 10000
	}
	return 0
}

// tokenize разбивает строку выражения на токены (числа и операторы).
func tokenize(expr string) []string {
	var tokens []string
	var number string

	for _, char := range expr {
		if unicode.IsDigit(char) || char == '.' {
			number += string(char)
		} else if char == ' ' {
			continue
		} else {
			if number != "" {
				tokens = append(tokens, number)
				number = ""
			}
			tokens = append(tokens, string(char))
		}
	}
	if number != "" {
		tokens = append(tokens, number)
	}

	return tokens
}

// precedence возвращает приоритет оператора.
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

// infixToPostfix преобразует инфиксное выражение в постфиксное.
func infixToPostfix(tokens []string) []string {
	var result []string
	var stack []string

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			result = append(result, token)
		} else {
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				result = append(result, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		result = append(result, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return result
}
