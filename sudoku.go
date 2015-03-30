package main

import (
	"fmt"
)

type Cell struct {
	Possible []int
}

type Game struct {
	Cells [9][9]Cell
}

const dataStr = "1..73...." +
	"..42....7" +
	"8...5...9" +
	".5...8..." +
	"..7...38." +
	"3...9.4.." +
	".61.....2" +
	"...5....." +
	"53....6.."

func parse(input string) (game Game) {
	var domain []int
	for i, d := range input {
		row := i / 9
		col := i % 9
		if d != '.' {
			game.Cells[row][col].Possible = []int{int(d - '0')}
		} else {
			domain = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
			game.Cells[row][col].Possible = domain
		}
	}
	return
}

func printGame(game Game) {
	for i, row := range game.Cells {
		for j, cell := range row {
			if len(cell.Possible) == 1 {
				fmt.Print(cell.Possible[0])
			} else {
				fmt.Print(".")
			}

			if j < 7 && j%3 == 2 {
				fmt.Print("|")
			}
		}
		fmt.Println("")
		if i < 7 && i%3 == 2 {
			fmt.Println("---|---|---")
		}
	}
}

func main() {
	game := parse(dataStr)
	//	fmt.Println("Solved %s", dataStr)
	printGame(game)
}
