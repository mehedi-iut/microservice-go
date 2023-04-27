In Go, an interface is a collection of method signatures. An interface defines a set of methods that a struct or other type must implement in order to be considered as "implementing" the interface. Interfaces are used to define a contract for the behavior of a type, and to allow for polymorphism and decoupling of code.

Here's an example of an interface definition:
```go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

This interface defines two methods, Area and Perimeter, which return a float64 value. Any struct or other type that implements these two methods can be considered as implementing the Shape interface.

Here's an example of a struct that implements the Shape interface:
```go
type Rectangle struct {
    width, height float64
}

func (r Rectangle) Area() float64 {
    return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
    return 2*(r.width + r.height)
}
```

In this example, the Rectangle struct has implemented the Area and Perimeter methods defined in the Shape interface.

An interface type can be used wherever a struct or other type is expected, because any struct or other type that implements the interface is considered "assignable" to the interface type.
```go
var s Shape
s = Rectangle{5.0, 4.0}
```

In this example, the variable s has the type Shape, but it is assigned a value of type Rectangle. This is possible because Rectangle implements the Shape interface. This enables you to write functions that can work with any type that implements the interface, without knowing the specific implementation.
```go
func Measure(s Shape) {
    fmt.Println("Area:", s.Area())
    fmt.Println("Perimeter:", s.Perimeter())
}
```

In this example, the Measure function takes a parameter of type Shape, and it can work with any struct or other type that implements the Shape interface, regardless of the specific implementation.

Interfaces also provide a way for structs to achieve polymorphism, which means that a single function can work with multiple types. Go does not have inheritance but interfaces do provide a way to achieve polymorphism in Go.