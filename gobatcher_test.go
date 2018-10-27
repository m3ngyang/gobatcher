package gobatcher

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func echoString(str string) error {
	fmt.Println(str)
	return nil
}

func addOne(num int) (int, error) {
	return num + 1, nil
}

func addOneErr(num int) (int, error) {
	return num + 1, errors.New("testerror")
}

type TmpObj struct {
	innerStr string
}

func (t *TmpObj) printStr(str string) error {
	fmt.Println(t.innerStr, str)
	return nil
}

func TestRunEchoString(t *testing.T) {
	strs := []string{"Hello", "World", "!"}
	goBatcher := New(echoString, strs, 2)
	err := goBatcher.Run()
	assert.NoError(t, err, "Err should be nil")
}

func TestRunAddOne(t *testing.T) {
	nums := []int{1, 2, 3}
	goBatcher := New(addOneErr, nums, 2)
	err := goBatcher.Run()
	assert.EqualError(t, err, "testerror", "Err should be testerror")
}

func TestRunObjFunc(t *testing.T) {
	strs := []string{"Call", "Object", "Function"}
	to := new(TmpObj)
	to.innerStr = "Inner"
	goBatcher := New(to.printStr, strs, 2)
	err := goBatcher.Run()
	assert.NoError(t, err, "Err should be nil")
}

func TestPanic(t *testing.T) {
	function := "func"
	params := []string{"Hello", "World", "!"}
	assert.PanicsWithValue(t, errArgNotFunc, func() { New(function, params, 1) })
	assert.PanicsWithValue(t, errArgNotSlice, func() { New(echoString, params[0], 1) })
	assert.PanicsWithValue(t, errArgNotPostive, func() { New(echoString, params, -1) })
}

func TestDefaultParam(t *testing.T) {
	nums := []interface{}{nil}
	goBatcher := New(addOne, nums, 2)
	err := goBatcher.Run()
	assert.Nil(t, err)
}
