package exchangers

// ошибка конвертации собития пришедшего от маршрутизатора
type eventConversionError struct {
}

func (e eventConversionError) Error() string {
	return "it is not possible to convert an event to an interface"
}
