# Sunfish

![Image of Sunfish](sunfish.jpg)

A CSV parser for Go.

## Usage

With a CSV file (for example, recipes.csv below):

```
1,How to boil an egg: Put it in hot water and wait for a while
2,How to cook a sunfish: You should never cook a sunfish, you monster
```

Parse the file into the struct by creating a reader and reading from a file. Only rows marked by `csv:"parse"` will be parsed.
The others can be set manually afterwards.
```go
package recipe

import "github.com/jamesneve/sunfish"

type Recipe struct {
        UserID int `csv:"parse"`
        Text string `csv:"parse"`
        Evaluation bool
}

func ReadRecipes() ([]Recipe, error) {
        var recipes []Recipe
        
        r := sunfish.NewReader()
        err := r.ReadCsvFromFile("./recipes.csv", &recipes)
        
        return recipes, err
}
```

You can also use `r.ReadCsv` directly if you already have an `io.Reader` instance.

## Future development

* Parse files by column header based on the `csv` parameter, not just by order
* Add support for more types including arrays, pointers and nested structs (currently only primitives)
* Better error handling

# License

Apache License, Version 2.0