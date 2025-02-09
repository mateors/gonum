// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build safe
// +build safe

// TODO(kortschak): Get rid of this rigmarole if https://golang.org/issue/50118
// is accepted.

package r3

import (
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

type array [9]float64

// At returns the value of a matrix element at row i, column j.
// At expects indices in the range [0,2].
// It will panic if i or j are out of bounds for the matrix.
func (m *Mat) At(i, j int) float64 {
	if uint(i) > 2 {
		panic(mat.ErrRowAccess)
	}
	if uint(j) > 2 {
		panic(mat.ErrColAccess)
	}
	if m.data == nil {
		m.data = new(array)
	}
	return m.data[i*3+j]
}

// Set sets the element at row i, column j to the value v.
func (m *Mat) Set(i, j int, v float64) {
	if uint(i) > 2 {
		panic(mat.ErrRowAccess)
	}
	if uint(j) > 2 {
		panic(mat.ErrColAccess)
	}
	if m.data == nil {
		m.data = new(array)
	}
	m.data[i*3+j] = v
}

// Eye returns the 3×3 Identity matrix
func Eye() *Mat {
	return &Mat{&array{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}}
}

// Skew returns the 3×3 skew symmetric matrix (right hand system) of v.
//                  ⎡ 0 -z  y⎤
//  Skew({x,y,z}) = ⎢ z  0 -x⎥
//                  ⎣-y  x  0⎦
func Skew(v Vec) (M *Mat) {
	return &Mat{&array{
		0, -v.Z, v.Y,
		v.Z, 0, -v.X,
		-v.Y, v.X, 0,
	}}
}

// Mul takes the matrix product of a and b, placing the result in the receiver.
// If the number of columns in a does not equal 3, Mul will panic.
func (m *Mat) Mul(a, b mat.Matrix) {
	ra, ca := a.Dims()
	rb, cb := b.Dims()
	switch {
	case ra != 3:
		panic(mat.ErrShape)
	case cb != 3:
		panic(mat.ErrShape)
	case ca != rb:
		panic(mat.ErrShape)
	}
	if m.data == nil {
		m.data = new(array)
	}
	if ca != 3 {
		// General matrix multiplication for the case where the inner dimension is not 3.
		var t mat.Dense
		t.Mul(a, b)
		copy(m.data[:], t.RawMatrix().Data)
		return
	}

	a00 := a.At(0, 0)
	b00 := b.At(0, 0)
	a01 := a.At(0, 1)
	b01 := b.At(0, 1)
	a02 := a.At(0, 2)
	b02 := b.At(0, 2)
	a10 := a.At(1, 0)
	b10 := b.At(1, 0)
	a11 := a.At(1, 1)
	b11 := b.At(1, 1)
	a12 := a.At(1, 2)
	b12 := b.At(1, 2)
	a20 := a.At(2, 0)
	b20 := b.At(2, 0)
	a21 := a.At(2, 1)
	b21 := b.At(2, 1)
	a22 := a.At(2, 2)
	b22 := b.At(2, 2)
	*(m.data) = array{
		a00*b00 + a01*b10 + a02*b20, a00*b01 + a01*b11 + a02*b21, a00*b02 + a01*b12 + a02*b22,
		a10*b00 + a11*b10 + a12*b20, a10*b01 + a11*b11 + a12*b21, a10*b02 + a11*b12 + a12*b22,
		a20*b00 + a21*b10 + a22*b20, a20*b01 + a21*b11 + a22*b21, a20*b02 + a21*b12 + a22*b22,
	}
}

// RawMatrix returns the blas representation of the matrix with the backing
// data of this matrix. Changes to the returned matrix will be reflected in
// the receiver.
func (m *Mat) RawMatrix() blas64.General {
	if m.data == nil {
		m.data = new(array)
	}
	return blas64.General{Rows: 3, Cols: 3, Data: m.data[:], Stride: 3}
}

// BUG(kortschak): Implementing NewMat without unsafe conversion or slice to
// array pointer conversion leaves it with semantics that do not match the
// sharing semantics that exist in the mat package.
func arrayFrom(vals []float64) *array {
	// TODO(kortschak): Use array conversion when go1.16 is no longer supported.
	var a array
	copy(a[:], vals)
	return &a
}
