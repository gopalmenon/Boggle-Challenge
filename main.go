package main

import (
	"flag"
	"fmt"
	"math/rand"
)

const (
	boardSide      = 4
	dieFaces       = 6
	dictionaryFile = "http://coursera.cs.princeton.edu/algs4/testing/boggle/dictionary-yawl.txt"
)

//Dice in the Boggle game. They are all in upper case and dictionary will also be in upper case. The letter U will be
//implicit after Q
var dice = []string{"AAEEGN", "AOOTTW", "DISTTY", "EIOSST", "ABBJOO", "CIMOTU", "EEGHNW", "ELRTTY", "ACHOPS", "DEILRX", "EEINSU", "HIMNQU", "AFFKPS", "DELRVY", "EHRTVW", "HLNNRZ"}

//Command line flags for simulated annealing parameters
var acceptWorse bool
var perturbationCount, iterations int
var initialTemp, coolingRate float64

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

//Print the boggle board to standard output
func printBoard(board [][]boardDie) {

	//For every row on the board
	for _, row := range board {
		//For every die in the row
		for _, die := range row {
			fmt.Printf("%c ", dice[die.dieNumber][die.dieFace])
		}
		fmt.Printf("\n")
	}

}

//Parse command line flags for whether to accept boards with worse scores for simulated annealing and if so,
//the number of perturbations, initial temperature and temperature cooling factor
func init() {

	flag.BoolVar(&acceptWorse, "a", true, "Accept perturbed board with worse score")
	flag.IntVar(&perturbationCount, "p", 1, "Number of dice to perturb to get next board")
	flag.Float64Var(&initialTemp, "t", 1000, "Initial temperature")
	flag.Float64Var(&coolingRate, "c", 0.9, "Cooling rate")
	flag.IntVar(&iterations, "i", 1000, "Number of iterations")

}

//Load dictionary into prefix tree
func loadDictionary() {

}

func main() {

	board := boggleBoard()
	printBoard(board)

}
