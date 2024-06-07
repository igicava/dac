package models

import (
	"sync"
)

var Expressions = make(map[string]Expression) // Мапа c выражениями
var Tasks = make(chan Task) // Канал с тасками
var Results = make(map[string]chan Result) // Канал с результатами
var Mu sync.Mutex // Мьютекс

func NewChan(id string) {
    Mu.Lock()
    Results[id] = make(chan Result)
    Mu.Unlock()
}

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

// Структура результатов
type Result struct {
    ID     string  `json:"id"`
    Result float64 `json:"result"`
}
