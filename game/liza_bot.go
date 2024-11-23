package game

type LizaBot struct {
	marker string
}

func NewLizaBot() *LizaBot {
	return &LizaBot{"?"}
}

func (b *LizaBot) getName() string {
	return "LizaBot"
}
func (b *LizaBot) startGame(marker string) {
	b.marker = marker
}
func (b *LizaBot) getMarker() string {
	return b.marker
}
func (b *LizaBot) step(board Board) (row, col int, marker string) {
	var prior int = 0
	var vrem_prior int = 0
	var max_prior bool = false
	marker = b.marker
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			vrem_prior, max_prior = prioriti(board, i, j, b.marker)
			if max_prior {
				row = i
				col = j
				return
			}
			if vrem_prior > prior {
				row = i
				col = j
				prior = vrem_prior
			}
		}
	}
	return
}

func prioriti(board Board, i, j int, marker string) (prior int, max_prior bool) {
	max_prior = false
	if board[i][j] != " " {
		return 0, max_prior
	}
	var prior_row, prior_col, prior_dogonali int = 0, 0, 0
	for n := 0; n < 3; n++ {
		if marker == board[i][n] { //Проверка для колонок
			prior_col += 5
		} else if board[i][n] == " " {
			prior_col++
		} else {
			prior_col += -4
		}
		if marker == board[n][j] { //Проверка для строк
			prior_row += 5
		} else if board[n][j] == " " {
			prior_row++
		} else {
			prior_row += -4
		}
		if i != 1 && j != 1 {
			if marker == board[n][n] { //Проверка для диогоналей
				prior_dogonali += 5
			} else if board[n][n] == " " {
				prior_dogonali++
			} else {
				prior_dogonali += -4
			}
		}
	}
	if prior_col == 13 || prior_row == 13 || prior_dogonali == 13 {
		max_prior = true
		return
	} else {
		prior = absInt(prior_col) + absInt(prior_row) + absInt(prior_dogonali)
		return
	}
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
