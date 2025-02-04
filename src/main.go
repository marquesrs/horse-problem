package main

import (
	"fmt"
	"log"
)

const BOARD_SIZE = 4

func RemoveUnordered[T any](s []T, idx int) []T {
	s[len(s) - 1], s[idx] = s[idx], s[len(s) - 1]
	return s[:len(s) - 1]
}

type Board struct {
    Cells [BOARD_SIZE * BOARD_SIZE]int
	LastIndex int
    Horse Position
}

type Position struct {
    x, y int
    score int
}

func (b Board) ValidPosition(pos Position) bool {
	return pos.x >= 0 && pos.x < BOARD_SIZE && pos.y >= 0 && pos.y < BOARD_SIZE
}

func (b Board) VisitedPostion(pos Position) bool {
	return b.Cells[pos.x + pos.y * BOARD_SIZE] > 0
}

// Possible horse moves with (x, y) as origin
// . X . x .
// x . . . x
// . . H . .
// x . . . x
// . x . x .
func (b Board) PossibleMoves(x, y int) []Position{
	positions := []Position{
		{x - 1, y + 2, 0},
		{x + 1, y + 2, 0},

		{x - 2, y + 1, 0},
		{x + 2, y + 1, 0},

		{x - 2, y - 1, 0},
		{x + 2, y - 1, 0},

		{x - 1, y - 2, 0},
		{x + 1, y - 2, 0},
	}

	for i := len(positions) - 1; i >= 0; i-- {
		okPos := b.ValidPosition(positions[i]) && !b.VisitedPostion(positions[i])
		if !okPos {
			positions = RemoveUnordered(positions, i)
		}
	}

    return positions
}

// Give a score for each position given current board states
func (b Board) GradePositions(positions []Position) []Position {
    return nil
}

// Place the horse
func (b *Board) PlaceHorse(x, y int) Position {
	pos := Position{x, y, 0}
	if !b.ValidPosition(pos) {
		log.Fatal("Invalid horse position:", x, y)
	}
	b.Cells[b.Horse.x + b.Horse.y * BOARD_SIZE] = b.LastIndex + 1
	//Miguel: Ao Corno sendo alérgico a matriz, como caralhos essa parte de cima faz sentido?
	b.LastIndex += 1
	b.Horse = pos
    return pos
}

func DisplayBoard(b Board){
	for y := 0; y < BOARD_SIZE; y++ {
		for x := 0; x < BOARD_SIZE; x++ {
			if x > 0 { fmt.Print(" ") }
			cell := b.Cells[x + y * BOARD_SIZE]

			if x == b.Horse.x && y == b.Horse.y {
				fmt.Print(" H")
			} else if cell == 0 {
				fmt.Print(" .")
			} else {
				fmt.Printf("%2d", cell)
			}
		}
		fmt.Println()
	}
}

//Replace the horse one movement (input the place you want the horse to be)
func (b *Board) ReturnHorse(x, y int){
	pos := Position{x, y, 0}
	if !b.ValidPosition(pos) {
		log.Fatal("Invalid horse position:", x, y)
	}
	b.Cells[b.Horse.x + b.Horse.y * BOARD_SIZE] = 0
	b.LastIndex -= 1
	b.Horse = pos
}

//Miguel: The Dumb Way (with recursion) still wrong
func (b *Board) RecursionWay(x, y int){
	lastPos := Position{x, y, 0}
	possibleMoves := b.PossibleMoves(x, y)
	if b.LastIndex == (BOARD_SIZE * BOARD_SIZE)-1{
		return 
	}
	for i := 0; i < len(possibleMoves); i++{
		if b.VisitedPostion(possibleMoves[i]) || !b.ValidPosition(possibleMoves[i]){//Miguel: só pra garantir, não sei se precisa mesmo
			fmt.Println("You Fucked Up BIG TIME")
			return
		}
		b.PlaceHorse(possibleMoves[i].x, possibleMoves[i].y)
		b.RecursionWay(possibleMoves[i].x, possibleMoves[i].y)
		if b.LastIndex == (BOARD_SIZE * BOARD_SIZE)-1{
			return 
		}
		b.ReturnHorse(lastPos.x, lastPos.y)
	}
	return
}

func main() {
	fmt.Println("pipi")
	b := Board{}
	DisplayBoard(b)
	fmt.Println("------------")
	b.RecursionWay(0, 0)
	DisplayBoard(b)
}
