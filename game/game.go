package game

import (
	"fiber_api_v1/models"
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

type Board [3][3]string

type Game struct {
	Board Board
	bot   Bot
}

type GamePlay struct {
	UserName string
	Symbol   string
	isover   bool
	Winner   string
	Game
}

var games map[int]GamePlay = make(map[int]GamePlay)

func (b *Board) set(row, col int, marker string) (ok bool) {
	if b[row][col] != " " {
		return false
	}
	b[row][col] = marker
	return true
}
func (b *Board) isGameOver() (ok bool, winner string) {
	for row := 0; row < 3; row++ {
		if b[row][0] != " " && b[row][0] == b[row][1] && b[row][0] == b[row][2] {
			return true, b[row][0]
		}
	}

	for col := 0; col < 3; col++ {
		if b[0][col] != " " && b[0][col] == b[1][col] && b[0][col] == b[2][col] {
			return true, b[0][col]
		}
	}

	if b[0][0] != " " && b[0][0] == b[1][1] && b[0][0] == b[2][2] {
		return true, b[0][0]
	}
	if b[0][2] != " " && b[0][2] == b[1][1] && b[0][2] == b[2][0] {
		return true, b[0][2]
	}

	isDraw := true
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if b[row][col] == " " {
				isDraw = false
				break
			}
		}
	}

	return isDraw, "-"
}
func (g *GamePlay) IsBotWinner() (ok bool) {
	if g.isover {
		return g.UserName != g.Winner
	}
	return false
}
func (g *GamePlay) isGameOver() (ok bool, winner string) {

	if !g.isover {
		if ok, w_symbol := g.Board.isGameOver(); ok {

			if w_symbol == g.Symbol {
				g.Winner = g.UserName
			} else if w_symbol == "-" {
				g.Winner = "draw"
			} else {
				g.Winner = g.bot.getName()
			}
			g.isover = ok
		}
	}
	return g.isover, g.Winner
}

func StartHandler(c *fiber.Ctx) error {

	id := c.Locals("user_id").(int)
	bot := c.Query("bot")
	var response models.StartGameResponse

	// if game, ok := games[id]; ok && !game.isover {
	// 	return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"error": "Game hase alredy started, call /step to continue"})
	// } else {

	curgame := GamePlay{}
	curgame.UserName = c.Locals("username").(string)
	curgame.isover = false
	curgame.Winner = ""

	if bot == "Liza" {
		curgame.bot = NewLizaBot()
	} else if bot == "Alex" {
		curgame.bot = NewAlexBot()
	} else if bot == "Rostik" {
		curgame.bot = NewRostikBot()
	} else {
		curgame.bot = NewRandomBot()
	}

	if rand.Int63n(2) == 0 {
		curgame.Symbol = "0"
		curgame.bot.startGame("X")
		curgame.Board = [3][3]string{{" ", " ", " "}, {" ", "X", " "}, {" ", " ", " "}}
	} else {
		curgame.Symbol = "X"
		curgame.bot.startGame("0")
		curgame.Board = [3][3]string{{" ", " ", " "}, {" ", " ", " "}, {" ", " ", " "}}
	}
	response.Symbol = curgame.Symbol
	response.Board = curgame.Board
	games[id] = curgame
	//}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GameStepHandler(c *fiber.Ctx) error {
	id := c.Locals("user_id").(int)

	var request models.GameStepRequest
	var response models.GameStepResponse
	var curgame GamePlay

	if val, ok := games[id]; ok {
		curgame = val
	} else {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"error": "No Game started, call /start first"})
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if request.Col > 3 || request.Row > 3 || request.Col < 0 || request.Row < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "row and col indexes shoul be more than 0 and less 3"})
	}
	if curgame.Board[request.Row][request.Col] != " " {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("position row:%d col:%d is taken please choose another ", request.Row, request.Col)})
	} else {
		curgame := games[id]
		var winner string
		var isover bool

		if isover, winner = curgame.isGameOver(); !isover {
			curgame.Board.set(request.Row, request.Col, curgame.Symbol)
			if isover, winner = curgame.isGameOver(); !isover {
				curgame.Board.set(curgame.bot.step(curgame.Board))
				isover, winner = curgame.isGameOver()
			}
		}
		response.IsOver = isover
		response.WinnerName = winner
		response.Board = curgame.Board
		response.IsBotWinner = curgame.IsBotWinner()
		games[id] = curgame
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
