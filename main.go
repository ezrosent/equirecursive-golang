package main

import "fmt"

type T func(int) int
type L func(L) L

// infinite loop without recursion!
func Omega() L {
	omega := func(x L) L { return x(x) }
	return omega(omega)
}

func Tests(oneF func(chan int) L) {

	curry := func(f func(L, L) L) L {
		return func(a L) L { return func(b L) L { return f(a, b) } }
	}
	ch := make(chan int, 1)
	counter := oneF(ch)
	ident := func(l L) L { return l }
	// Church bools
	tru := curry(func(a L, b L) L { return a })
	fals := curry(func(a L, b L) L { return b })
	iif := func(c, a, b L) L { return c(a)(b) }
	iif(tru, fals, fals)

	// Church numerals
	zero := fals
	one := curry(func(f L, a L) L { return f(a) })
	plus := curry(func(m L, n L) L { return curry(func(f L, x L) L { return m(f)(n(f)(x)) }) })
	plus(one)(zero)
	// off by two, but this works
	two := plus(one)(plus(one)(one))
	four := plus(two)(two)
	eight := plus(four)(four)
	mfour := iif(tru, eight, two)

	go mfour(counter)(ident)
	sum := 0
	for i := range ch {
		sum += i
		fmt.Println(i, sum)
	}
}

// The following would work if go evaluated things lazily
// it causes infinite loops though
// func YL(f func(T)T) T {
//   return f(YL(f))
// }

func YRec(f func(T) T) T {
	return func(x int) int {
		return f(YRec(f))(x)
	}
}

func Factorial(f T) T {
	return func(i int) int {
		if i <= 0 {
			return 1
		} else {
			return i * f(i-1)
		}
	}
}

func Summer(ch chan int) L {
	return func(l L) L {
		ch <- 1
		return Summer(ch)
	}
}

func Base(ch chan int) L {
	return func(l L) L {
		ch <- 0
		return Base(ch)
	}
}

func main() {

	fact := YRec(Factorial)
	fmt.Println(fact(4))
	// Omega() stack overflow
	Tests(Summer)
}
