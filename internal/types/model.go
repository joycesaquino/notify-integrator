package types

// TODO Aqui usa a tipagem que quiser para ser o body enviado para a API da integração.
type Body struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Idade int    `json:"idade"`
}
