package pyencoder

import (
	"fmt"
	"runtime"

	"github.com/DataDog/go-python3"
)

func PyEncode(query string) ([]float32, error) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	gstate := python3.PyGILState_Ensure()
	defer python3.PyGILState_Release(gstate)

	goList, err := pyEncode(query)

	return goList, err
}

func pyEncode(query string) ([]float32, error) {

	queryBytes := python3.PyBytes_FromString(query)
	defer queryBytes.DecRef()

	encodeFuncRes := encodePyFunc.CallFunctionObjArgs(queryBytes)

	if !(encodeFuncRes != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return nil, fmt.Errorf("error calling function encode")
	}
	defer encodeFuncRes.DecRef()

	golist, err := goSliceFromPylist(encodeFuncRes)
	return golist, err
}

func goSliceFromPylist(pylist *python3.PyObject) ([]float32, error) {

	seq := pylist.GetIter() //ret val: New reference
	if !(seq != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return nil, fmt.Errorf("error creating iterator for list")
	}
	defer seq.DecRef()
	tNext := seq.GetAttrString("__next__") //ret val: new ref
	if !(tNext != nil && python3.PyCallable_Check(tNext)) {
		return nil, fmt.Errorf("iterator has no __next__ function")
	}
	defer tNext.DecRef()

	var golist []float32

	pylistLen := pylist.Length()
	if pylistLen == -1 {
		return nil, fmt.Errorf("error getting list length")
	}

	var itemType *python3.PyObject
	for i := 1; i <= pylistLen; i++ {
		item := tNext.CallObject(nil) //ret val: new ref
		if item == nil && python3.PyErr_Occurred() != nil {
			python3.PyErr_Print()
			return nil, fmt.Errorf("error getting next item in sequence")
		}
		itemType = item.Type()
		if itemType == nil && python3.PyErr_Occurred() != nil {
			python3.PyErr_Print()
			return nil, fmt.Errorf("error getting item type")
		}

		itemGo := python3.PyFloat_AsDouble(item)
		if itemGo != -1 && python3.PyErr_Occurred() == nil {
			golist = append(golist, float32(itemGo))
		} else {
			return golist, fmt.Errorf("error parsing list element As Double")
		}

		if item != nil {
			item.DecRef()
			item = nil
		}
	}
	defer itemType.DecRef()

	return golist, nil
}
