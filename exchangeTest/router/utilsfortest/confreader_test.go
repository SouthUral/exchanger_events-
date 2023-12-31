package utilsfortest

import (
	"testing"
)

func TestConfReader1(t *testing.T) {
	pub := Publisher{
		Name: "publisher_1",
		TypeMess: []string{
			"type_1",
			"type_2",
		},
	}

	pubSub := Publisher{
		Name: "publisher_1",
	}

	confTest := Config{
		Publishers: []Publisher{
			pub,
		},
		Consumers: []Consumer{
			Consumer{
				Name: "consumer_1",
				Publishers: []Publisher{
					pubSub,
				},
			},
		},
	}

	conf, err := LoadConf("./testdata/fixt_1.json")
	if err != nil {
		t.Fatal(err)
	}

	checkConsumers(confTest.Consumers, conf.Consumers, t)
	checkPublishers(confTest.Publishers, conf.Publishers, t)
}

func TestConfReader2(t *testing.T) {
	publishers := []Publisher{
		Publisher{
			Name: "publisher_1",
			TypeMess: []string{
				"type_1",
				"type_2",
			},
		},
		Publisher{
			Name: "publisher_2",
			TypeMess: []string{
				"type_2",
				"type_3",
			},
		},
	}
	consumers := []Consumer{
		Consumer{
			Name: "consumer_1",
			Publishers: []Publisher{
				Publisher{
					Name: "publisher_1",
				},
			},
		},
		Consumer{
			Name: "consumer_2",
			Publishers: []Publisher{
				Publisher{
					Name: "publisher_2",
					TypeMess: []string{
						"type_3",
					},
				},
			},
		},
		Consumer{
			Name: "consumer_3",
			Types: []string{
				"type_2",
			},
		},
	}

	confTest := Config{
		Consumers:  consumers,
		Publishers: publishers,
	}

	conf, err := LoadConf("./testdata/fixt_2.json")
	if err != nil {
		t.Fatal(err)
	}

	checkConsumers(confTest.Consumers, conf.Consumers, t)
	checkPublishers(confTest.Publishers, conf.Publishers, t)
}

func checkPublishers(publishersTest, publishersLoad []Publisher, t *testing.T) {
	numsPubTest := len(publishersTest)
	numsPubLoad := len(publishersLoad)
	if numsPubTest != numsPubLoad {
		t.Errorf("Не совпадает количество отправителей, ожидается %d, получено %d", numsPubTest, numsPubLoad)
	}

	for idx, testPub := range publishersTest {
		loadPub := publishersLoad[idx]
		if testPub.Name != loadPub.Name {
			t.Errorf("ожидалось имя отправителя %s, получено %s", testPub.Name, loadPub.Name)
		} else {
			if len(testPub.TypeMess) != len(loadPub.TypeMess) {
				t.Errorf("Количество типов событий у отправителя %s не совпадает", testPub.Name)
			}
			for idxType, typeTest := range testPub.TypeMess {
				typeLoad := loadPub.TypeMess[idxType]
				if typeTest != typeLoad {
					t.Errorf("у отправителя %s не совпадает тип события, получен %s, ожидалось %s", testPub.Name, typeLoad, typeTest)
				}
			}
		}
	}
}

func checkConsumers(consumersTest, consumersLoad []Consumer, t *testing.T) {
	numsConsTest := len(consumersTest)
	numsConsLoad := len(consumersLoad)
	if numsConsTest != numsConsLoad {
		t.Errorf("Не совпадает количество подписчиков, ожидается %d, получено %d", numsConsTest, numsConsLoad)
	}

	for idx, consumer := range consumersLoad {
		consumerTest := consumersTest[idx]
		if consumer.Name != consumerTest.Name {
			t.Errorf("имя подписчика не совпадает %s != %s", consumer.Name, consumerTest.Name)
		}
		if consumerTest.AllEvent != consumer.AllEvent {
			t.Errorf("не совпдает флаг AllEvent, ожидалось получить %t", consumerTest.AllEvent)
		}
		if consumerTest.Publishers != nil {
			if consumer.Publishers != nil {
				numsSubPubTest := len(consumerTest.Publishers)
				numsSubPubLoad := len(consumer.Publishers)
				if numsSubPubTest != numsSubPubLoad {
					t.Errorf("не совпадает количество отправителей у подписчика %s, ожидалось %d, получено %d", consumerTest.Name, numsSubPubTest, numsSubPubLoad)
				}
				for idxPub, PubTest := range consumerTest.Publishers {
					Pub := consumer.Publishers[idxPub]
					if PubTest.Name != Pub.Name {
						t.Errorf("не совпадает имя отправителя %s у подписчика %s, получено %s", PubTest.Name, consumerTest.Name, Pub.Name)
					} else {
						if PubTest.TypeMess != nil {
							if Pub.TypeMess != nil {
								numsSubPubTypesTest := len(PubTest.TypeMess)
								numsSubPubTypesLoad := len(Pub.TypeMess)
								if numsSubPubTypesTest != numsSubPubTypesLoad {
									t.Errorf("не совпадает количество типов у %s.%s, ожидается %d, получено %d", consumerTest.Name, PubTest.Name, numsSubPubTypesTest, numsSubPubTypesLoad)
								}
								for idxType, TypeTest := range PubTest.TypeMess {
									Type := Pub.TypeMess[idxType]
									if TypeTest != Type {
										t.Errorf("ожидается тип %s.%s.%s, получен %s.%s.%s", consumerTest.Name, PubTest.Name, TypeTest, consumer.Name, Pub.Name, Type)
									}
								}
							} else {
								t.Errorf("%s.%s должен быть ненулевой список типов событий", consumerTest.Name, PubTest.Name)
							}
						} else {
							if Pub.TypeMess != nil {
								t.Errorf("%s.%s не должно быть списка типов", consumerTest.Name, PubTest.Name)
							}
						}
					}
				}
			} else {
				t.Errorf("у подписчика %s должен быть не нулевой список Publishers", consumerTest.Name)
			}
		} else {
			if consumer.Publishers != nil {
				t.Errorf("у подписчика %s не должно быть списка Publishers", consumerTest.Name)
			}
		}

		if consumerTest.Types != nil {
			if consumer.Types != nil {
				numsConsTypesTest := len(consumerTest.Types)
				numsConsTypesLoad := len(consumer.Types)
				if numsConsTypesTest != numsConsTypesLoad {
					t.Errorf("не совпадает количество типов у подписчика %s, ожидается %d, получено %d", consumerTest.Name, numsConsTypesTest, numsConsTypesLoad)
				}
				for idxType, TestType := range consumerTest.Types {
					Type := consumer.Types[idxType]
					if TestType != Type {
						t.Errorf("типы у подписчика %s не совпадают, ожидалось %s, получено %s", consumerTest.Name, TestType, Type)
					}
				}

			} else {
				t.Errorf("у подписчика %s должен быть ненулевой список типов", consumerTest.Name)
			}
		} else {
			if consumer.Types != nil {
				t.Errorf("у подписчика %s не должно быть списка типов", consumerTest.Name)
			}
		}
	}
}
