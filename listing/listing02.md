Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
    "fmt"
)

func test() (x int) {
    defer func() {
        x++
    }()
    x = 1
    return
}

func anotherTest() int {
    var x int
    defer func() {
        x++
    }()
    x = 1
    return x
}

func main() {
    fmt.Println(test())
    fmt.Println(anotherTest())
}
```

Ответ:
```
2
1
```
`defer` откладывает выполнение инструкций в `func(...){}()` до завершения функции
(после выставления возвращаемого значения, но до передачи управления вызывающей функции),
при этом аргументы вычисляются немедленно в момент объявления `defer`.

Алгоритм работы defer:
1. Вызов defer и оценка параметров функции
2. Выполнение defer -> добавление функции в стек
3. Выполнение функций из стека (LIFO) перед возвратом управления (return или panic)

`fmt.Println(anotherTest()) -> 1`:
При вызове `defer` переменные оцениваются немедленно: компилятор ищет кандидата на добавление 
в замыкание, просматривая код вверх по тексту, до первого символа с именем `x`, видимого в 
точке определения отложенной функции. И встретит `var x int` - это внутренняя переменная, анонимная 
функция в defer ее изменит, но это произойдет уже после того как ее значение установится как возвращаемое.
Поэтому на результат это не повлияет.

`fmt.Println(test()) -> 2`:
В данном случае компилятор свяжет `x` внутри отложенной функции с возвращаемым значением.
Поэтому результат вычисления `x++` будет записан в ту область стека где хранится возвращаемое значение.









