# weakmap

A generic cache map implementation with [weak
pointers](https://pkg.go.dev/weak#Pointer) as described [on the Go
blog](https://go.dev/blog/cleanups-and-weak#weak-pointers).

## Usage

To initialize a `Map` instance, provide a method for instantiating new values
given a key.

```go
m := Map[int, string]{
    New: func(key int) (*string, error) {
        value := strconv.Itoa(key)
        return &value, nil
    },
}

k := 1
v, err := m.Load(k) // *v == "1", err == nil
```

Instances are safe for concurrent use but should not be copied.
