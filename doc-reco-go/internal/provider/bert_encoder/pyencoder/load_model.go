package pyencoder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DataDog/go-python3"
)

var (
	oImport, oModule, oDict, args, loadModelPyFunc, loadModelRes, encodePyFunc *python3.PyObject
)

var state *python3.PyThreadState

func LoadModel(modelName string) error {

	python3.Py_InitializeEx(false)
	if !python3.Py_IsInitialized() {
		return fmt.Errorf("error initializing the python interpreter")
	}

	runtime.LockOSThread()

	state = python3.PyEval_SaveThread()

	gstate := python3.PyGILState_Ensure()

	err := loadModel(modelName)

	python3.PyGILState_Release(gstate)

	runtime.UnlockOSThread()
	return err
}

func loadModel(modelName string) error {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = filepath.Join(dir, "/internal/provider/bert_encoder/pyencoder")
	if err != nil {
		return err
	}

	ret := python3.PyRun_SimpleString("import sys\nsys.path.append(\"" + dir + "\")")

	if ret != 0 {
		log.Fatalf("error appending '%s' to python sys.path", dir)
	}

	oImport = python3.PyImport_ImportModule("pyencoder") //ret val: new ref
	if !(oImport != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return fmt.Errorf("failed to import module 'pyencoder'")
	}

	oModule = python3.PyImport_AddModule("pyencoder") //ret val: borrowed ref (from oImport)
	if !(oModule != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return fmt.Errorf("failed to add module 'pyencoder'")
	}

	oDict = python3.PyModule_GetDict(oModule) //ret val: Borrowed
	if !(oDict != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return fmt.Errorf("could not get dict for module")
	}

	///////////////////////////// LOADING MODEL /////////////////////////
	args = python3.PyTuple_New(1)

	ret = python3.PyTuple_SetItem(args, 0, python3.PyBytes_FromString(modelName))
	if ret != 0 {
		if python3.PyErr_Occurred() != nil {
			python3.PyErr_Print()
		}
	}
	loadModelPyFunc = python3.PyDict_GetItemString(oDict, "init_model")

	loadModelRes = loadModelPyFunc.CallObject(args)
	if !(loadModelRes != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return fmt.Errorf("error calling py function init_model")
	}

	encodePyFunc = python3.PyDict_GetItemString(oDict, "encode")
	return nil
}

func DeferPyEncoder() {
	loadModelRes.DecRef()
	loadModelPyFunc.DecRef()
	args.DecRef()
	oDict.DecRef()
	oModule.DecRef()
	oImport.DecRef()

	python3.PyEval_RestoreThread(state)
	python3.Py_Finalize()
}
