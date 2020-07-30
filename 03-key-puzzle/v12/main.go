package main

import "fmt"

// Solution is a wrapper object for puzzle solution
type Solution struct {
	A int
	B int
	C int
	X int
	Y int
}

func main() {
	puzzle := [6][8]bool{
		{false, false, false, false, false, false, false, false},
		{false, true, true, true, true, true, true, false},
		{false, true, false, false, false, true, true, false},
		{false, true, true, true, false, true, false, false},
		{false, true, false, true, true, true, true, false},
		{false, false, false, false, false, false, false, false},
	}

	// Set starting position
	startCol := 1
	startRow := 4

	// Setup solutions
	solutions := make([]Solution, 0)

	y := startRow
	x := startCol

	a := 0
	for {
		a++

		x = startCol
		y = startRow - a
		if y < 0 || !puzzle[y][x] {
			break
		}

		b := 0
		for {
			b++

			x = startCol + b
			y = startRow - a
			if x >= len(puzzle[0]) || !puzzle[y][x] {
				break
			}

			c := 0
			for {
				c++

				x = startCol + b
				y = startRow - a + c
				if y >= len(puzzle) || !puzzle[y][x] {
					break
				}

				solutions = append(solutions, Solution{
					A: a,
					B: b,
					C: c,
					X: x,
					Y: y,
				})
			}
		}
	}

	if len(solutions) > 0 {
		for _, solution := range solutions {
			fmt.Printf("Possible solution found at position [%d, %d] with A = %d, B = %d, C = %d\n", solution.X, solution.Y, solution.A, solution.B, solution.C)
		}
		fmt.Printf("Possible solutions count: %d\n", len(solutions))
	} else {
		fmt.Println("No possible solutions found.")
	}
}
