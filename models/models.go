package models

// LoginRequest — структура для запроса авторизации
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse — структура для ответа после авторизации
type LoginResponse struct {
	Message string `json:"message"`
	Token   string
}

// RegisterRequest - структура для запроса на регистрацию нового пользователя
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// StartGameResponse
type StartGameResponse struct {
	Symbol string `json:"symbol"`
	Board  [3][3]string
}

// GameStepRequest
type GameStepRequest struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// GameStepResponse
type GameStepResponse struct {
	Board       [3][3]string
	IsOver      bool   `json:"isover"`
	WinnerName  string `json:"winnername"`
	IsBotWinner bool   `json:"isbotwinner"`
}
