package CoR

import "fmt"

/*
Цепочка обязанностей — это поведенческий паттерн проектирования, который
позволяет передавать запросы последовательно по цепочке обработчиков. Каждый
последующий обработчик решает, может ли он обработать запрос сам и стоит ли
передавать запрос дальше по цепи.
Цепочка обязанностей базируется на том, чтобы превратить отдельные поведения в объекты.

Паттерн предлагает связать объекты обработчиков в одну цепь. Каждый из них будет
иметь ссылку на следующий обработчик в цепи. Таким образом, при получении
запроса обработчик сможет не только сам что-то с ним сделать, но и передать
обработку следующему объекту в цепочке.

Применимость:
Когда программа должна обрабатывать разнообразные запросы несколькими способами, но заранее неизвестно, какие конкретно запросы будут приходить и какие обработчики для них понадобятся.
	- С помощью Цепочки обязанностей вы можете связать потенциальных обработчиков в одну цепь и при получении запроса поочерёдно спрашивать каждого из них, не хочет ли он обработать запрос.
Когда важно, чтобы обработчики выполнялись один за другим в строгом порядке.
	- Цепочка обязанностей позволяет запускать обработчиков последовательно один за другим в том порядке, в котором они находятся в цепочке.
Когда набор объектов, способных обработать запрос, должен задаваться динамически.
	- В любой момент вы можете вмешаться в существующую цепочку и переназначить связи так, чтобы убрать или добавить новое звено.

Преимущества:
- Уменьшает зависимость между клиентом и обработчиками.
- Реализует принцип единственной обязанности.
- Реализует принцип открытости/закрытости.

Недостатки:
- Запрос может остаться никем не обработанным.

*/

/*
Цепочка обязанностей на примере приложения больницы. Госпиталь может иметь разные помещения, например:
Приемное отделение
Доктор
Комната медикаментов
Кассир
*/

func main() {

	cashier := &Cashier{}

	//Set next for medical department
	medical := &Medical{}
	medical.setNext(cashier)

	//Set next for doctor department
	doctor := &Doctor{}
	doctor.setNext(medical)

	//Set next for reception department
	reception := &Reception{}
	reception.setNext(doctor)

	patient := &Patient{name: "abc"}
	//Patient visiting
	reception.execute(patient)

	//Reception registering patient
	//Doctor checking patient
	//Medical giving medicine to patient
	//Cashier getting money from patient patient
}

// объект
type Patient struct {
	name              string
	registrationDone  bool
	doctorCheckUpDone bool
	medicineDone      bool
	paymentDone       bool
}

// интерфейс обработчика
type Department interface {
	execute(*Patient)
	setNext(Department)
}

// конкретный обработчик
type Reception struct {
	next Department
}

func (r *Reception) execute(p *Patient) {
	if p.registrationDone {
		fmt.Println("Patient registration already done")
		r.next.execute(p)
		return
	}
	fmt.Println("Reception registering patient")
	p.registrationDone = true
	r.next.execute(p)
}

func (r *Reception) setNext(next Department) {
	r.next = next
}

// конкретный обработчик
type Doctor struct {
	next Department
}

func (d *Doctor) execute(p *Patient) {
	if p.doctorCheckUpDone {
		fmt.Println("Doctor checkup already done")
		d.next.execute(p)
		return
	}
	fmt.Println("Doctor checking patient")
	p.doctorCheckUpDone = true
	d.next.execute(p)
}

func (d *Doctor) setNext(next Department) {
	d.next = next
}

// конкретный обработчик
type Medical struct {
	next Department
}

func (m *Medical) execute(p *Patient) {
	if p.medicineDone {
		fmt.Println("Medicine already given to patient")
		m.next.execute(p)
		return
	}
	fmt.Println("Medical giving medicine to patient")
	p.medicineDone = true
	m.next.execute(p)
}

func (m *Medical) setNext(next Department) {
	m.next = next
}

// конкретный обработчик
type Cashier struct {
	next Department
}

func (c *Cashier) execute(p *Patient) {
	if p.paymentDone {
		fmt.Println("Payment Done")
	}
	fmt.Println("Cashier getting money from patient patient")
}

func (c *Cashier) setNext(next Department) {
	c.next = next
}
