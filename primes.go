package generator

type naturals struct {
	GeneratorBase[LoopDirective, uint]
	start uint
}

func (n *naturals) Loop() {
	counter := n.start

	for {
		if LoopDirectiveBreak == n.Yield(counter) {
			break
		}
		counter++
	}
}

// NewNaturals produces a generator generating
// a sequence of natural numbers start with start.
func NewNaturals(start uint) Generator[LoopDirective, uint] {
	return &naturals{start: start}
}

// IsNonDivisible produces a predicate determining
// whether its argument is not divisible by divisor
func IsNonDivisible(divisor uint) Predicate[uint] {
	return func(
		elem uint) bool {
		return elem%divisor != 0
	}
}

func primeIterMap() GeneratorMap[LoopDirective, uint] {
	alreadyFound := []uint{}
	var l uint = 2
	return CreateIterativeMap(
		func(
			outStream Generator[LoopDirective, uint],
			elem uint) LoopDirective {

			if elem > l*l {
				l++
			}
			isNextPrime := true
			for j := uint(0); j < uint(len(alreadyFound)); j++ {
				if alreadyFound[j] > l {
					break
				}
				if elem%alreadyFound[j] == 0 {
					isNextPrime = false
					break
				}
			}
			if isNextPrime {
				alreadyFound = append(alreadyFound, elem)
				if LoopDirectiveBreak == outStream.Yield(elem) {
					return LoopDirectiveBreak
				}
			}

			return LoopDirectiveContinue
		})
}

func primeRecMap(alreadyFound []uint, lastSent uint) GeneratorMap[LoopDirective, uint] {
	return CreateRecursiveMap(
		func(
			outStream Generator[LoopDirective, uint],
			elem uint) (
			LoopDirective,
			[]GeneratorMap[LoopDirective, uint]) {

			if LoopDirectiveBreak == outStream.Yield(elem) {
				return LoopDirectiveBreak,
					[]GeneratorMap[LoopDirective, uint]{}
			}

			_alreadyFound := append(alreadyFound, elem)
			_lastSent := lastSent
			maps := []GeneratorMap[LoopDirective, uint]{}
			if lastSent*lastSent < elem {
				_lastSent = _alreadyFound[0]
				maps = append(
					maps,
					Filter[uint](IsNonDivisible(_lastSent)))
				_alreadyFound = _alreadyFound[1:]
			}
			maps = append(
				maps,
				primeRecMap(_alreadyFound, _lastSent))
			return LoopDirectiveContinue, maps
		})

}

// Primes is a generator generating all prime numbers.
func Primes() Generator[LoopDirective, uint] {
	generator := NewNaturals(2)
	m := primeIterMap()
	return m.Apply(generator)
}

// RecPrimes is a generator generating all prime numbers.
//
// It is the same function as [Primes].  However we use
// recursive definition of generator map.
// RecPrimes is less performant than [Primes].
func RecPrimes() Generator[LoopDirective, uint] {
	generator := NewNaturals(2)
	m := primeRecMap([]uint{}, 0)
	return m.Apply(generator)
}
