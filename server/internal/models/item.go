package models

type Item struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
}

type Result struct {
	Distance  int    `json:"distance"`
	Signature string `json:"signature"`
}
