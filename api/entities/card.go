package entities

type CreditCard struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	ExpireMonth int    `json:"expire_month"`
	ExpireYear  int    `json:"expire_year"`
	CVV         int    `json:"cvv"`
}
