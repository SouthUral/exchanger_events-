package router_test

import (
	"testing"
	"time"

	rt "github.com/SouthUral/exchangeTest/router"
	ut "github.com/SouthUral/exchangeTest/router/utilsfortest"
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
	}{
		{"test_1", "./testdata/fixt_1.json", 100},
		{"test_2", "./testdata/fixt_2.json", 100},
		{"test_3", "./testdata/fixt_3.json", 100},
		{"test_4", "./testdata/fixt_4.json", 100},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			conf, err := ut.LoadConf(d.pass)
			if err != nil {
				t.Fatal(err.Error())
			}
			eventCh, subCh, cancelRouter := rt.InitRouter()

			// получатели запускаются первыми
			cancelSub, reportCh := ut.SubscribersWork(conf.Consumers, subCh, d.numbersMsg)

			ctx := ut.StartPublishers(conf.Publishers, eventCh, d.numbersMsg)

			for {
				select {
				case rep := <-reportCh:
					if rep.AllEvent {
						checkNumsReceivedEvents(&rep, t)
					} else {
						checkNumsReceivedEvents(&rep, t)
						checkUnexpectedEvents(&rep, t)
						checkExpectedEvents(&rep, t)
						checkNumberMissedEvents(&rep, t, d.numbersMsg)
					}

				case <-ctx.Done():
					time.Sleep(20 * time.Millisecond)
					cancelSub()
					cancelRouter()
					return
				}
			}
		})
	}
}

// проверка, все ли события получены
func checkNumsReceivedEvents(rep *ut.ReportSubTest, t *testing.T) {
	resCount := rep.CounterReceivedMess
	expCount := rep.ExpectedNumberMess
	if resCount != expCount {
		t.Errorf("получатель %s ожидал %d количество событий, получено %d", rep.NameSub, expCount, resCount)
		return
	}

	t.Logf("Проверка количества всех полученных событий пройдена! Получатель %s", rep.NameSub)
}

// проверка, есть ли событи которые получатель не ожидает
func checkUnexpectedEvents(rep *ut.ReportSubTest, t *testing.T) {
	numsUnexpEvent := len(rep.UnexpectedListCountEvents)
	if numsUnexpEvent > 0 {
		t.Errorf("получатель %s получил %d неожидаемых типов событий", rep.NameSub, numsUnexpEvent)
		return
	}
	t.Logf("Проверка, есть ли событи которые получатель не ожидает пройдена! Получатель %s", rep.NameSub)
}

// проверка, все ли типы событий, которые ожидает получатель, были получены
func checkExpectedEvents(rep *ut.ReportSubTest, t *testing.T) {
	for key, countEvent := range rep.ExpectedlistCountEvents {
		if countEvent == 0 {
			t.Errorf("тип события %s, не был получен получателем %s", key, rep.NameSub)
			return
		}
	}
	t.Logf("Проверка, все ли типы событий, которые ожидает получатель, были получены пройдена! Получатель %s", rep.NameSub)
}

// проверка числа недополученных событий у каждого типа событий
func checkNumberMissedEvents(rep *ut.ReportSubTest, t *testing.T, numbersMsg int) {
	for key, countEvent := range rep.ExpectedlistCountEvents {
		if countEvent != numbersMsg {
			t.Errorf("получатель %s ожидал получить %d количество событий типа %s, было получено %d событий",
				rep.NameSub,
				numbersMsg,
				key,
				countEvent,
			)
		}
	}
}
