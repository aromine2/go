package main

import "fmt"

// declare some type
type someType struct {
    someTypeProperty float64
}

// declare what you want functions that type has
type thingsForSomeTypeToDo interface{
    functionToDo() float64
}

// declare what each of those functions do
func (v someType) functionToDo() float64 {
    return v.someTypeProperty
}

// declare a function to use your type and the functions on that type
func overArchingFunction(v thingsForSomeTypeToDo) {
    fmt.Println(v.functionToDo())
}

func main() {
	// create an instance of your type you created
    v := someType{1}

    // pass the type into the regular function
    overArchingFunction(v)
}