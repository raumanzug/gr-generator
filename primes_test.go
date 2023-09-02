package generator

import (
	"testing"
)

func TestIterPrimes(t *testing.T) {
	loopNr := 1
	var result uint

	Foreach(
		Primes(),
		func(elem uint) LoopDirective {
			if loopNr >= 12 {
				result = elem
				return LoopDirectiveBreak
			}
			loopNr++
			return LoopDirectiveContinue
		})

	if result != 37 {
		t.Fatalf("primes does not work.")
	}
}

func TestRecPrimes(t *testing.T) {
	loopNr := 1
	var result uint

	Foreach(
		RecPrimes(),
		func(elem uint) LoopDirective {
			if loopNr >= 12 {
				result = elem
				return LoopDirectiveBreak
			}
			loopNr++
			return LoopDirectiveContinue
		})

	if result != 37 {
		t.Fatalf("primes does not work.")
	}
}
