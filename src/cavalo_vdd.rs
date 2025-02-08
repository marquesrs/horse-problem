use std::fmt::Debug;
use std::io::{self, Write};
use std::str::FromStr;
use std::fmt;
use std::time::{Duration, Instant};
use std::sync::{Arc, Mutex};
use std::thread;

const BOARD_SIZE: usize = 8;

static mut ITERATION_COUNTER: usize = 0;

enum Heuristic {
    None,
}

#[derive(Clone)]
struct Position {
    x: usize,
    y: usize,
}

#[derive(Clone)]
struct Board {
    cells: [usize; BOARD_SIZE*BOARD_SIZE],
    horse: Position,
    start_position: Position,
    total_moves: usize,
}

impl Board {
    fn new(x: usize, y: usize) -> Self {
        return Self {
            cells: [0; BOARD_SIZE * BOARD_SIZE],
            horse: Position{x, y},
            start_position: Position{x, y},
            total_moves: 0,
        }
    }

    fn is_solved(&self) -> bool {
        return self.total_moves == BOARD_SIZE * BOARD_SIZE - 1;   
    }    
    
    fn place_horse(&mut self, x: usize, y: usize) -> Position {
        let pos = Position {x, y};
        if !self.valid_position(&pos) || self.visited_position(&pos) {
            eprintln!("Invalid horse position x{} y{}", x, y);
        }
        self.total_moves += 1;
        self.cells[self.horse.x + self.horse.y * BOARD_SIZE] = self.total_moves;
        self.horse = pos.clone();
        self.cells[self.horse.x + self.horse.y * BOARD_SIZE] = BOARD_SIZE * BOARD_SIZE;  // Mark horse position with a unique value
        
        return pos;
    }
    
    fn valid_position(&self, pos: &Position) -> bool {
        return (pos.x < BOARD_SIZE) && (pos.y < BOARD_SIZE);
    }

    fn visited_position(&self, pos: &Position) -> bool {
        return self.cells[pos.x + pos.y * BOARD_SIZE] != 0;
    }
    
    // Possible horse moves with (x, y) as origin
    // . X . x .
    // x . . . x
    // . . H . .
    // x . . . x
    // . x . x .
    fn possible_moves(&self, x: usize, y: usize) -> Vec<Position> {
        let position = [
            Position{x: x.wrapping_sub(1), y: y + 2},
            Position{x: x + 1, y: y + 2},
            Position{x: x + 1, y: y.wrapping_sub(2)},
            Position{x: x.wrapping_sub(1), y: y.wrapping_sub(2)},
            Position{x: x.wrapping_sub(2), y: y + 1},
            Position{x: x + 2, y: y + 1},
            Position{x: x + 2, y: y.wrapping_sub(1)},
            Position{x: x.wrapping_sub(2), y: y.wrapping_sub(1)},
        ];
    
        let mut valid_positions: Vec<Position> = Vec::new();
    
        for pos in position.iter() {
            let ok_pos = self.valid_position(&pos) && !self.visited_position(&pos);
            if ok_pos {
                valid_positions.push(pos.clone());
            }
        }
    
        valid_positions
    }
    
    fn is_closed(&self) -> bool {
        let mut last_pos = Position { x: 0, y: 0 };
        let mut start_pos = Position { x: 0, y: 0 };
        
        for y in 0..BOARD_SIZE {
            for x in 0..BOARD_SIZE {
                let cell = self.cells[x + y * BOARD_SIZE];
                if cell == BOARD_SIZE * BOARD_SIZE {
                    last_pos.x = x;
                    last_pos.y = y;
                } else if cell == 1 {
                    start_pos.x = x;
                    start_pos.y = y;
                }
            }    
        }

        let tmp = Board::new(last_pos.x, last_pos.y);
        let possible = tmp.possible_moves(last_pos.x, last_pos.y);
        for p in possible {
            if p.x == start_pos.x && p.y == start_pos.y {
                return true;
            }
        }
        false
    }
}

fn display_board(b: Board) {
    for y in 0..BOARD_SIZE {
        for x in 0..BOARD_SIZE {
            if x > 0 {
                print!(" ");
            }
            let cell = b.cells[x + y * BOARD_SIZE];
            if cell == BOARD_SIZE * BOARD_SIZE {
                print!("  H");
            }
            else if cell == 0 {
                print!("  .");
            }
            else {
                print!("{:>3}", cell);
            }
        } 
        println!();
    }
}


fn increment_counter() {
    unsafe {
        ITERATION_COUNTER += 1;
    }
}

fn decrement_counter() {
    unsafe {
        ITERATION_COUNTER = 0;
    }
}

fn brute_force_solve(base_x: usize, base_y: usize, h: Heuristic, visualize: bool) -> (Board, bool) {
    match h {
        Heuristic::None => {
            if let (s, t) = brute_force_rec(Board::new(base_x, base_y), 0) {
                return (s, t);
            }
        }
    }
    
    return (Board::new(base_x, base_y), false)
}

fn brute_force_rec(mut b: Board, level: i32) -> (Board, bool) {
    increment_counter();
    
    if b.is_solved() {
        return (b, true);
    }
    
    let possible = b.possible_moves(b.horse.x, b.horse.y);
    
    for pos in possible {
        b.place_horse(pos.x, pos.y);
        
        let (board, ok) = brute_force_rec(b.clone(), level + 1);
        if ok {
            return (board, true);    
        }
    }
    
    return (b, false);
}


fn input<T: FromStr + fmt::Debug>(msg: &str) -> T 
    where <T as FromStr>::Err: Debug
{
    print!("{}", msg);
    io::stdout().flush().unwrap();
    
    let mut input = String::new();
    io::stdin().read_line(&mut input).unwrap();

    return input.trim().parse::<T>().unwrap();
}


fn main() {
    let x: usize = input::<usize>("x: ");
    let y: usize = input::<usize>("y: ");    
    
    println!("Begin solve for x: {} y: {}", x, y);
    
    let running = Arc::new(Mutex::new(true));
    
    let running_clone = Arc::clone(&running);
    
    thread::spawn(move || {
        let start = Instant::now();
        let (b, solved) = brute_force_solve(x, y, Heuristic::None, true);
        let elapsed = start.elapsed();
        
        let status: &str;
        
        if solved {
            status = "[ Solved ]";
        } else {
            status = "[ Unsolved ]";
        }
        
        {
            let mut running = running_clone.lock().unwrap();
            *running = false;
        }
        
        println!(
            "{} Took: {:?} Iterations: {} Closed? {}",
            status, 
            elapsed, 
            unsafe {ITERATION_COUNTER}, 
            b.is_closed()
        );
        
        display_board(b);
        decrement_counter();
    });
    
    let start_point = Instant::now();
    while *running.lock().unwrap() {
        let iteration_count = unsafe { ITERATION_COUNTER };  
        let it_per_sec = iteration_count as f64 / start_point.elapsed().as_secs_f64();  
        print!(
            "\r                                                                 \rIter:{} Iter/s: {:.0}",
            iteration_count,
            it_per_sec
        );
        thread::sleep(Duration::from_millis(5)); 
    }
}

