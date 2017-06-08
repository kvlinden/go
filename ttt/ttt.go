package main

import "fmt"

type TicTacToeBoard [][]int

func NewTicTacToeBoard() TicTacToeBoard {
	result := TicTacToeBoard(make([][]int, 3))
	for i := 0; i < 3; i++ {
		result[i] = make([]int, 3)
	}
	return result
}

func (tttBoard TicTacToeBoard) PPrint() {
	pieces := []string{"X", "-", "O"}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Print(pieces[tttBoard[i][j]+1])
		}
		fmt.Println()
	}
}

func main() {
	tttBoard := NewTicTacToeBoard()
	tttBoard[0][0] = -1
	tttBoard[1][1] = 1
	tttBoard.PPrint()
}
