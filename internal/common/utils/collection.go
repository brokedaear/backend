// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package utils

// Find is a higher-order function that iterates through a collection to find
// a value. If the value is found, it returns the value and a boolean
// that signifies the value exists in the collection. Otherwise, the zero
// value of the value is returned along with false.
//
// Retrieved from: https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/revisiting-arrays-and-slices-with-generics.
func Find[A any](collection []A, finder func(A, A) bool, target A) (A, bool) {
	for _, item := range collection {
		if finder(item, target) {
			return item, true
		}
	}

	var zero A
	return zero, false
}

// Reduce takes a collection of elements A and applies a reduction function
// on the elements, reducing the collection into a single value of type B.
//
// Retrieved from: https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/revisiting-arrays-and-slices-with-generics.
func Reduce[A, B any](collection []A, reductionFunc func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, item := range collection {
		result = reductionFunc(result, item)
	}
	return result
}
