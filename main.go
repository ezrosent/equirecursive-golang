package main

import "fmt"

type L func(L) L

// strict fixed-point combinator without recursion
func Z(f L) L {
	help := func(x L) L {
		return f(func(v L) L {
			return x(x)(v)
		})
	}
	return help(help)
}

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
	})(func(l L) L { return l })
	return i
}
func curry(f func(L, L) L) L {
	return func(a L) L {
		return func(b L) L {
			return f(a, b)
		}
	}
}
func plus(m, n L) L {
	return curry(func(f L, x L) L {
		return m(f)(n(f)(x))
	})
}

func pred(n L) L {
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

func Tests() {

	// Church bools
	tru := curry(func(a L, b L) L { return a })
	fals := curry(func(a L, b L) L { return b })
	iif := func(c, a, b L) L { return c(a)(b) }

	iif(tru, fals, fals)
	iszero := func(n L) L { return n(func(x L) L { return fals })(tru) }

	// Church numerals
	zero := fals
	one := curry(func(f L, a L) L { return f(a) })
	s := func(l L) L { return plus(one, l) }
	// off by two, but this works
	lif := func(c L, a, b func() L) L {
		// need laziness here, more lambdas will implement it, but can't get it right atm
		if ChurchToInt(c(zero)(one)) == 0 {
			return c(a())(zero)
		} else {
			return c(zero)(b())
		}
	}

	two := s(one)
	four := plus(two, two)
	eight := plus(four, four)
	eight_if := iif(iszero(zero), eight, two)
	fmt.Println(ChurchToInt(zero), ChurchToInt(one), ChurchToInt(pred(eight_if)))
	summorial := Z(func(f L) L {
		return func(n L) L {
			return lif(iszero(n),
				func() L { return zero },
				func() L { return plus(n, f(pred(n))) })
		}
	})
	fmt.Println(ChurchToInt(summorial(eight)))
}

// The following would work if go evaluated things lazily
// it causes infinite loops though
// func YL(f func(T)T) T {
//   return f(YL(f))
// }

func main() {
	// Omega() stack overflow
	Tests()
}
