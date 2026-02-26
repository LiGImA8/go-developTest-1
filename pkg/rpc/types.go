package rpc

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid    bool   `json:"valid"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type PlaceOrderRequest struct {
	Token    string `json:"token"`
	ItemName string `json:"item_name"`
	Quantity int32  `json:"quantity"`
}

type PlaceOrderResponse struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}
