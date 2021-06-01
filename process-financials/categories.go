package main

// use the function, not the struct directly
type category struct {
name string
filters []string
values map[string]float64
}

func newCategory(name string, filters []string) category {
    return category{
      name: name,
      filters: filters,
      values: make(map[string]float64),
    }
}
