package generator

type arrayGenerator[T any] struct {
	GeneratorBase[LoopDirective, T]
	data []T
}

func (g *arrayGenerator[T]) Loop() {
	for _, item := range g.data {
		if LoopDirectiveBreak == g.Yield(item) {
			break
		}
	}
}

// Array2Generator transform a slice to a generator.
func Array2Generator[T any](data []T) Generator[LoopDirective, T] {
	return &arrayGenerator[T]{data: data}
}
