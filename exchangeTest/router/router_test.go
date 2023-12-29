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
		{"test_1", "./testdata/fixt_1.json", 10},
		{"test_2", "./testdata/fixt_2.json", 10},
		{"test_3", "./testdata/fixt_3.json", 10},
		{"test_4", "./testdata/fixt_4.json", 10},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			conf, err := ut.LoadConf(d.pass)
			if err != nil {
				t.Fatal(err.Error())
			}
			eventCh, subCh, cancelRouter := rt.InitRouter()

			// получатели запускаются первыми
			reportCh, ctxSub := ut.SubscribersWork(conf.Consumers, subCh, d.numbersMsg)

			ut.StartPublishers(conf.Publishers, eventCh, d.numbersMsg)
			var numsReport int
			for {
				select {
				case rep := <-reportCh:
					numsReport++
					if rep.AllEvent {
						if !checkNumsReceivedEvents(&rep, t) {
							t.Fatalf("получатель %s ожидал %d количество событий, получено %d", rep.NameSub, rep.ExpectedNumberMess, rep.CounterReceivedMess)
						}
					} else {
						checkNumsReceivedEvents(&rep, t)
						checkUnexpectedEvents(&rep, t)
						checkExpectedEvents(&rep, t)
						checkNumberMissedEvents(&rep, t, d.numbersMsg)
					}

				case <-ctxSub.Done():
					time.Sleep(5 * time.Second)
					cancelRouter()
					if numsReport == 0 {
						t.Error("Ни одного отчета не было получено")
					}
					return
				}
			}
		})
	}
}

// проверка, все ли события получены
func checkNumsReceivedEvents(rep *ut.ReportSubTest, t *testing.T) bool {
	res := true
	resCount := rep.CounterReceivedMess
	expCount := rep.ExpectedNumberMess
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

// проверка числа недополученных событий у каждого типа событий
func checkNumberMissedEvents(rep *ut.ReportSubTest, t *testing.T, numbersMsg int) {
	res := true
	for key, countEvent := range rep.ExpectedlistCountEvents {
		if countEvent != numbersMsg {
			t.Errorf("получатель %s ожидал получить %d количество событий типа %s, было получено %d событий",
				rep.NameSub,
				numbersMsg,
				key,
				countEvent,
			)
			res = false
		}
	}
	if res {
		t.Logf("#4 Проверка получателя %s количества ожидаемых событий по каждому типу пройдена", rep.NameSub)
	}
}
