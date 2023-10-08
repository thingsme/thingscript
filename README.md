[![codecov](https://codecov.io/gh/thingsme/thingscript/graph/badge.svg?token=9SCBJ4SXOF)](https://codecov.io/gh/thingsme/thingscript)

ThingScript is a script interpreter that implemented in Go without zero dependency.
It can be embedded in your Go application as a library.

## Data Types

### STRING

Double quoted text in unicode.

```go
var str1 = "hello"
str2 := "world"
str3 := str1 + " "+ str2
```

### BOOLEAN

`true` and `false`

```go
var b1 = true
b2 := false
b3 := 10 < 20 // true
```

### INTEGER

64 bit integer number

```go
var i1 = 123
zero := 0
```

### FLOAT

64 bit floating point number

```go
var f1 = 3.14
zero := 0.0
```

## Control flow

### IF-ELSE

```go
var n = 10
var m = 20

max := if n > m { n } else { m }
```

### Function

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

