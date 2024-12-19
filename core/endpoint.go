package core

type Endpoint struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  bool   `json:"status"`
	Date    int64  `json:"date"`
}

type Setter interface {
	Save(e Endpoint) error
}
