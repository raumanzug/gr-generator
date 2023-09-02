package generator

import (
	"testing"
)

func Test_ArrayGenerator(t *testing.T) {
	data := []uint{1, 2, 3}
	generator := Array2Generator(data)

	fstRun := true
	var result uint = 0

	Foreach(
		generator,
		func(elem uint) LoopDirective {
			if fstRun {
				fstRun = false
				return LoopDirectiveContinue
			}
			result = elem
			return LoopDirectiveBreak
		})

	if result != 2 {
		t.Fatalf("expected: 2, got: %d\n", result)
	}
}
