[![codecov](https://codecov.io/gh/thingsme/thingscript/graph/badge.svg?token=9SCBJ4SXOF)](https://codecov.io/gh/thingsme/thingscript)

ThingScript is a script interpreter that implemented in Go without zero dependency.
It can be embedded in your Go application as a library.

## Data Types

### STRING

### BOOLEAN

### INTEGER

### FLOAT

## Control flow

### IF-ELSE

```go
var n = 10
var m = 20

max := if n > m { n } else { m }
```

### Functions

```go

func inc(x) {
    x + 1 // 'return' is optional
}

func dec(x) {
    return x - 1
}

eleven := inc(10)
nine := dec(10)
```

