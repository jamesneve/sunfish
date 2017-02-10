# Sunfish

![Image of Sunfish](sunfish.jpg)

A CSV parser for Go.

## Usage

With a CSV file (for example, recipes.csv below):

```
user_id,recipe
33,How to boil an egg: Put it in hot water and wait for a while
45,How to cook a sunfish: You should never cook a sunfish, you monster
```

Parse the file into the struct by creating a reader and reading from a file. Only rows marked by `csv:"parse"` will be parsed.
The others can be set manually afterwards.
```go
package recipe

import "github.com/jamesneve/sunfish/parser"

type Recipe struct {
        UserID int `csv:"user_id"`
        Text string `csv:"recipe"`
        Evaluation bool
}

func ReadRecipes() ([]Recipe, error) {
        var recipes []Recipe
        
        r := parser.NewParser()
        err := r.ReadCsvFromFileWithHeaders("./recipes.csv", &recipes)
        
        return recipes, err
}
```

For a CSV without headers e.g.

```
33,How to boil an egg: Put it in hot water and wait for a while
45,How to cook a sunfish: You should never cook a sunfish, you monster
```

Define your struct with the fields you want to parse labelled `parse`

```go
type Recipe struct {
        UserID int `csv:"parse"`
        Text string `csv:"parse"`
        Evaluation bool
}
```

And call `ReadCsvFromFileInOrder`

If you don't want to read from a file, but have an io.Reader instance from somewhere else, you can also call `ReadCsvInOrder`, 
and `ReadCsvWithHeaders` passing the reader directly.

## Future development

* Add support for more types including arrays, pointers and nested structs (currently only primitives)
* Better error handling

# License

Apache License, Version 2.0