package generator

// GeneratorBase type helps to implement
// Generator[R, T] interface by an
// implementation of Yield and setYield
// method.
//
// Embed this type into structs which should
// implement Generator[R, T] interface.
// Complete implementing Generator[R, T] by
// implementing Loop method.
type GeneratorBase[R any, T any] struct {
	yield func(elem T) R
}

// Generator represents master code
// for using generators.
//
// Loop contains master code.  It uses
// method Yield in order to send data items
// to slave code.
type Generator[R any, T any] interface {

	// Sends data item elem to slave code.
	// slave code replies with data of type R.
	// Use this method in implementation of
	// method Loop.
	Yield(elem T) R

	setYield(func(elem T) R)

	// Contains master code.
	Loop()
}

func (b *GeneratorBase[R, T]) Yield(elem T) R {
	return b.yield(elem)
}

func (b *GeneratorBase[R, T]) setYield(yield func(elem T) R) {
	b.yield = yield
}

// Foreach defines slave code.
//
// Foreach demands for two parameters:
//   - g     master code
//   - yield slave code
func Foreach[R any, T any](g Generator[R, T], yield func(elem T) R) {
	g.setYield(yield)
	g.Loop()
}

// GeneratorMap represent generator maps, i.e.
// maps eating generators and producing generators.
//
// Apply method executes the such a generator transformation
type GeneratorMap[R any, T any] interface {
	Apply(g Generator[R, T]) Generator[R, T]
}

type sGenerator[R any, T any] struct {
	GeneratorBase[R, T]
	inStream Generator[R, T]
}

func (g *sGenerator[R, T]) Loop() {
	g.inStream.Loop()
}

type mapGenerator[R any, T any] func(inStream, outStream Generator[R, T], elem T) R

func (g mapGenerator[R, T]) Apply(inStream Generator[R, T]) Generator[R, T] {
	outStream := &sGenerator[R, T]{inStream: inStream}
	_action := func(elem T) R {
		return g(inStream, outStream, elem)
	}
	inStream.setYield(_action)
	return outStream
}

// CreateIterativeMap defines generator map in
// iterative manner.
func CreateIterativeMap[R any, T any](
	action func(
		outStream Generator[R, T],
		elem T) R) GeneratorMap[R, T] {
	var retval mapGenerator[R, T] = func(_, outStream Generator[R, T], elem T) R {
		return action(outStream, elem)
	}
	return retval
}

// CreateRecursiveMap defines generator map using
// recursive way.
func CreateRecursiveMap[R any, T any](
	action func(
		outStream Generator[R, T],
		elem T) (
		R,
		[]GeneratorMap[R, T])) GeneratorMap[R, T] {
	var retval mapGenerator[R, T] = func(inStream, outStream Generator[R, T], elem T) R {
		_retval, maps := action(outStream, elem)

		_stream := inStream
		for _, m := range maps {
			_stream = m.Apply(_stream)
		}
		_stream.setYield(outStream.Yield)

		return _retval
	}
	return retval
}
