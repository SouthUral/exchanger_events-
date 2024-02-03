package router_test

import (
	"testing"
	"time"

	rt "github.com/SouthUral/exchangerEvents/internal/core/router"
	ut "github.com/SouthUral/exchangerEvents/internal/core/router/utilsfortest"
)

// Тест с проверкой базовой работоспособности маршрутизатора
// Один отправитель, один получатель (все события отправителя)
func Test_complex(t *testing.T) {
	// количество сообщений которое будет сгенерировано каждым типом
	// reports := make([]ut.ReportSubTest, 0, 30)

	data := []struct {
		name       string
		pass       string
		numbersMsg int
		numsEvents map[string]int
	}{
		{"test_1", "./testdata/fixt_1.json", 10, map[string]int{"consumer_1": 20}},
		{"test_2", "./testdata/fixt_2.json", 10, map[string]int{"consumer_1": 20, "consumer_2": 20, "consumer_3": 20}},
		{"test_3", "./testdata/fixt_3.json", 10, map[string]int{"consumer_1": 70}},
		{"test_4", "./testdata/fixt_4.json", 10, map[string]int{"consumer_1": 20, "consumer_2": 30, "consumer_3": 10, "consumer_4": 10, "consumer_5": 70}},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			conf, err := ut.LoadConf(d.pass)
			if err != nil {
				t.Fatal(err.Error())
			}
			eventCh, subCh, cancelRouter := rt.InitRouter()

			// получатели запускаются первыми
			numsConsConf := len(conf.Consumers)
			reportCh, ctxSub := ut.SubscribersWork(conf.Consumers, subCh, d.numbersMsg)
			time.Sleep(2 * time.Second)

			ut.StartPublishers(conf.Publishers, eventCh, d.numbersMsg)
			var numsReport int
			<-ctxSub.Done()

			for i := 0; i < numsConsConf; i++ {
				select {
				case rep := <-reportCh:
					numsEvent, ok := d.numsEvents[rep.NameSub]
					if !ok {
						t.Errorf("ошибка тестов, %s нет в списке теста", rep.NameSub)
					}
					numsReport++
					if rep.AllEvent {
						checkNumsReceivedEvents(&rep, t, numsEvent)
					} else {
						checkNumsReceivedEvents(&rep, t, numsEvent)
						checkUnexpectedEvents(&rep, t)
						checkExpectedEvents(&rep, t)
					}
				}
			}
			cancelRouter()
		})
	}
}

// проверка, все ли события получены
func checkNumsReceivedEvents(rep *ut.ReportSubTest, t *testing.T, expCount int) bool {
	res := true
	resCount := rep.CounterReceivedMess
	if resCount != expCount {
		t.Errorf("получатель %s ожидал %d количество событий, получено %d", rep.NameSub, expCount, resCount)
		res = false
		return res
	}

	t.Logf("#1 Проверка количества всех полученных событий пройдена! Получатель %s. Ожидалось %d / получено %d", rep.NameSub, expCount, resCount)
	return true
}

// проверка, есть ли событи которые получатель не ожидает
func checkUnexpectedEvents(rep *ut.ReportSubTest, t *testing.T) {
	numsUnexpEvent := len(rep.UnexpectedListCountEvents)
	if numsUnexpEvent > 0 {
		t.Errorf("получатель %s получил %d неожидаемых типов событий", rep.NameSub, numsUnexpEvent)
		return
	}
	t.Logf("#2 Проверка, есть ли событи которые получатель не ожидает пройдена! Получатель %s", rep.NameSub)
}

// проверка, все ли типы событий, которые ожидает получатель, были получены
func checkExpectedEvents(rep *ut.ReportSubTest, t *testing.T) {
	for key, countEvent := range rep.ExpectedlistCountEvents {
		if countEvent == 0 {
			t.Errorf("тип события %s, не был получен получателем %s", key, rep.NameSub)
			return
		}
	}
	t.Logf("#3 Проверка, все ли типы событий, которые ожидает получатель, были получены пройдена! Получатель %s", rep.NameSub)
}
