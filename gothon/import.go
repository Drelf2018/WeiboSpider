package main

// code from: https://blog.csdn.net/weixin_43988498/article/details/120631081
import (
	"fmt"
	"os"

	python3 "github.com/go-python/cpy3"
)

func init() {
	// 初始化python环境
	python3.Py_Initialize()
	if !python3.Py_IsInitialized() {
		fmt.Println("Error initializing the python interpreter")
		os.Exit(1)
	}
}

func main() {
	// 提前结束环境，后文出现 defer 多为保护内存安全
	defer python3.Py_Finalize()

	InsertBeforeSysPath("C:\\Users\\drelf\\AppData\\Local\\Programs\\Python\\Python37\\Lib\\site-packages")

	// 导入 hello 模块
	hello := ImportModule(".\\gothon", "hello")
	defer hello.DecRef()

	// pyObject => string 解析结果
	helloRepr, err := pythonRepr(hello)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[MODULE] repr(hello) = %s\n", helloRepr)
	// 4. 获取变量
	a := hello.GetAttrString("a")
	aString, err := pythonRepr(a)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[VARS] a = %#v\n", aString)
	// 5. 获取函数方法
	SayHello := hello.GetAttrString("SayHello")
	// 设置调用的参数（一个元组）
	args := python3.PyTuple_New(1)                                        // 创建存储空间
	python3.PyTuple_SetItem(args, 0, python3.PyUnicode_FromString("xwj")) // 设置值
	res := SayHello.Call(args, python3.Py_None)                           // 调用
	fmt.Printf("[FUNC] res = %s\n", python3.PyUnicode_AsUTF8(res))
	// 6. 调用第三方库sklearn
	sklearn := hello.GetAttrString("aiohttp")
	skVersion := sklearn.GetAttrString("__version__")
	sklearnRepr, err := pythonRepr(sklearn)
	if err != nil {
		panic(err)
	}
	skVersionRepr, err := pythonRepr(skVersion)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[IMPORT] sklearn = %s\n", sklearnRepr)
	fmt.Printf("[IMPORT] sklearn version =  %s\n", skVersionRepr)
}

func InsertBeforeSysPath(p string) {
	sysModule := python3.PyImport_ImportModule("sys")
	path := sysModule.GetAttrString("path")
	pObject := python3.PyUnicode_FromString(p)
	defer pObject.DecRef()
	python3.PyList_Append(path, pObject)
}

// ImportModule
// @Description: 导入一个包
func ImportModule(dir, name string) *python3.PyObject {
	sysModule := python3.PyImport_ImportModule("sys")                 // import sys
	path := sysModule.GetAttrString("path")                           // path = sys.path
	python3.PyList_Insert(path, 0, python3.PyUnicode_FromString(dir)) // path.insert(0, dir)
	return python3.PyImport_ImportModule(name)                        // return __import__(name)
}

// pythonRepr
// @Description: PyObject转换为string
func pythonRepr(o *python3.PyObject) (string, error) {
	if o == nil {
		return "", fmt.Errorf("object is nil")
	}
	s := o.Repr()
	if s == nil {
		python3.PyErr_Clear()
		return "", fmt.Errorf("failed to call Repr object method")
	}
	defer s.DecRef()

	return python3.PyUnicode_AsUTF8(s), nil
}
