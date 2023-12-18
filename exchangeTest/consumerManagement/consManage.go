package consumermanagement

import (
	"context"

	api "github.com/SouthUral/exchangeTest/api/api_v1"

	models "github.com/SouthUral/exchangeTest/models"
)

// Инициализирует менеджера получателей, который управляет внутренними получателями
// TODO: нужно передать сюда состояние, которое будет выгружено из БД при запуске
// либо придет пустая структура (значит что состояния нет в БД), либо придет структура и по ней будут инициализированы внутренние получатели

// TODO: Нужно сделать три канала которые будут переданы сюда от модуля состояния:
// Каналы на отправку:
// 1-й канал - нужен для запросов на восстановление очереди // по каналу отправляется запрос со структурой в которой будет канал для обратного возврата событий
// 2-й канал - нужен для отправки измененного состояния получателей (внутренний и промежуточный) и для отправки данных мониторинга очередей
// 3-й канал - нужен для записи событий в БД

// Каналы на получение:
// TODO: Нужен канал от HTTP API, по которому будут поступать команды для получателей
func InitConsManager(apiCh <-chan models.SubMess, subCh chan<- models.SubscriberMess) func() {
	storageInfoCons := make(map[string]internalConsum)

	ctx, cancel := context.WithCancel(context.Background())

	// TODO: здесь должен быть метод инициализации внутренних получателей

	go func() {
		for {
			select {
			case mess := <-apiCh:
				// прием сообщение от API gRPC
				switch mess.TypeMess {

				case api.ConnectingConsumer:
					// тип сообщения (подключение потребителя)
					confKey := mess.Conf.ConfSubscribe.GenStringVue()
					cons, ok := internalConsumers[confKey]
					// если внутренний получатель с такой конфигурацией уже есть то ему добавляется получатель
					if ok {
						cons.AddConsumer(subscriber{
							name:  mess.Conf.Name,
							subCh: mess.RevCh,
						})
					} else {
						internalSub := InitInternalConsumer(mess.Conf.ConfSubscribe, subCh)
						internalSub.AddConsumer(subscriber{
							name:  mess.Conf.Name,
							subCh: mess.RevCh,
						})
					}
					// TODO: нужно реализовать проверки конфига, есть ли уже работающий внутренний потребитель с таким конфигом

					// TODO: если конфиг найден, нужна проверка, есть ли этого конфига такой подписчик, если нет, то нужно создать подписчика
					// TODO: если подписчик есть, нужна проверка в каком он состоянии, если подписчик активен, то нужно отправить ошибку, что такой подписчик уже есть.
					// TODO: если подписчик есть и он не активен, то нужно привязать канал к этому подписчику и перевести подписчика в состояние активен

				case api.StoppingSending:
					// тип сообщения (остановка отправления событий)
					// перевод подписчика в состояние ожидания

				}
			case <-done:
				return
			}
		}
	}()
	// TODO: должен принимать канал для связи с сервером, по которому будут передаваться данные о новых потребителях или об их остановке
	// TODO: должен принимать канал для связи с API для управления и мониторинга получателей
	// TODO: должен принимать канал для отрпавки сообщений регистрации получателей в маршрутизаторе
	return cancel
}

// Менеджер внутренних получателей
type ConsManager struct {
	InternalConsumers map[string]InternalConsumer  // словарь для хранения экземпляров внутренних получателей
	ApiCh             <-chan models.SubMess        // канал по которому приходят события от gRPC API
	SubCh             chan<- models.SubscriberMess // канал в который нужно отправить сообщение для регистрации получателя в router
}

// Запуск работы менеджера получателей
func (M *ConsManager) runWork(ctx context.Context) {

}

// Добавление внутреннего получателя в ConsManager.InternalConsumers
func (M *ConsManager) addInternalConsumer(conf models.ConfSub) {

}
