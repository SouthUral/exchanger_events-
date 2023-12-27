package router_test

import (
	"testing"
	"time"

	rt "github.com/SouthUral/exchangeTest/router"
	ut "github.com/SouthUral/exchangeTest/router/utilsfortest"
)

// Тест с проверкой базовой работоспособности маршрутизатора
// Один отправитель, один получатель (все события отправителя)
func Test1(t *testing.T) {
	// количество сообщений которое будет сгенерировано каждым типом
	// reports := make([]ut.ReportSubTest, 0, 30)

	numbersMsg := 10
	pass := "./testdata/fixt_1.json"

	conf, err := ut.LoadConf(pass)
	if err != nil {
		t.Fatal(err.Error())
	}
	eventCh, subCh, cancelRouter := rt.InitRouter()

	// получатели запускаются первыми
	cancelSub, reportCh := ut.SubscribersWork(conf.Consumers, subCh, numbersMsg)

	ctx := ut.StartPublishers(conf.Publishers, eventCh, numbersMsg)

	for {
		select {
		case rep := <-reportCh:
			if rep.AllEvent {
				checkNumsReceivedEvents(&rep, t)
			} else {
				checkNumsReceivedEvents(&rep, t)
				checkUnexpectedEvents(&rep, t)
				checkExpectedEvents(&rep, t)
				checkNumberMissedEvents(&rep, t, numbersMsg)
			}

		case <-ctx.Done():
			time.Sleep(20 * time.Millisecond)
			cancelSub()
			cancelRouter()
			break
		}
	}
}

// проверка, все ли события получены
func checkNumsReceivedEvents(rep *ut.ReportSubTest, t *testing.T) {
	resCount := rep.CounterReceivedMess
	expCount := rep.ExpectedNumberMess
	if resCount != expCount {
		t.Errorf("получатель %s ожидал %d количество событий, получено %d", rep.NameSub, expCount, resCount)
	}
}

// проверка, есть ли событи которые получатель не ожидает
func checkUnexpectedEvents(rep *ut.ReportSubTest, t *testing.T) {
	numsUnexpEvent := len(rep.UnexpectedListCountEvents)
	if numsUnexpEvent > 0 {
		t.Errorf("получатель %s получил %d неожидаемых типов событий", rep.NameSub, numsUnexpEvent)
	}
}

// проверка, все ли типы событий, которые ожидает получатель, были получены
func checkExpectedEvents(rep *ut.ReportSubTest, t *testing.T) {
	for key, countEvent := range rep.ExpectedlistCountEvents {
		if countEvent == 0 {
			t.Errorf("тип события %s, не был получен получателем %s", key, rep.NameSub)
		}
	}

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
