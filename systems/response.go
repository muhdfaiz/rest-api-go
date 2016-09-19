package systems

type Response struct {
}

type ResponseFormat struct {
	Data `json:"data"`
}

type Data struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

// func (r Response) FormatResponseData(responseData interface{}) ResponseFormat {
// 	data := Data{}
// 	return nil
// }
