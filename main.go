package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/derekparker/trie"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

const (
	boardSide      = 4
	dieFaces       = 6
	dictionaryFile = "http://coursera.cs.princeton.edu/algs4/testing/boggle/dictionary-yawl.txt"
)

//Dice in the Boggle game. They are all in upper case and dictionary will also be in upper case. The letter U will be
//implicit after Q
var dice = []string{"AAEEGN", "AOOTTW", "DISTTY", "EIOSST", "ABBJOO", "CIMOTU", "EEGHNW", "ELRTTY", "ACHOPS", "DEILRX", "EEINSU", "HIMNQU", "AFFKPS", "DELRVY", "EHRTVW", "HLNNRZ"}
var dieNumbers = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

//Scores corresponding to word length
var scores = []int{0, 0, 0, 1, 1, 2, 3, 5, 11, 11, 11, 11, 11, 11, 11, 11, 11}

//Command line flags for simulated annealing parameters
var acceptWorse, logProgress bool
var perturbationCount, iterations int
var initialTemp, coolingRate float64

//Die number and face number of a die on the board
type boardDie struct {
	dieNumber int
	dieFace   int
}

//Row and column coordinates for a die
type dieCoordinates struct {
	row    int
	column int
}

type bestBoard struct {
	board string
	score int
	words string
}

var bestYet = bestBoard{board: "", score: -1, words: ""}

var runLog strings.Builder

//Find face showing on the die
func dieFace(dieNumber int, board [][]boardDie) int {

	for _, row := range board {
		for _, die := range row {
			if die.dieNumber == dieNumber {
				return die.dieFace
			}
		}
	}

	fmt.Printf("%v\n", board)
	log.Fatal("Could not find die ", dieNumber)
	return 0
}

//Perturb a die
func perturbDie(board [][]boardDie, dieNumber, newDieFace int) {

	for _, row := range board {
		for dieIndex := range row {
			if row[dieIndex].dieNumber == dieNumber {
				row[dieIndex].dieFace = newDieFace
				return
			}
		}
	}
	fmt.Printf("%v\n", board)
	log.Fatal("Could not find die ", dieNumber)
}

//Perturb board in the hope of getting one with a better score
func perturbBoard(board [][]boardDie) [][]boardDie {

	dieFaces := []int{0, 1, 2, 3, 4, 5}

	//Make a copy of the board
	boardCopy := make([][]boardDie, len(board))
	copy(boardCopy, board)

	remainingDice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	for counter := 0; counter < perturbationCount; counter++ {

		//Pick a die to perturb
		dieIndex := rand.Intn(len(remainingDice))
		dieNumber := remainingDice[dieIndex]
		dieFace := dieFace(dieNumber, boardCopy)

		//Find new die face
		remainingDieFaces := append(dieFaces[:dieFace], dieFaces[dieFace+1:]...)
		newDieFace := remainingDieFaces[rand.Intn(len(remainingDieFaces))]

		//Perturb the die
		perturbDie(boardCopy, dieNumber, newDieFace)

		//Remove perturbed die from future perturbations
		remainingDice = append(remainingDice[:dieIndex], remainingDice[dieIndex+1:]...)

	}

	return boardCopy
}

//Generate a random boggle board and return it
func boggleBoard() [][]boardDie {

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

//Get neighbors of die that have already not been used
func neighbors(row, column int, alreadyUsed map[dieCoordinates]bool) []dieCoordinates {

	dieNeighbors := make([]dieCoordinates, 0, 8)

	//Add neighbors above
	if row != 0 {
		if column != 0 {
			candidate := dieCoordinates{row: row - 1, column: column - 1}
			if _, present := alreadyUsed[candidate]; !present {
				dieNeighbors = append(dieNeighbors, candidate)
			}
		}
		candidate := dieCoordinates{row: row - 1, column: column}
		if _, present := alreadyUsed[candidate]; !present {
			dieNeighbors = append(dieNeighbors, candidate)
		}
		if column < boardSide-1 {
			candidate := dieCoordinates{row: row - 1, column: column + 1}
			if _, present := alreadyUsed[candidate]; !present {
				dieNeighbors = append(dieNeighbors, candidate)
			}
		}
	}

	//Add neighbors below
	if row < boardSide-1 {
		if column != 0 {
			candidate := dieCoordinates{row: row + 1, column: column - 1}
			if _, present := alreadyUsed[candidate]; !present {
				dieNeighbors = append(dieNeighbors, candidate)
			}
		}
		candidate := dieCoordinates{row: row + 1, column: column}
		if _, present := alreadyUsed[candidate]; !present {
			dieNeighbors = append(dieNeighbors, candidate)
		}
		if column < boardSide-1 {
			candidate := dieCoordinates{row: row + 1, column: column + 1}
			if _, present := alreadyUsed[candidate]; !present {
				dieNeighbors = append(dieNeighbors, candidate)
			}
		}
	}

	//Add neighbors at the sides
	if column != 0 {
		candidate := dieCoordinates{row: row, column: column - 1}
		if _, present := alreadyUsed[candidate]; !present {
			dieNeighbors = append(dieNeighbors, candidate)
		}
	}
	if column < boardSide-1 {
		candidate := dieCoordinates{row: row, column: column + 1}
		if _, present := alreadyUsed[candidate]; !present {
			dieNeighbors = append(dieNeighbors, candidate)
		}
	}

	return dieNeighbors
}

//Add to already used dice map
func add(alreadyUsed map[dieCoordinates]bool, neighbor dieCoordinates) map[dieCoordinates]bool {

	newMap := make(map[dieCoordinates]bool, len(alreadyUsed)+1)
	for k, v := range alreadyUsed {
		newMap[k] = v
	}

	newMap[neighbor] = true

	return newMap
}

//Compute score by adding neighbor
func scoreWithNeighbor(prefix string, neighbor dieCoordinates, alreadyUsed map[dieCoordinates]bool,
	dict *trie.Trie, board [][]boardDie, alreadyScored map[string]bool) int {

	score := 0
	neighborDie := board[neighbor.row][neighbor.column]

	var str strings.Builder
	str.WriteString(prefix)
	str.WriteString(string(dice[neighborDie.dieNumber][neighborDie.dieFace]))
	newPrefix := str.String()

	//Check if any words start with the new prefix
	if !dict.HasKeysWithPrefix(newPrefix) {
		return 0
	}

	//Add score if new prefix is a valid word and has not already been scored
	if _, found := dict.Find(newPrefix); found {
		if _, scored := alreadyScored[newPrefix]; !scored {
			score += scores[len(newPrefix)]
			alreadyScored[newPrefix] = true
		}
	}

	updatedAlreadyUsed := add(alreadyUsed, neighbor)

	//Check for words starting with new prefix
	dieNeighbors := neighbors(neighbor.row, neighbor.column, updatedAlreadyUsed)
	for _, neighbor := range dieNeighbors {
		score += scoreWithNeighbor(newPrefix, neighbor, updatedAlreadyUsed, dict, board, alreadyScored)
	}

	return score

}

//Compute score for the board
func score(board [][]boardDie, dict *trie.Trie, alreadyScored map[string]bool) int {

	boardScore := 0

	//Loop through each position on the board
	for rowNumber, row := range board {
		for columnNumber, die := range row {

			alreadyUsed := make(map[dieCoordinates]bool)
			alreadyUsed[dieCoordinates{row: rowNumber, column: columnNumber}] = true
			dieNeighbors := neighbors(rowNumber, columnNumber, alreadyUsed)

			//Check if any words start with current die
			dieFace := dice[die.dieNumber][die.dieFace]
			if !dict.HasKeysWithPrefix(string(dieFace)) {
				continue
			}

			//Check for words starting with die in current position
			for _, neighbor := range dieNeighbors {
				boardScore += scoreWithNeighbor(string(dieFace), neighbor, alreadyUsed, dict, board, alreadyScored)
			}
		}
	}

	return boardScore

}

//Parse command line flags for whether to accept boards with worse scores for simulated annealing and if so,
//the number of perturbations, initial temperature and temperature cooling factor
func init() {

	flag.BoolVar(&acceptWorse, "a", true, "Accept perturbed board with worse score")
	flag.BoolVar(&logProgress, "l", false, "Log progress of best board search process")
	flag.IntVar(&perturbationCount, "p", 1, "Number of dice to perturb to get next board")
	flag.Float64Var(&initialTemp, "t", 1000, "Initial temperature")
	flag.Float64Var(&coolingRate, "c", 0.99, "Cooling rate")
	flag.IntVar(&iterations, "i", 10000, "Number of iterations of simulated annealing")
	flag.Parse()

}

//Load list of valid words into a prefix tree
func loadIntoPrefixTree() *trie.Trie {

	//Read dictionary text file
	resp, err := http.Get(dictionaryFile)
	if err != nil {
		log.Fatal(err)
	}

	//Close response body after use
	defer resp.Body.Close()

	//Load valid words in dictionary into prefix tree
	t := trie.New()
	s := bufio.NewScanner(resp.Body)
	for tok := s.Scan(); tok != false; tok = s.Scan() {
		words := strings.Fields(s.Text())
		for _, word := range words {
			t.Add(word, 1)
		}
	}

	return t

}

//Print best board along with associated score and words in the board
func showBestBoard() {

	fmt.Println("\nBest board:")
	fmt.Println(bestYet.board)
	fmt.Printf("Best board score: %d\n", bestYet.score)
	fmt.Println("Best board words:")
	fmt.Println(bestYet.words)

}

//Save board as best yet board
func saveBestYet(board [][]boardDie, boardScore int, boardWords map[string]bool) {

	var boardStr strings.Builder

	//For every row on the board
	for _, row := range board {
		//For every die in the row
		for _, die := range row {
			boardStr.WriteString(string(dice[die.dieNumber][die.dieFace]))
			boardStr.WriteString(" ")
		}
		boardStr.WriteString("\n")
	}
	bestYet.board = boardStr.String()

	bestYet.score = boardScore

	var wordsStr strings.Builder

	for k := range boardWords {
		wordsStr.WriteString(k)
		wordsStr.WriteString(" ")
	}

	bestYet.words = wordsStr.String()

}

//Return a copy of the board
func copyBoard(board [][]boardDie) [][]boardDie {

	boardCopy := make([][]boardDie, len(board))
	copy(boardCopy, board)
	return boardCopy

}

func main() {

	//Generate random board
	prevBoard := boggleBoard()

	//Valid words from online list
	validWords := loadIntoPrefixTree()

	//Compute score for the board
	boardWords := make(map[string]bool)
	boardScore := score(prevBoard, validWords, boardWords)

	//Save as best yet board
	saveBestYet(prevBoard, boardScore, boardWords)

	boardToPerturb := copyBoard(prevBoard)

	prevBoardScore := boardScore
	currentTemperature := initialTemp

	//Run for predetermined number of iterations
	for counter := 0; counter < iterations; counter++ {

		//Generate a new board by perturbing the previous one
		perturbedBoard := perturbBoard(boardToPerturb)

		perturbedBoardWords := make(map[string]bool)
		newBoardScore := score(perturbedBoard, validWords, perturbedBoardWords)

		//Save if this is the best yet
		if newBoardScore > bestYet.score {
			saveBestYet(perturbedBoard, newBoardScore, perturbedBoardWords)
			if logProgress {
				runLog.WriteString("\nSaved best yet board at iteration ")
				runLog.WriteString(strconv.Itoa(counter))
			}
		}

		//Keep the perturbed board if its better than the previous one
		if newBoardScore > prevBoardScore {
			boardToPerturb = copyBoard(perturbedBoard)
			prevBoard = copyBoard(perturbedBoard)
			prevBoardScore = newBoardScore
			if logProgress {
				runLog.WriteString("\nNew board better than previous board at iteration ")
				runLog.WriteString(strconv.Itoa(counter))
			}
		} else {
			//Keep the worse board depending on probability
			if acceptWorse && rand.Float64() < math.Exp(-1*float64(prevBoardScore-boardScore)/currentTemperature) {
				boardToPerturb = copyBoard(perturbedBoard)
				prevBoard = copyBoard(perturbedBoard)
				prevBoardScore = newBoardScore
				if logProgress {
					runLog.WriteString("\nAccepted worse board at iteration ")
					runLog.WriteString(strconv.Itoa(counter))
				}
			} else {
				//Reject the perturbed board
				boardToPerturb = copyBoard(prevBoard)
				if logProgress {
					runLog.WriteString("\nRejected worse board at iteration ")
					runLog.WriteString(strconv.Itoa(counter))
				}
			}
		}

		currentTemperature *= coolingRate

	}

	if logProgress {
		fmt.Println(runLog.String())
	}
	showBestBoard()

}
