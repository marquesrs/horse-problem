package main

import (
	"fmt"
	"log"
	"math"
	"slices"
	"time"
)

const MAXIMUM_ITER_COUNT = (1 << 32) - 1

const BOARD_SIZE = 10

func RemoveUnordered[T any](s []T, idx int) []T {
	s[len(s)-1], s[idx] = s[idx], s[len(s)-1]
	return s[:len(s)-1]
}

func abs[T int | float32 | float64](x T) T {
	if x < T(0) {
		return -x
	}
	return x
}

type BoardCallback func(b *Board)

type Board struct {
	Cells         [BOARD_SIZE * BOARD_SIZE]int
	Horse         Position
	StartPosition Position
	TotalMoves    int
	PlaceCallback  BoardCallback
}

func boardDefaultCallback(b *Board){}

func boardVisualizeCallback(b *Board){
	fmt.Print("\033[H\033[2J")
	DisplayBoard(*b)
	time.Sleep(16 * time.Millisecond)
}

func NewBoard(x, y int) Board {
	return Board{
		Horse:         Position{x, y},
		StartPosition: Position{x, y},
		TotalMoves:    0,
		PlaceCallback: boardDefaultCallback,
	}
}

type Position struct {
	x, y int
}

func (b Board) IsClosed() bool {
	lastPos := Position{}
	startPos := Position{}

	for y := 0; y < BOARD_SIZE; y++ {
		for x := 0; x < BOARD_SIZE; x++ {
			cell := b.Cells[x+y*BOARD_SIZE]
			if cell == -1 {
				lastPos.x, lastPos.y = x, y
			} else if cell == 1 {
				startPos.x, startPos.y = x, y
			}
		}
	}

	tmp := NewBoard(lastPos.x, lastPos.y)
	possible := tmp.PossibleMoves(lastPos.x, lastPos.y)
	for _, p := range possible {
		if p == startPos {
			return true
		}
	}
	return false
}

func (b Board) ValidPosition(pos Position) bool {
	return (pos.x >= 0 && pos.x < BOARD_SIZE) && (pos.y >= 0 && pos.y < BOARD_SIZE)
}

func (b Board) VisitedPostion(pos Position) bool {
	return b.Cells[pos.x+pos.y*BOARD_SIZE] != 0
}

// Possible horse moves with (x, y) as origin
// . X . x .
// x . . . x
// . . H . .
// x . . . x
// . x . x .
func (b Board) PossibleMoves(x, y int) []Position {
	positions := []Position{
		{x - 1, y + 2},
		{x + 1, y + 2},

		{x + 1, y - 2},
		{x - 1, y - 2},

		{x - 2, y + 1},
		{x + 2, y + 1},

		{x + 2, y - 1},
		{x - 2, y - 1},
	}

	validPositions := make([]Position, 0, len(positions))

	for _, pos := range positions {
		okPos := b.ValidPosition(pos) && (!b.VisitedPostion(pos))

		if okPos {
			validPositions = append(validPositions, pos)
		}
	}

	// fmt.Println("POSITIONS", positions);
	// fmt.Println("VALID", validPositions);
	return validPositions
}

// Place the horse
func (b *Board) PlaceHorse(x, y int) Position {
	pos := Position{x, y}
	if !b.ValidPosition(pos) || b.VisitedPostion(pos) {
		log.Fatal("Invalid horse position:", x, y)
	}
	b.TotalMoves += 1
	b.Cells[b.Horse.x+b.Horse.y*BOARD_SIZE] = b.TotalMoves

	b.Horse = pos
	b.Cells[b.Horse.x+b.Horse.y*BOARD_SIZE] = -1

	b.PlaceCallback(b)
	return pos
}

func (b Board) IsSolved() bool {
	return b.TotalMoves == BOARD_SIZE*BOARD_SIZE-1
}

func DisplayBoard(b Board) {
	for y := 0; y < BOARD_SIZE; y++ {
		for x := 0; x < BOARD_SIZE; x++ {
			if x > 0 {
				fmt.Print(" ")
			}
			cell := b.Cells[x+y*BOARD_SIZE]

			if cell == -1 {
				fmt.Print("  H")
			} else if cell == 0 {
				fmt.Print("  .")
			} else {
				fmt.Printf("%3d", cell)
			}
		}
		fmt.Println()
	}
}

type Heuristic uint32

const (
	None Heuristic = iota
	PreferCorners
	PreferRunning
)

func ManhattanDist(a, b Position) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func EuclidianDist(a, b Position) float64 {
	return math.Sqrt(float64(a.x-b.x)*float64(a.x-b.x) + (float64(a.y-b.y) * float64(a.y-b.y)))
}

func BruteForceSolve(baseX, baseY int, h Heuristic, visualize bool) (Board, bool) {
	b, solved := NewBoard(baseX, baseY), false

	if visualize {
		b.PlaceCallback = boardVisualizeCallback
	}

	switch h {
	case None:
		b, solved = BruteForceRec(b, 0)
	case PreferCorners:
		b, solved = BruteForcePreferCorners(b, 0)
	case PreferRunning:
		b, solved = BruteForcePreferRunning(b, 0)
	}

	return b, solved
}


func BruteForceRec(b Board, level int) (Board, bool) {
	IterationCounter += 1
	if b.IsSolved() {
		return b, true
	}

	ok := false

	possible := b.PossibleMoves(b.Horse.x, b.Horse.y)
	for _, move := range possible {
		board := b
		board.PlaceHorse(move.x, move.y)


		if board, ok = BruteForceRec(board, level+1); ok {
			return board, true
		}
	}

	return b, false
}

// How close is a point to a corner
func CornerScore(p Position) int {
	topLeft := ManhattanDist(p, Position{0, 0})
	topRight := ManhattanDist(p, Position{BOARD_SIZE - 1, 0})
	bottomLeft := ManhattanDist(p, Position{0, BOARD_SIZE - 1})
	bottomRight := ManhattanDist(p, Position{BOARD_SIZE - 1, BOARD_SIZE - 1})
	return min(topLeft, topRight, bottomLeft, bottomRight)
}

var IterationCounter = 0

func BruteForcePreferCorners(b Board, level int) (Board, bool) {
	IterationCounter += 1
	if b.IsSolved() {
		return b, true
	}

	ok := false

	possible := b.PossibleMoves(b.Horse.x, b.Horse.y)

	slices.SortFunc(possible, func(a, b Position) int {
		return CornerScore(a) - CornerScore(b)
	})

	for _, move := range possible {
		board := b
		board.PlaceHorse(move.x, move.y)

		if board, ok = BruteForcePreferCorners(board, level+1); ok {
			return board, true
		}
	}

	return b, false
}

func BruteForcePreferRunning(b Board, level int) (Board, bool) {
	IterationCounter += 1
	if b.IsSolved() {
		return b, true
	}

	ok := false

	possible := b.PossibleMoves(b.Horse.x, b.Horse.y)

	slices.SortFunc(possible, func(p0, p1 Position) int {
		return ManhattanDist(p0, b.StartPosition) - ManhattanDist(p1, b.StartPosition)
	})

	for _, move := range possible {
		board := b
		board.PlaceHorse(move.x, move.y)

		if board, ok = BruteForcePreferRunning(board, level+1); ok {
			return board, true
		}
	}

	return b, false
}


func main() {
	// Closed (4, 4) Board: 8
	X, Y := 5, 5
	fmt.Println("Begin solve for", X, Y)
	running := true
	go func() {
		start := time.Now()
		b, solved := BruteForceSolve(X, Y, PreferCorners, true)
		elapsed := time.Since(start)

		status := "[ Solved ] "
		if !solved {
			status = "[ Unsolved ] "
		}
		running = false

		fmt.Println(status, "Took:", elapsed, "Iterations:", IterationCounter, "Closed?", b.IsClosed())
		DisplayBoard(b)
		_ = b
		IterationCounter = 0
	}()

	for running {
		// elapsed := time.Since(startPoint)
		// fmt.Printf("\r                                              \rN:%#v It/s: %v", IterationCounter, IterationCounter/(1+int(elapsed/time.Second)))
		time.Sleep(100 * time.Millisecond)
	}
}
