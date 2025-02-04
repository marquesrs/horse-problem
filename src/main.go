package main

import (
	"fmt"
	"log"
)

const BOARD_SIZE = 8

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

func main() {
	b := Board{}
	DisplayBoard(b)
	fmt.Println("------------")
	// NOTE: Obviously incorrect moves, just testing
	b.PlaceHorse(4, 4)
	b.PlaceHorse(1, 2)
	b.PlaceHorse(2, 7)
	b.PlaceHorse(3, 1)
	DisplayBoard(b)
	fmt.Println(b.PossibleMoves(4, 4))
}
