package main

import "fmt"

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
	solutions := make([]int, 0)

	// Limit iteration to startRow, more than that and going north will go out of puzzle boundary
	for y := 1; y <= startRow; y++ {
		fmt.Printf("Trying Y = %d -- ", y)
		isPathOpen := true

		// Going north
		north := 0
		for north < y {
			north++
			isPathOpen = puzzle[startRow-north][startCol]
			if !isPathOpen {
				fmt.Printf("Path blocked at row %d col %d when going north.\n", startRow-north, startCol)
				break
			}
		}
		if !isPathOpen {
			continue
		}

		// Going east
		east := 0
		for east < y {
			east++
			isPathOpen = puzzle[startRow-north][startCol+east]
			if !isPathOpen {
				fmt.Printf("Path blocked at row %d col %d when going east.\n", startRow-north, startCol+east)
				break
			}
		}
		if !isPathOpen {
			continue
		}

		// Going south
		south := 0
		for south < y {
			south++
			isPathOpen = puzzle[startRow-north+south][startCol+east]
			if !isPathOpen {
				fmt.Printf("Path blocked at row %d col %d when going south.\n", startRow-north+south, startCol+east)
				break
			}
		}

		if isPathOpen {
			solutions = append(solutions, y)
			fmt.Println("Possible solution found.")
		}
	}

	fmt.Println("\nSUMMARY:")
	if len(solutions) > 0 {
		for _, y := range solutions {
			fmt.Printf("Possible solution found at Y = %d.\n", y)
		}
	} else {
		fmt.Println("No possible solutions found.")
	}
}
