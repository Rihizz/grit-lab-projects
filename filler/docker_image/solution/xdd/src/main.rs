use std::io;
use std::io::Write;
use std::io::BufRead;

fn main() {
    let stdin = io::stdin();
    let mut anfield: Vec<Vec<char>> = Vec::new();
    let mut piece: Vec<Vec<char>> = Vec::new();
    let mut p: String = String::new();

    let mut lines = stdin.lock().lines();
    loop {
        while let Some(result) = lines.next() {
            let line = result.expect("Failed to read line");

            if line.starts_with("$$$") {
                let p_parts: Vec<&str> = line.split_whitespace().collect();
                let player: String = p_parts[2].parse().expect("Failed to parse player");
                p = player;
            }

            if line.starts_with("Anfield") {
                let an_parts: Vec<&str> = line.split_whitespace().collect();
                let an_height: usize = an_parts[2]
                    .trim_end_matches(':')
                    .parse()
                    .expect("Failed to parse height");

                // Skip the next line which only has the numbers
                let _ = lines.next().expect("Failed to skip numbers line");

                for _ in 0..an_height {
                    if let Some(board_line_result) = lines.next() {
                        let board_line = board_line_result.expect("Failed to read board line");
                        let chars: Vec<char> = board_line[4..].chars().collect();
                        anfield.push(chars);
                    }
                }
            }

            if line.starts_with("Piece") {
                let parts: Vec<&str> = line.split_whitespace().collect();
                let height: usize = parts[2]
                    .trim_end_matches(':')
                    .parse()
                    .expect("Failed to parse height");

                for _ in 0..height {
                    if let Some(piece_line_result) = lines.next() {
                        let piece_line = piece_line_result.expect("Failed to read piece line");
                        let chars: Vec<char> = piece_line.chars().collect();
                        piece.push(chars);
                    }
                }

                let enemy: char;
                let enemy2: char;
                if p == "p1" {
                    enemy = 's';
                    enemy2 = '$';
                } else {
                    enemy = 'a';
                    enemy2 = '@';
                }

                let valid_positions = place_piece_on_board(&anfield, &piece, &p);
                let enemy_positions: Vec<_> = anfield
                    .iter()
                    .enumerate()
                    .flat_map(|(i, row)| {
                        row.iter()
                            .enumerate()
                            .filter(|&(_, &ch)| (ch == enemy || ch == enemy2))
                            .map(move |(j, _)| (i, j))
                    })
                    .collect();

                if !valid_positions.is_empty() {
                    let (chosen_x, chosen_y) = valid_positions
                        .into_iter()
                        .min_by(|&(x1, y1), &(x2, y2)| {
                            let min_distance1 = enemy_positions
                                .iter()
                                .map(|&(ex, ey)| distance(x1, y1, ex, ey))
                                .min_by(|&a, &b|
                                    a.partial_cmp(&b).unwrap_or(std::cmp::Ordering::Equal)
                                )
                                .unwrap_or(f64::INFINITY);
                            let min_distance2 = enemy_positions
                                .iter()
                                .map(|&(ex, ey)| distance(x2, y2, ex, ey))
                                .min_by(|&a, &b|
                                    a.partial_cmp(&b).unwrap_or(std::cmp::Ordering::Equal)
                                )
                                .unwrap_or(f64::INFINITY);
                            min_distance1
                                .partial_cmp(&min_distance2)
                                .unwrap_or(std::cmp::Ordering::Equal)
                        })
                        .unwrap();

                    println!("{} {}", chosen_y, chosen_x); // Note: swapped x and y as per your original print order
                    io::stdout().flush().unwrap();
                    piece.clear();
                    anfield.clear();
                } else {
                    println!("0 0");
                    io::stdout().flush().unwrap();
                    piece.clear();
                    anfield.clear();
                }
            }
        }
    }
}

fn can_place_piece(
    board: &Vec<Vec<char>>,
    piece: &Vec<Vec<char>>,
    x: usize,
    y: usize,
    enemy: char,
    enemy2: char,
    p1: char,
    p2: char
) -> bool {
    let mut overlap_count = 0;

    for i in 0..piece.len() {
        for j in 0..piece[i].len() {
            let board_x = x + i;
            let board_y = y + j;

            // Boundary check
            if board_x >= board.len() || board_y >= board[0].len() {
                if piece[i][j] != '.' {
                    // Only consider non-dot parts of the piece for boundary checks
                    return false;
                }
                continue;
            }

            // overlap with enemy check
            if
                piece[i][j] != '.' &&
                (board[board_x][board_y] == enemy || board[board_x][board_y] == enemy2)
            {
                return false;
            }

            // Connection check
            if
                piece[i][j] != '.' &&
                (board[board_x][board_y] == p1 || board[board_x][board_y] == p2)
            {
                overlap_count += 1;
            }
        }
    }

    overlap_count == 1
}

fn place_piece_on_board(
    board: &Vec<Vec<char>>,
    piece: &Vec<Vec<char>>,
    player: &String
) -> Vec<(usize, usize)> {
    let enemy: char;
    let enemy2: char;
    let p1: char;
    let p2: char;
    let mut valid_positions = Vec::new();

    if player == "p1" {
        enemy = 's';
        enemy2 = '$';
        p1 = '@';
        p2 = 'a';
    } else {
        enemy = 'a';
        enemy2 = '@';
        p1 = 's';
        p2 = '$';
    }

    for x in 0..board.len() {
        for y in 0..board[0].len() {
            if can_place_piece(&board, &piece, x, y, enemy, enemy2, p1, p2) {
                valid_positions.push((x, y));
            }
        }
    }

    valid_positions
}

// Euclidean distance
fn distance(x1: usize, y1: usize, x2: usize, y2: usize) -> f64 {
    let dx = (x2 as f64) - (x1 as f64);
    let dy = (y2 as f64) - (y1 as f64);
    (dx * dx + dy * dy).sqrt()
}
