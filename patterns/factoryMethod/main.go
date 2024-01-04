package factoryMethod

import "fmt"

/*
Фабричный метод — это порождающий паттерн проектирования, который определяет
общий интерфейс для создания объектов в суперклассе, позволяя подклассам
изменять тип создаваемых объектов.

Паттерн Фабричный метод предлагает создавать объекты не напрямую, используя оператор new, а через вызов особого фабричного метода.
теперь вы сможете переопределить фабричный метод в подклассе, чтобы изменить тип создаваемого продукта.

Чтобы эта система заработала, все возвращаемые объекты должны иметь общий
интерфейс. Подклассы смогут производить объекты различных классов, следующих
одному и тому же интерфейсу.

Применимость:
Когда заранее неизвестны типы и зависимости объектов, с которыми должен работать код.
	- Фабричный метод отделяет код производства продуктов от
	остального кода, который эти продукты использует. Благодаря этому, код
	производства можно расширять, не трогая основной. Так, чтобы добавить поддержку
	нового продукта, вам нужно создать новый подкласс и определить в нём фабричный
	метод, возвращая оттуда экземпляр нового продукта.

Когда вы хотите дать возможность пользователям расширять части вашего фреймворка или библиотеки.

Когда вы хотите экономить системные ресурсы, повторно используя уже созданные объекты, вместо порождения новых.
	-Такая проблема обычно возникает при работе с тяжёлыми ресурсоёмкими объектами, такими, как подключение к базе данных, файловой системе и т. д.

Преимущества:
- Избавляет класс от привязки к конкретным классам продуктов.
- Выделяет код производства продуктов в одно место, упрощая поддержку кода.
- Упрощает добавление новых продуктов в программу.
- Реализует принцип открытости/закрытости.

Недостатки
- Может привести к созданию больших параллельных иерархий классов, так как для каждого класса продукта надо создать свой подкласс создателя.

*/

// В этом примере создаются разные типы оружия при помощи структуры фабрики.
func main() {
	ak47, _ := getGun("ak47")
	musket, _ := getGun("musket")

	printDetails(ak47)
	printDetails(musket)
}

func printDetails(g IGun) {
	fmt.Printf("Gun: %s", g.getName())
	fmt.Println()
	fmt.Printf("Power: %d", g.getPower())
	fmt.Println()
}

// интерфейс продукта
type IGun interface {
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
}

// конкретный продукт
type Gun struct {
	name  string
	power int
}

func (g *Gun) setName(name string) {
	g.name = name
}

func (g *Gun) getName() string {
	return g.name
}

func (g *Gun) setPower(power int) {
	g.power = power
}

func (g *Gun) getPower() int {
	return g.power
}

// конкретный продукт
type Ak47 struct {
	Gun
}

func newAk47() IGun {
	return &Ak47{
		Gun: Gun{
			name:  "AK47 gun",
			power: 4,
		},
	}
}

// конкретный продукт
type musket struct {
	Gun
}

func newMusket() IGun {
	return &musket{
		Gun: Gun{
			name:  "Musket gun",
			power: 1,
		},
	}
}

// фабрика
func getGun(gunType string) (IGun, error) {
	if gunType == "ak47" {
		return newAk47(), nil
	}
	if gunType == "musket" {
		return newMusket(), nil
	}
	return nil, fmt.Errorf("Wrong gun type passed")
}
