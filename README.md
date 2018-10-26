# gobatcher
[![Build Status](https://travis-ci.org/m3ngyang/gobatcher.svg?branch=master)](https://travis-ci.org/m3ngyang/gobatcher)
[![Coverage Status](https://coveralls.io/repos/github/m3ngyang/gobatcher/badge.svg?branch=master)](https://coveralls.io/github/m3ngyang/gobatcher?branch=master)

`gobatcher` is a light-weighted golang library for batch, which operates the function on the given parameter list concurrently. It fails once any execution returns error.

## Install
```shell
go get -u github.com/m3ngyang/gobatcher
```

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