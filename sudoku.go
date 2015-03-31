package main

import (
	"domain"
	"fmt"
)

//////////////////////////////////////////
// Game data

type Cell struct {
	domain domain.Domain
}

type Game struct {
	Cells [9][9]*Cell // [row][col]
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
	for i, d := range input {
		row := i / 9
		col := i % 9
		var cell Cell
		game.Cells[row][col] = &cell
		game.Cells[row][col].domain = domain.New()
		if d != '.' {
			game.Cells[row][col].domain.Add(int(d - '0'))
		} else {
			for i := 1; i <= 9; i += 1 {
				game.Cells[row][col].domain.Add(i)
			}
		}
	}
	return
}

func row(game Game, i int) []*Cell {
	return game.Cells[i][0:]
}

func col(game Game, i int) []*Cell {
	cells := make([]*Cell, 9)
	for j, row := range game.Cells {
		cells[j] = row[i]
	}
	return cells
}

// 0|1|2
// 3|4|5
// 6|7|8
func block(game Game, i int) []*Cell {
	startRow := (i / 3) * 3

	startCol := (i % 3) * 3
	cells := make([]*Cell, 9)
	for row := 0; row < 3; row += 1 {
		offset := row * 3
		slice := game.Cells[startRow+row][startCol : startCol+3]
		for i, cell := range slice {
			cells[offset+i] = cell
		}
	}
	return cells
}

func printCells(cells []*Cell) {
	fmt.Println("*****")
	for _, cell := range cells {
		fmt.Println(cell.domain)
	}
	fmt.Println("*****")
}

func printGame(game Game) {
	for i, row := range game.Cells {
		for j, cell := range row {
			if cell.domain.Size() == 1 {
				val, _ := cell.domain.Any()
				fmt.Print(val)
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

func isSolved(game Game) bool {
	for i := range game.Cells {
		for j := range game.Cells[i] {
			if game.Cells[i][j].domain.Size() > 1 {
				return false
			}
		}
	}
	return true
}



//////////////////////////////////////////
// Constraint propagation (CSP)

type constraintFn func(values []int) bool

// Constraints are still coupled to the game via Cell. I guess this
// could be removed by introducing an interface that had Domain() and
// SetNewDomain() functions
type constraint struct {
	Variables  []*Cell
	constraint constraintFn
}

type domainFunc func([]int)

func applyTo(variables []*Cell, list []int, cb domainFunc) {
	if len(variables) == 0 {
		cb(list)
	} else {
		variables[0].domain.ForAll(func(val int) {
			var newList = append(list, val)
			applyTo(variables[1:], newList, cb)
		})
	}
}

func propagateConstraint(constraint constraint) {
	newDomains := make([]domain.Domain, len(constraint.Variables))
	for i := 0; i < len(constraint.Variables); i += 1 {
		newDomains[i] = domain.New()
	}

	applyTo(constraint.Variables, []int{}, func(values []int) {
		if constraint.constraint(values) {
			for i, value := range values {
				newDomains[i].Add(value)
			}
		}
	})

	for i, variable := range constraint.Variables {
		variable.domain = newDomains[i]
	}
}

func propagateConstraints(constraints []constraint) {
	for _, constraint := range constraints {
		propagateConstraint(constraint)
	}
}

//////////////////////////////////////////
// Game<->CSP adapter

func hasOneToNine(values []int) bool {
	var set [10]int // coming in as 1..9, not 0..8
	for _, val := range values {
		set[val] = val
	}
	for i := 1; i <= 9; i += 1 {
		if set[i] == 0 {
			return false
		}
	}
	return true
}

func notEqual(values []int) bool {
	return values[0] != values[1]
}

func blockConstraints(game Game) []constraint {
	var constraints = make([]constraint, 3*len(game.Cells))
	for i := 0; i < len(game.Cells); i += 1 {
		var rows = row(game, i)
		var cols = col(game, i)
		var blocks = block(game, i)
		constraints[3*i].Variables = append([]*Cell(nil), blocks[:]...)
		constraints[3*i].constraint = hasOneToNine
		constraints[3*i+1].Variables = append([]*Cell(nil), rows[:]...)
		constraints[3*i+1].constraint = hasOneToNine
		constraints[3*i+2].Variables = append([]*Cell(nil), cols[:]...)
		constraints[3*i+2].constraint = hasOneToNine
	}
	return constraints
}

func arcConstraintsOnCells(cells []*Cell) []constraint {
	var constraints []constraint
	for this := range cells {
		for that := range cells {
			if (this != that) {
				var constraint constraint
				constraint.Variables = []*Cell{cells[this], cells[that]}
				constraint.constraint = notEqual
				constraints = append(constraints, constraint)
			}
		}
	}
	return constraints
}

func arcConstraints(game Game) []constraint {
	var constraints []constraint
	for i := 0; i < len(game.Cells); i += 1 {
		constraints = append(constraints, arcConstraintsOnCells(row(game, i))...)
		constraints = append(constraints, arcConstraintsOnCells(col(game, i))...)
		constraints = append(constraints, arcConstraintsOnCells(block(game, i))...)
	}
	return constraints
}

func main() {
	game := parse(dataStr)
	printGame(game)

	constraints := append(arcConstraints(game), blockConstraints(game)...)

	fmt.Println("\nSolving...\n")

	for !isSolved(game) {
		propagateConstraints(constraints)
	}

	printGame(game)
}
