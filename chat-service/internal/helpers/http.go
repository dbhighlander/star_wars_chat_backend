package helpers

type Response struct {
	Result  string      `json:"result"`
	Details interface{} `json:"details"`
}
