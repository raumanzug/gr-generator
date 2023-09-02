// Package gr-generator implements generator support.
//
// This package provides generators for Go as they already exist
// in some other programming languages as C# and Python.  Generators
// handle possibly infinite streams of data via a master-slave approach.
// There is code serving as master producing sequences,
// e.g. natural numers, 0, 1, 2, 3, ...  Slave code reads in item by item
// and process them.
//
// Example slave code:
//
//	import "fmt"
//
//	 ... // a lot of code
//
//	Foreach(
//	   NewNaturals(0),
//	   func (elem uint) LoopDirective {
//	      fmt.Println(elem)
//	      return LoopDirectiveContinue
//	   })
//
// This code prints each natural number.
//
//	NewNaturals(0)
//
// denotes the master which produces the sequence of natural numbers.
// The function below treats item by item.  Its return type [LoopDirective]
// can have two instances:
//   - [LoopDirectiveBreak] recommends the master to
//     stop producing items.
//   - [LoopDirectiveContinue] recommends the master
//     to hand over the next item.
//
// Slave functions do not need to return with
// [LoopDirective] return code.  Other types are possible.
//
// Now, turn our attention to providing master code:  [NewNaturals]
// can be introduced as follows:
//
//	type naturals struct {
//	   GeneratorBase[LoopDirective, uint]
//	   start uint
//	}
//
//	func (n *naturals) Loop() {
//	   counter :=  n.start
//
//	   for {
//	      if LoopDirectiveBreak == n.Yield(counter) {
//	         break
//	      }
//	      counter++
//	   }
//	}
//
//	func NewNaturals(start uint) Generator[LoopDirective, uint] {
//	   return &naturals(start: start)
//	}
//
// Master implements interface Generator[R, T].  Herein R denotes the
// return type of slave functions and T denotes the type of the items.
// Aforementioned type naturals implement this interface.  It embeds
// GeneratorBase[R, T]. We implement the Loop method which contains the
// master code.  n.Yield(counter) herein sends counter as a data item
// to slave and it replies with an item of type R which master
// processes respectively.
//
// # Generator maps.
//
// [gr-generator] supports generator maps, i.e. maps which eat generators and
// produce other generators from them.  An example for it are filters.
// The following filter filters out even numbers from generator:
//
//	Filter[uint](isNonDivisible(2))
//
// This map can be applied as follows:
//
//	numbers := NewNaturals(0)
//	myFilter := Filter[uint](isNonDivisible(2))
//
//	oddNumbers :=  myFilter.Apply(numbers)
//
// oddNumbers generates odd numbers by extracting them from all
// natural numbers.
//
//	isNonDivisible(2)
//
// is a predicate, a function defined as follows:
//
//	func isNonDivisible(divisor uint) Predicate[uint] {
//		return func(
//			elem uint) bool {
//			return elem%divisor != 0
//		}
//	}
//
// # Defining generator maps - the iterative method
//
// This package introduces two methods of defining generator maps.
// The iterative way is the more easy way.  The following code
// shows how the aforementioned Filter map can be defined in this
// manner:
//
//	func Filter[T any](predicate Predicate[T]) GeneratorMap[LoopDirective, T] {
//		return CreateIterativeMap(
//			func(
//				outStream Generator[LoopDirective, T],
//				elem T) LoopDirective {
//
//				if predicate(elem) {
//					return outStream.Yield(elem)
//				}
//				return LoopDirectiveContinue
//			})
//	}
//
// [CreateIterativeMap] produces generator maps in iterative way.
// Two generators play a role herein.  This code is slave for
// the inStream generator and the master of outStream generator.
// [CreateIterativeMap] is set up similarly to aforementioned [Foreach]
// statement but also contains outStream.Yield(elem) statement
// which sends elem to outStream.
//
// # Defining generator maps - the recursive way
//
// In iteratively defined generator maps generator maps process
// item by item sent from inStream.  However, [gr-generator]
// package provides another way to define generator maps.  The
// recursive one.
//
// Eratosthenes' sieve should serve as an example.  This approach
// to enumerate prime numbers start with the sequence of natural
// numbers starting from 2.
//
//	startStream :=  NewNaturals(2)
//
// A generator map
//
//	primeMap()
//
// to be defined extract the prime numbers
// from it.  At first we split the stream in head and tail.  The first
// number in startStream, 2, is the head, the remaining numbers
// beginning from 3 is the tail.  A function eating the head has to be
// defined which generates a return value from return type R of
// GeneratorMap[R, T] and a sequence of GeneratorMap[R, T] values.
// These generator maps will be applied to the tail of startStream after
// executing the defined function.  That's the recursive way of
// defining generator maps.  The following code shows how to do it:
//
//	func primeMap() GeneratorMap[LoopDirective, uint] {
//	   return CreateRecursiveMap(
//	      func(
//	         outStream Generator[LoopDirective, uint],
//	         elem uint) (
//	         LoopDirective,
//	         []GeneratorMap[LoopDirective, uint]) {
//
//	         return outStream.Yield(elem),
//	            []GeneratorMap[LoopDirective, uint]{
//	               Filter[uint](isNonDivisible(elem)),
//	               primeMap()}
//	      })
//	}
//
// [CreateRecursiveMap] is the function which performs recursive
// definition of generator maps.  We see, that the first
// return value, outStream.Yield(elem), send the head, i.e. 2,
// to its slave.  2 is the first prime number.  After doing that
// two generator maps are applied to the tail.  At first,
// Filter[uint](isNonDivisible(elem)) drops all even numbers so
// that only odd numbers remain.  We obtain the sequence
// 3, 5, 7, 9, 11, ...  The next generator map obtained by the
// defining function is primeMap itself, i.e. the process whill
// be repeated with this sequence.  3 will be sent after sending 2.
// 3 is the next prime number.  Then all numbers divisable by
// 3 are dropped in the sequence coming after 3 so that we obtain
// the sequence 5, 7, 11, 13, 17, ... and so on and so forth.
//
// The following code defines a generator, primes, generating
// all prime numbers:
//
//	myMap :=  primeMap()
//	primes := myMap.Apply(startStream)
//
// Be aware that this implementation is far from efficient.
// It should serve as a demonstration sample introducing
// recursive defining generator map.
package generator
