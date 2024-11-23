package game

type Bot interface {
	getName() string
	startGame(marker string)
	getMarker() string
	step(board Board) (row, col int, marker string)
}

type AlexBot struct {
	marker string
}

func NewAlexBot() *AlexBot {
	return &AlexBot{"?"}
}

func (b *AlexBot) getName() string {
	return "Папин бот"
}
func (b *AlexBot) startGame(marker string) {
	b.marker = marker
}
func (b *AlexBot) getMarker() string {
	return b.marker
}
func getpriority(board Board, row, col int, marker string) (res int) {
	res = 0
	var row_weight, col_weight, diag1_weight, diag2_weight int
	if row == 1 && col == 1 {
		res = 5
	}

	if row == col {
		for i := 0; i < 3; i++ {
			if row == i {
				continue
			}
			if board[i][i] == marker {
				diag1_weight += 5
			} else if board[i][i] != " " {
				diag1_weight += -4
			}
		}
		if diag1_weight < 0 {
			diag1_weight *= -1
		}
	}

	if row == 2-col {
		for i := 0; i < 3; i++ {
			if row == i {
				continue
			}
			if board[i][2-i] == marker {
				diag2_weight += 5
			} else if board[i][2-i] != " " {
				diag2_weight += -4
			}
		}
		if diag2_weight < 0 {
			diag2_weight *= -1
		}
	}

	for i := 0; i < 3; i++ {
		if i == row {
			continue
		}
		if board[i][col] == marker {
			row_weight += 5
		} else if board[i][col] != " " {
			row_weight += -4
		}
	}
	if row_weight < 0 {
		row_weight *= -1
	}

	for j := 0; j < 3; j++ {
		if j == col {
			continue
		}
		if board[row][j] == marker {
			col_weight += 5
		} else if board[row][j] != " " {
			col_weight += -4
		}
	}
	if col_weight < 0 {
		col_weight *= -1
	}

	// if res < row_weight {
	// 	res = row_weight
	// }
	// if res < col_weight {
	// 	res = col_weight
	// }

	// if res < diag1_weight {
	// 	res = diag1_weight
	// }
	// if res < diag2_weight {
	// 	res = diag2_weight
	// }
	col_weight *= col_weight
	row_weight *= row_weight
	diag1_weight *= diag1_weight
	diag2_weight *= diag2_weight

	res = res + col_weight + row_weight + diag1_weight + diag2_weight

	return
}

func (b *AlexBot) step(board Board) (row, col int, marker string) {

	var cur_priority int
	marker = b.marker

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] != " " {
				continue
			}

			priority := getpriority(board, i, j, marker)
			if cur_priority < priority {
				cur_priority = priority
				row = i
				col = j
			}
		}
	}
	return
}
