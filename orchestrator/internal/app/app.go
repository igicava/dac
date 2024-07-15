package app

import (
	"log"
	"os"
	"strconv"
	"unicode"

	"dac/orchestrator/models"
)

// ProcessExpression:
// Токенизирует выражение и преобразует его в постфиксную нотацию.
// Передает постфиксное выражение в computePostfix для вычисления.
func ProcessExpression(expr models.Expression, user_name string) {
	tokens := tokenize(expr.Expression)
	postfix := infixToPostfix(tokens)
	computePostfix(postfix, expr.ID, user_name)
}

// computePostfix:
// Проходит по каждому токену постфиксного выражения.
// Если токен является числом, помещает его в стек.
// Если токен является оператором, извлекает два числа из стека, создает задачу и отправляет ее в канал задач.
// Ожидает результат выполнения задачи агентом и помещает его обратно в стек.
// В конце обновляет результат и статус выражения в модели.
func computePostfix(tokens []string, exprID string, user_name string) {
	var stack []float64
	models.NewChan(exprID, user_name)
	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
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
			out := <-models.Results[task.ID]
			stack = append(stack, out.Result)
		}
	}

	if len(stack) == 1 {
		models.Mu.Lock()
		defer models.Mu.Unlock()
		err := models.UpdateResult(exprID, stack[0])
		if err != nil {
			log.Printf("In app.go 58 : %s", err)
		}
		err = models.UpdateStatus(exprID)
		if err != nil {
			log.Printf("In app.go 62: %s", err)
		}
	}
}


// getOperationTime возвращает время выполнения для каждой операции
func getOperationTime(op string) int {
	switch op {
	case "+":
		return getEnvAsInt("TIME_ADDITION_MS", 1000)
	case "-":
		return getEnvAsInt("TIME_SUBTRACTION_MS", 1000)
	case "*":
		return getEnvAsInt("TIME_MULTIPLICATIONS_MS", 2000)
	case "/":
		return getEnvAsInt("TIME_DIVISIONS_MS", 2000)
	}
	return 0
}

// getEnvAsInt возвращает значение переменной среды или значение по умолчанию
func getEnvAsInt(name string, defaultValue int) int {
    valueStr := os.Getenv(name)
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        log.Printf("Invalid value for %s: %s. Using default: %d\n", name, valueStr, defaultValue)
        return defaultValue
    }
    return value
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
