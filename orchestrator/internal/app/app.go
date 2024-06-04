package app

import (
	"time"
	"math"
	"strconv"
	"unicode"
	"os"
	"log"

	"dac/orchestrator/models"
)

// ProcessExpression:
// Токенизирует выражение и преобразует его в постфиксную нотацию.
// Передает постфиксное выражение в computePostfix для вычисления.
func ProcessExpression(expr models.Expression) {
	tokens := tokenize(expr.Expression)
	postfix := infixToPostfix(tokens)
	computePostfix(postfix, expr.ID)
}

// computePostfix:
// Проходит по каждому токену постфиксного выражения.
// Если токен является числом, помещает его в стек.
// Если токен является оператором, извлекает два числа из стека, создает задачу и отправляет ее в канал задач.
// Ожидает результат выполнения задачи агентом и помещает его обратно в стек.
// В конце обновляет результат и статус выражения в модели.
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

/*
	Почему waitForResult важная функция?

	Если убрать waitForResult, оркестратор не сможет дождаться выполнения задачи агентом 
	перед продолжением обработки следующего оператора. 
	Следовательно, цепочка вычислений прервется на первой задаче, что приводит к некорректному результату.

	Но это не мешает агенту корректно выполнять задачи и отдавать их результат обратно и также без агента
	оркестратор работать не будет
*/
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
