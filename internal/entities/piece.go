package entities

type Piece struct {
	Content string   `json:"content"`
	Proofs  []string `json:"proof"`
}
