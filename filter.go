package generator

// LoopDirective is usually used as return value
// in Yield operations.
//
// LoopDirective can be used to make master code
// abort if slave demands it.
type LoopDirective uint

// LoopDirective can be one of the following instances:
const (
	LoopDirectiveBreak    = iota // recommend master to leave loop
	LoopDirectiveContinue        // recommend master to produce more data items
)

// Predicate represent a predicate on an item.
type Predicate[T any] func(elem T) bool

// Filter is a generator map for filtering out items in streams
// which does not satisfy a predicate.
func Filter[T any](predicate Predicate[T]) GeneratorMap[LoopDirective, T] {
	return CreateIterativeMap(
		func(
			outStream Generator[LoopDirective, T],
			elem T) LoopDirective {

			if predicate(elem) {
				return outStream.Yield(elem)
			}
			return LoopDirectiveContinue
		})
}
