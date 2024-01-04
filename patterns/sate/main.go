package sate

import (
	"fmt"
	"log"
)

/*
Состояние — это поведенческий паттерн проектирования, который позволяет объектам
менять поведение в зависимости от своего состояния. Извне создаётся впечатление,
что изменился класс объекта.

Паттерн Состояние невозможно рассматривать в отрыве от концепции машины
состояний, также известной как стейт-машина или конечный автомат Основная идея в
том, что программа может находиться в одном из нескольких состояний, которые всё
время сменяют друг друга. Набор этих состояний, а также переходов между ними,
предопределён и конечен. Находясь в разных состояниях, программа может
по-разному реагировать на одни и те же события, которые происходят с ней.

Такой подход можно применить и к отдельным объектам. Например, объект Документ может принимать три состояния: Черновик, Модерация или Опубликован. В каждом из этих состоянии метод опубликовать будет работать по-разному:
Из черновика он отправит документ на модерацию.
Из модерации — в публикацию, но при условии, что это сделал администратор.
В опубликованном состоянии метод не будет делать ничего.

Паттерн Состояние предлагает создать отдельные классы для каждого состояния, в котором может пребывать объект, а затем вынести туда поведения, соответствующие этим состояниям.

Применимость:
Когда у вас есть объект, поведение которого кардинально меняется в зависимости от внутреннего состояния, причём типов состояний много, и их код часто меняется.
 - Паттерн предлагает выделить в собственные классы все поля и методы, связанные
 с определёнными состояниями. Первоначальный объект будет постоянно ссылаться на
 один из объектов-состояний, делегируя ему часть своей работы. Для изменения
 состояния в контекст достаточно будет подставить другой объект-состояние.

Когда код класса содержит множество больших, похожих друг на друга, условных операторов, которые выбирают поведения в зависимости от текущих значений полей класса.
 - Паттерн предлагает переместить каждую ветку такого условного оператора в
 собственный класс. Тут же можно поселить и все поля, связанные с данным
 состоянием.

Когда вы сознательно используете табличную машину состояний, построенную на условных операторах, но вынуждены мириться с дублированием кода для похожих состояний и переходов.
 - Паттерн Состояние позволяет реализовать иерархическую машину состояний,
 базирующуюся на наследовании. Вы можете отнаследовать похожие состояния от
 одного родительского класса и вынести туда весь дублирующий код.

Преимущества:
	Избавляет от множества больших условных операторов машины состояний.
	Концентрирует в одном месте код, связанный с определённым состоянием.
	Упрощает код контекста.

Недостатки:
	Может неоправданно усложнить код, если состояний мало и они редко меняются.
*/

/*
паттерн проектирования Состояние в контексте торговых автоматов. Для упрощения задачи представим, что торговый автомат может выдавать только один товар. Также представим, что автомат может пребывать только в одном из четырех состояний:

hasItem (имеетПредмет)
noItem (неИмеетПредмет)
itemRequested (выдаётПредмет)
hasMoney (получилДеньги)
*/

func main() {
	vendingMachine := newVendingMachine(1, 10)

	err := vendingMachine.requestItem()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.insertMoney(10)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.dispenseItem()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	err = vendingMachine.addItem(2)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	err = vendingMachine.requestItem()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.insertMoney(10)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.dispenseItem()
	if err != nil {
		log.Fatalf(err.Error())
	}

	//Item requestd
	//Money entered is ok
	//Dispensing Item
	//
	//Adding 2 items
	//
	//Item requestd
	//Money entered is ok
	//Dispensing Item

}

// интерфейс состояния
type State interface {
	addItem(int) error
	requestItem() error
	insertMoney(money int) error
	dispenseItem() error
}

// конкретный интерфейс
type NoItemState struct {
	vendingMachine *VendingMachine
}

func (i *NoItemState) requestItem() error {
	return fmt.Errorf("Item out of stock")
}

func (i *NoItemState) addItem(count int) error {
	i.vendingMachine.incrementItemCount(count)
	i.vendingMachine.setState(i.vendingMachine.hasItem)
	return nil
}

func (i *NoItemState) insertMoney(money int) error {
	return fmt.Errorf("Item out of stock")
}
func (i *NoItemState) dispenseItem() error {
	return fmt.Errorf("Item out of stock")
}

// конкретный интерфейс
type HasItemState struct {
	vendingMachine *VendingMachine
}

func (i *HasItemState) requestItem() error {
	if i.vendingMachine.itemCount == 0 {
		i.vendingMachine.setState(i.vendingMachine.noItem)
		return fmt.Errorf("No item present")
	}
	fmt.Printf("Item requestd\n")
	i.vendingMachine.setState(i.vendingMachine.itemRequested)
	return nil
}

func (i *HasItemState) addItem(count int) error {
	fmt.Printf("%d items added\n", count)
	i.vendingMachine.incrementItemCount(count)
	return nil
}

func (i *HasItemState) insertMoney(money int) error {
	return fmt.Errorf("Please select item first")
}
func (i *HasItemState) dispenseItem() error {
	return fmt.Errorf("Please select item first")
}

// конкретный интерфейс
type ItemRequestedState struct {
	vendingMachine *VendingMachine
}

func (i *ItemRequestedState) requestItem() error {
	return fmt.Errorf("Item already requested")
}

func (i *ItemRequestedState) addItem(count int) error {
	return fmt.Errorf("Item Dispense in progress")
}

func (i *ItemRequestedState) insertMoney(money int) error {
	if money < i.vendingMachine.itemPrice {
		return fmt.Errorf("Inserted money is less. Please insert %d", i.vendingMachine.itemPrice)
	}
	fmt.Println("Money entered is ok")
	i.vendingMachine.setState(i.vendingMachine.hasMoney)
	return nil
}
func (i *ItemRequestedState) dispenseItem() error {
	return fmt.Errorf("Please insert money first")
}

// конкретный интерфейс
type HasMoneyState struct {
	vendingMachine *VendingMachine
}

func (i *HasMoneyState) requestItem() error {
	return fmt.Errorf("Item dispense in progress")
}

func (i *HasMoneyState) addItem(count int) error {
	return fmt.Errorf("Item dispense in progress")
}

func (i *HasMoneyState) insertMoney(money int) error {
	return fmt.Errorf("Item out of stock")
}
func (i *HasMoneyState) dispenseItem() error {
	fmt.Println("Dispensing Item")
	i.vendingMachine.itemCount = i.vendingMachine.itemCount - 1
	if i.vendingMachine.itemCount == 0 {
		i.vendingMachine.setState(i.vendingMachine.noItem)
	} else {
		i.vendingMachine.setState(i.vendingMachine.hasItem)
	}
	return nil
}

// контекст
type VendingMachine struct {
	hasItem       State
	itemRequested State
	hasMoney      State
	noItem        State

	currentState State

	itemCount int
	itemPrice int
}

func newVendingMachine(itemCount, itemPrice int) *VendingMachine {
	v := &VendingMachine{
		itemCount: itemCount,
		itemPrice: itemPrice,
	}
	hasItemState := &HasItemState{
		vendingMachine: v,
	}
	itemRequestedState := &ItemRequestedState{
		vendingMachine: v,
	}
	hasMoneyState := &HasMoneyState{
		vendingMachine: v,
	}
	noItemState := &NoItemState{
		vendingMachine: v,
	}

	v.setState(hasItemState)
	v.hasItem = hasItemState
	v.itemRequested = itemRequestedState
	v.hasMoney = hasMoneyState
	v.noItem = noItemState
	return v
}

func (v *VendingMachine) requestItem() error {
	return v.currentState.requestItem()
}

func (v *VendingMachine) addItem(count int) error {
	return v.currentState.addItem(count)
}

func (v *VendingMachine) insertMoney(money int) error {
	return v.currentState.insertMoney(money)
}

func (v *VendingMachine) dispenseItem() error {
	return v.currentState.dispenseItem()
}

func (v *VendingMachine) setState(s State) {
	v.currentState = s
}

func (v *VendingMachine) incrementItemCount(count int) {
	fmt.Printf("Adding %d items\n", count)
	v.itemCount = v.itemCount + count
}
