package main

import (
	"fmt"
	"math/rand"
)

const (
	boardSide = 4
	dieFaces  = 6
)

//Dice in the Boggle game
var dice = []string{"aaeegn", "aoottw", "distty", "eiosst", "abbjoo", "cimotu", "eeghnw", "elrtty", "achops", "deilrx", "eeinsu", "himnqu", "affkps", "delrvy", "ehrtvw", "hlnnrz"}

//Die number and face number of a die on the board
type boardDie struct {
	dieNumber int
	dieFace   int
}

//Generate a random boggle board and return it
func boggleBoard() [][]boardDie {

	dieNumbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	board := make([][]boardDie, boardSide)

	dieCounter := 0
	rowCounter := -1

	//While there are dice that have not been placed on board
	for len(dieNumbers) > 0 {

		//Check for new row on board
		if dieCounter%boardSide == 0 {
			rowCounter += 1
		}

		//Choose from remaining dice
		dieNumber := rand.Intn(len(dieNumbers))
		die := dieNumbers[dieNumber]

		//Choose a die face
		dieFace := rand.Intn(dieFaces)

		//Place die on board in next position
		board[rowCounter] = append(board[rowCounter], boardDie{dieNumber: die, dieFace: dieFace})

		//Remove chosen die from list of dice
		dieNumbers = append(dieNumbers[:dieNumber], dieNumbers[dieNumber+1:]...)

		dieCounter += 1

	}

	return board
}

func main() {

	board := boggleBoard()

	fmt.Println(board)
}
