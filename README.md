[![codecov](https://codecov.io/gh/thingsme/thingscript/graph/badge.svg?token=9SCBJ4SXOF)](https://codecov.io/gh/thingsme/thingscript)

ThingScript is a script interpreter that implemented in Go with zero dependency.
It can be embedded in your Go application as a library, and install as an executable binary.

## Hello World

```go
out := import("fmt")
name := "World"
out.println("Hello", name, "?")
```

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

### ARRAY

```go
var arr = [ 0, 1, true, "hello world"]
```

### MAP

```go
var tbl = { "key1": 0, "key2": 1, "key3": true, "key4": "hello world"}
```

## Control Flow

### IF-ELSE

```go
var n = 10
var m = 20

max := if n > m { n } else { m }
```

### WHILE

```go
var n = 0
var sum = 0
while n < 10 {
    sum += n
    n += 1
}
```

### DO-WHILE

```go
var n = 0
var sum = 0
do {
    sum += n
    n += 1
} while n < 10;
```

### FOREACH

```go
sum := 0
[1,2,3].foreach(func(idx,elm){
    sum += elm
})
// sum = 6
```

```go
sum := ""
["1","2","3"].foreach(func(idx,elm){
    sum += elm
})
// sum = "123"
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

