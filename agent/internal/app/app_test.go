package app

import (
	"testing"
)

func TestComputeTask(t *testing.T) {
	t.Run("ADDITION", func(t *testing.T) {
		expression := Task{
			Arg1: 2.2,
			Arg2: 2.2,
			Operation: "+",
		}
		test := ComputeTask(expression)
		if test != (2.2 + 2.2) {
			t.Fail()
		}
	})

	t.Run("SUBTRACTION", func(t *testing.T) {
		expression := Task{
			Arg1: 21,
			Arg2: 1,
			Operation: "-",
		}
		test := ComputeTask(expression)
		if test != (21 - 1) {
			t.Fail()
		}
	})

	t.Run("MULTIPLICATION", func(t *testing.T) {
		expression := Task{
			Arg1: 2.2,
			Arg2: 2,
			Operation: "*",
		}
		test := ComputeTask(expression)
		if test != (2.2 * 2) {
			t.Fail()
		}
	})

	t.Run("DIVISION", func(t *testing.T) {
		expression := Task{
			Arg1: 2.2,
			Arg2: 2,
			Operation: "/",
		}
		test := ComputeTask(expression)
		if test != (2.2 / 2) {
			t.Fail()
		}
	})

	t.Run("Division by zero", func(t *testing.T) {
		expression := Task{
			Arg1: 2.2,
			Arg2: 0,
			Operation: "/",
		}
		test := ComputeTask(expression)
		if test != 0 {
			t.Fail()
		}
	})

	t.Run("Unknown operation", func(t *testing.T) {
		expression := Task{
			Arg1: 2.2,
			Arg2: 0,
			Operation: "%",
		}
		test := ComputeTask(expression)
		if test != 0 {
			t.Fail()
		}
	})
}