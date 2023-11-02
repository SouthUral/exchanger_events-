package confreader

// Структура для парсинга json
type Config struct {
	Publishers []Publisher `json:"publishers"`
	Consumers  []Consumer  `json:"consumers"`
}

type Consumer struct {
	Name       string      `json:"name"`
	Publishers []Publisher `json:"publishers"`
	Types      []string    `json:"types"`
	AllEvent   bool        `json:"all_event"`
}

type Publisher struct {
	Name     string   `json:"name"`
	TypeMess []string `json:"type_mess"`
}
