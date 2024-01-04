Что выведет программа? Объяснить вывод программы.
```go
package main
 
type customError struct {
     msg string
}
 
func (e *customError) Error() string {
    return e.msg
}
 
func test() *customError {
     {
         // do something
     }
     return nil
}
 
func main() {
    var err error
    err = test()
    if err != nil {
        println("error")
        return
    }
    println("ok")
}

```
Ответ:
```
error
```
- Объявление переменной `err` типа интерфейса `error`
- Присвоение конкретного типа `*customError` интерфейсной переменной, но значение интерфейсного типа `nil`.
Т.е. у переменной `err` поле описывающее тип `tab *itab` теперь не нулевое, но поле описывающее данные `data unsafe.Pointer` все еще `nil`
  - Это можно проверить `fmt.Printf("type: %T\tvalue: %v\n", err, err)` -> `type: *main.customError value: <nil>`
- И так как значение интерфейсного типа == `nil` только тогда, когда у него нет, ни конкретного типа, ни значения конкретного типа, 
а в данном случает конкретный тип у нас есть, поэтому логика проваливается в `if err != nil {...}` и выводит "error".
