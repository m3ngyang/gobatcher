gobatcher
===

`gobatcher` is a light-weighted golang library for batch, which operates the function on the given parameter list concurrently. It fails once any execution returns error.

## How to use
Here comes an easy example.

```golang
// self-defined function must return error in the result list
func echoString(str string) error {
    fmt.Println(str)
    return nil
}

strs := []string{"Hello", "World", "!"}
goBatcher := New(echoString, strs, 2)
goBatcher.Run()
```

More examples are in the test file.