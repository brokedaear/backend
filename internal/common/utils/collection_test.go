// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package utils_test

import (
	"testing"

	"backend.brokedaear.com/internal/common/utils"
	"backend.brokedaear.com/pkg/assert"
)

func TestFind(t *testing.T) {
	t.Run(
		"name string", func(t *testing.T) {
			fridayNightSquad := []string{"neo", "luke", "kenny", "ian", "trevor"}
			sieve := func(x, y string) bool {
				return x == y
			}
			want := "ian"

			v, ok := utils.Find(fridayNightSquad, sieve, want)

			assert.Equal(t, true, ok)
			assert.Equal(t, v, want)
		},
	)

	type player struct {
		kills  uint
		deaths uint
	}

	ian := player{kills: 41, deaths: 1}
	neo := player{kills: 4, deaths: 4}
	luke := player{kills: 6, deaths: 4}
	kenny := player{kills: 24, deaths: 4}
	trevor := player{kills: 3, deaths: 5}

	t.Run(
		"best player", func(t *testing.T) {
			fridayNightSquad := []player{neo, luke, ian, kenny, trevor}
			hasMostCrackedOut := func(x, y player) bool {
				kd := x.kills / x.deaths
				target := y.kills / y.deaths
				return kd == target
			}

			t.Run(
				"exists", func(t *testing.T) {
					v, ok := utils.Find(fridayNightSquad, hasMostCrackedOut, ian)

					assert.True(t, ok)
					assert.Equal(t, v, ian)
				},
			)

			fridayNightSquad = []player{neo, luke, kenny, trevor}

			t.Run(
				"does not exist", func(t *testing.T) {
					v, ok := utils.Find(fridayNightSquad, hasMostCrackedOut, ian)

					assert.False(t, ok)
					assert.NotEqual(t, v, ian)
				},
			)
		},
	)
}

func TestReduce(t *testing.T) {
	t.Run(
		"multiplication of all elements", func(t *testing.T) {
			multiply := func(x, y int) int {
				return x * y
			}

			assert.Equal(t, utils.Reduce([]int{1, 2, 3}, multiply, 1), 6)
		},
	)

	t.Run(
		"concatenate strings", func(t *testing.T) {
			concatenate := func(x, y string) string {
				return x + y
			}

			assert.Equal(t, utils.Reduce([]string{"a", "b", "c"}, concatenate, ""), "abc")
		},
	)

	t.Run(
		"add up items in cart", func(t *testing.T) {
			addItemValues := func(x int, y mockShopItem) int {
				return x + y.price
			}

			cart := []mockShopItem{
				{"microwave", 30},
				{"daw candy pack", 70},
				{"tuna fish", 10},
			}

			assert.Equal(t, utils.Reduce(cart, addItemValues, 0), 110)
		},
	)
}

type mockShopItem struct {
	name  string
	price int
}
