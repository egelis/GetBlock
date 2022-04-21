package common

type MostChanged struct {
	Address string `json:"address"`
	Value   string `json:"value"`
}

type Block struct {
	*BlockContent `json:"result"`
}

type BlockContent struct {
	BlockNumber  string         `json:"number"`
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type LastBlockNum struct {
	Number string `json:"result"`
}
