package systems

type Transformer struct {
}

type CreateTransformer struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
	Links      map[string]string `json:"links"`
}
