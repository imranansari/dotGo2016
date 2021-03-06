// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "fmt"

const boundsChecks = true

type T float64 // for convenience

type Vector struct {
	array       []T // may be longer than len
	len, stride int
}

func (x *Vector) addr(i int) *T {
	if boundsChecks && uint(i) >= uint(x.len) {
		panic("index out of bounds")
	}
	return &x.array[i*x.stride]
}

func (x *Vector) Len() int        { return x.len }
func (x *Vector) [] (i int) T     { return *x.addr(i) }
func (x *Vector) []= (i int, t T) { *x.addr(i) = t }

// dot-product
func (x *Vector) * (y *Vector) T {
	if x.Len() != y.Len() {
		panic("incompatible vector lengths")
	}
	var t T
	for i := x.Len() - 1; i >= 0; i-- {
		t += x[i] * y[i]
	}
	return t
}

func (x *Vector) GoSlice() []T {
	if x.stride == 1 {
		return x.array[:x.len]
	}
	s := make([]T, x.len)
	j := 0
	for i := range s {
		s[i] = x.array[j]
		j += x.stride
	}
	return s
}

type dim [2]int

func (d dim) transpose() dim { return dim{d[1], d[0]} }

type Matrix struct {
	array       []T
	len, stride dim
}

func (m *Matrix) addr(i, j int) *T {
	if boundsChecks && uint(i) >= uint(m.len[0]) || uint(j) >= uint(m.len[1]) {
		panic("index out of bounds")
	}
	return &m.array[i*m.stride[0]+j*m.stride[1]]
}

func (m *Matrix) Len() (int, int)    { return m.len[0], m.len[1] }
func (m *Matrix) [] (i, j int) T     { return *m.addr(i, j) }
func (m *Matrix) []= (i, j int, x T) { *m.addr(i, j) = x }

func (m *Matrix) Row(i int) *Vector { return &Vector{m.array[i*m.stride[0]:], m.len[1], m.stride[1]} }
func (m *Matrix) Col(j int) *Vector { return &Vector{m.array[j*m.stride[1]:], m.len[0], m.stride[0]} }

func (a *Matrix) Transpose() *Matrix {
	return &Matrix{
		a.array,
		a.len.transpose(),
		a.stride.transpose(),
	}
}

func NewMatrix(n, m int) *Matrix {
	if n < 0 || m < 0 {
		panic("invalid length")
	}
	return &Matrix{
		array:  make([]T, n*m),
		len:    dim{n, m},
		stride: dim{m, 1}, // row-major
	}
}

func (a *Matrix) Set(coeff ...T) {
	n, m := a.Len()
	if len(coeff) != n*m {
		panic("incorrect number of coefficients")
	}
	k := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			a[i, j] = coeff[k]
			k++
		}
	}
}

func (a *Matrix) Print() {
	n, m := a.Len()
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Printf(" %5g", a[i, j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func (a *Matrix) * (b *Matrix) *Matrix {
	n, m := a.Len()
	o, p := b.Len()
	if m != o {
		panic("incompatible matrix sizes")
	}
	c := NewMatrix(n, p)
	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			var t T
			for k := 0; k < m; k++ {
				t += a[i, k] * b[k, j]
			}
			c[i, j] = t
		}
	}
	return c
}

func (a *Matrix) Mul(b *Matrix) *Matrix {
	n, m := a.Len()
	o, p := b.Len()
	if m != o {
		panic("incompatible matrix sizes")
	}
	c := NewMatrix(n, p)
	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			a := a.Row(i)
			b := b.Col(j)
			c[i, j] = a * b
		}
	}
	return c
}

func main() {
	a := NewMatrix(4, 5)
	a.Set(
		4, 2, 7, 9, 1,
		5, 0, 1, 8, 3,
		5, 6, 3, 2, 1,
		7, 9, 0, 1, 2,
	)

	b := NewMatrix(5, 3)
	b.Set(
		3, 4, 5,
		0, 3, 1,
		3, 2, 1,
		8, 2, 6,
		2, 7, 1,
	)

	(a * b).Print()

	c := a.Mul(b)
	c.Print()

	c.Transpose().Print()

	n, m := c.Len()
	for i := 0; i < n; i++ {
		fmt.Println(c.Row(i).GoSlice())
	}

	for j := 0; j < m; j++ {
		fmt.Println(c.Col(j).GoSlice())
	}
}
