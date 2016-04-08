# equirecursive-golang

Brief demonstration of Go's equirecursive types. One little-known aspect of Go's
type system is that it support
[equirecursive](https://en.wikipedia.org/wiki/Recursive_data_type#Equirecursive_types)
types. The only place that I have seen this in use is with state machine-like
computations, with types like:
~~~~~{.go}
type State func(Label)State
~~~~~

This repo shows how the subset of go *without* looping constructs or explicit
recursion (or if) but with recursion is Turing-complete, by showing how to do
some standard computations with the un(i)typed lambda calculus.
