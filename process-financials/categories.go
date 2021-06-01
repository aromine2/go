package main

import (
    "fmt"
    "sort"
)

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

func printCategoryResults(category category) {
    values := category.values
    name := category.name

    // make new slice of strings so keys can be sorted
    keys := make([]string, 0, len(values))
    for key := range category.values {
        keys = append(keys, key)
    }
    sort.Strings(keys)

    fullTotal := 0.00
    // print sorted values from above step
    for _, k := range keys {
        fmt.Printf("%s %s: %.2f\n", name, k, values[k])
        fullTotal = fullTotal + values[k]
    }
    fmt.Printf("%s full total: %.2f\n", name, fullTotal)
}

