/// Most examples taken from wikipedia articles on fixed-point combinators and church encoding
package main

import "fmt"

// type of all values in the lambda calculus
type L func(L) L

// infinite loop without recursion!
func Omega() L {
	omega := func(x L) L { return x(x) }
	return omega(omega)
}

// convert a church numeral into an integer by counting the number of times a function is applied.
// There doesn't appear to be a nicer way to do this, unfortunately
func ChurchToInt(n L) int {
	i := 0
	n(func(l L) L {
		i += 1
		return l
	})(Ident)
	return i
}

// adaptation from the haskell function
func curry(f func(L, L) L) L {
	return func(a L) L {
		return func(b L) L {
			return f(a, b)
		}
	}
}

var One L = curry(func(f, x L) L { return f(x) })

// add two church numerals
func Plus(m, n L) L {
	return curry(func(f L, x L) L {
		return m(f)(n(f)(x))
	})
}

func Pred(n L) L {
	return curry(
		func(f, x L) L {
			return n(
				curry(
					func(g, h L) L {
						return h(g(f))
					}))(
				func(u L) L { return x })(
				func(u L) L {
					return u
				})
		})
}

// dummy value
func Ident(f L) L {
	return f
}

// strict fixed-point combinator without recursion
func Z(f L) L {
	help := func(x L) L {
		return f(func(v L) L {
			return x(x)(v)
		})
	}
	return help(help)
}


func Examples() {

	// Church bools
	tru := curry(func(a L, b L) L { return a })
	fals := curry(func(a L, b L) L { return b })
	iif := func(c, a, b L) L { return c(a)(b) }

	iif(tru, fals, fals)
	iszero := func(n L) L { return n(func(x L) L { return fals })(tru) }

	// Church numerals
	zero := fals
	s := func(l L) L { return Plus(One, l) }

	two := s(One)
	four := s(s(two))
	eight := Plus(four, four)
	eight_if := iif(iszero(zero), eight, two)
	fmt.Println(ChurchToInt(zero), ChurchToInt(One), ChurchToInt(Pred(eight_if)))

	// for lazy computations that can be evaluated strictly
	pad := func(f L) L { return func(g L) L { return f } }
	// recusively compute (sum_{i=0}^n i)
	summorial := Z(func(f L) L {
		return func(n L) L {
			return iif(iszero(n), pad(zero),
				func(g L) L { // need to bind this lazily
					return Plus(n, f(Pred(n)))
				})(Ident) // force the thunk
		}
	})
	fmt.Println(ChurchToInt(summorial(eight))) // 36
}

func main() {
	// Omega() stack overflow
	Examples()
}
