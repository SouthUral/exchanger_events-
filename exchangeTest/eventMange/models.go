package eventmange

import (
	rt "github.com/SouthUral/exchangeTest/router"
)

// Канал для отправки событий в БД, временно находится здесь
type EventDBCh chan rt.Event
