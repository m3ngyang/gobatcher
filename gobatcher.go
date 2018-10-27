package gobatcher

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	errArgNotFunc    = errors.New("Function Arg is Not a function type")
	errArgNotSlice   = errors.New("Parameters Arg is Not a slice type")
	errArgNotPostive = errors.New("Concurrency Arg is Not a positive int")
)

// GoBatcher is a batcher that operates the same func on given arg set
type GoBatcher struct {
	funct        interface{}   // called function
	vars         []interface{} // arg set
	maxConcurNum int           // concurrency
}

// New creates a new batcher object
func New(f interface{}, v interface{}, m int) *GoBatcher {
	ft := reflect.ValueOf(f)
	fType := ft.Type()
	if fType.Kind() != reflect.Func {
		panic(errArgNotFunc)
	}

	slice := reflect.ValueOf(v)
	if slice.Kind() != reflect.Slice {
		panic(errArgNotSlice)
	}

	if m < 1 {
		panic(errArgNotPostive)
	}

	args := make([]interface{}, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		args[i] = slice.Index(i).Interface()
	}

	gb := new(GoBatcher)
	gb.maxConcurNum = m
	gb.funct = f
	gb.vars = args
	return gb
}

// Run executes the batcher
func (gb *GoBatcher) Run() error {
	lock := sync.Mutex{}
	var row []interface{}
	rows := [][]interface{}{}

	for idx, i := range gb.vars {
		if idx%gb.maxConcurNum == 0 {
			row = []interface{}{}
		}
		row = append(row, i)
		if len(row) == gb.maxConcurNum || idx == len(gb.vars)-1 {
			rows = append(rows, row)
		}
	}

	for i := 0; i < len(rows); i++ {
		succ := 0
		errCh := make(chan error, 1)
		for _, element := range rows[i] {
			// fmt.Println("row:", i, ", element:", element)
			go func(e interface{}) {
				funct := reflect.ValueOf(gb.funct)
				paras := make([]reflect.Value, 1)
				if e == nil {
					paras[0] = reflect.Zero(funct.Type().In(0))
				} else {
					paras[0] = reflect.ValueOf(e)
				}
				out := funct.Call(paras)
				retErr := out[len(out)-1].Interface()

				if retErr != nil {
					errInfo := new(bytes.Buffer)
					fmt.Fprint(errInfo, retErr)
					errCh <- errors.New(errInfo.String())
					return
				}

				lock.Lock()
				succ++
				lock.Unlock()
				if succ == len(rows[i]) {
					errCh <- nil
				}
			}(element)
		}

		happenedErr := <-errCh
		if happenedErr != nil {
			return happenedErr
		}
		continue
	}
	return nil
}
