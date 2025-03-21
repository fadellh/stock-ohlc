package ohlcEntity

type OrderType string

const (
	A OrderType = "A"
	P OrderType = "P"
	E OrderType = "E"
)

type OhlcMessage struct {
	OrderBook        string    `json:"order_book"`
	Price            int       `json:"price,string,omitempty"`
	StockCode        string    `json:"stock_code"`
	Type             OrderType `json:"type"`
	Quantity         int       `json:"quantity,string,omitempty"`
	ExecutedQuantity int       `json:"executed_quantity,string,omitempty"`
	ExecutionPrice   int       `json:"execution_price,string,omitempty"`
}

type OhlcStock struct {
	StockCode     string    `json:"stock_code"`
	PreviousPrice int       `json:"prev"`
	OpenPrice     int       `json:"open"`
	HighestPrice  int       `json:"highest"`
	LowestPrice   int       `json:"lowest"`
	ClosePrice    int       `json:"close"`
	Volume        int64     `json:"volume"`
	Value         string    `json:"value"`
	AveragePrice  int       `json:"average"`
	Type          OrderType `json:"type"`
}
