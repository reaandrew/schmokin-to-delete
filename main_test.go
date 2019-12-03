package main_test

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
)

func method() {
	fmt.Println("Boom")
}

func MyObject(call goja.ConstructorCall) *goja.Object {
	// call.This contains the newly created object as per http://www.ecma-international.org/ecma-262/5.1/index.html#sec-13.2.2
	// call.Arguments contain arguments passed to the function

	call.This.Set("method", method)

	//...

	// If return value is a non-nil *Object, it will be used instead of call.This
	// This way it is possible to return a Go struct or a map converted
	// into goja.Value using runtime.ToValue(), however in this case
	// instanceof will not work as expected.
	return nil
}

func TestDoesSomething(t *testing.T) {
	vm := goja.New()

	vm.Set("MyObject", MyObject)
	vm.RunString(`
new MyObject().method();
`)

}
