package gothon

// code from: https://blog.csdn.net/weixin_43988498/article/details/120631081
import (
	"fmt"
	"os"

	python3 "github.com/go-python/cpy3"
)

type Mblog struct {
	Bid        string   `json:"bid"`
	CreatedAt  string   `json:"created_at"`
	IsLongText bool     `json:"isLongText"`
	Mid        string   `json:"mid"`
	PicIds     []string `json:"pic_ids"`
	Source     string   `json:"source"`
	Text       string   `json:"text"`
	User       struct {
		AvatarHd          string `json:"avatar_hd"`
		Description       string `json:"description"`
		FollowCount       int    `json:"follow_count"`
		FollowersCountStr string `json:"followers_count_str"`
		ID                int    `json:"id"`
		ScreenName        string `json:"screen_name"`
	} `json:"user"`
	RetweetedStatus *Mblog `json:"retweeted_status"`
}

type Result struct {
	Data struct {
		Cards []struct {
			CardGroup []struct {
				CardType string `json:"card_type"`
				Desc     string `json:"desc,omitempty"`
				Desc1    string `json:"desc1,omitempty"`
				Mblog    Mblog  `json:"mblog,omitempty"`
			} `json:"card_group,omitempty"`
		} `json:"cards"`
	} `json:"data"`
}

var d2p, create_new_img *python3.PyObject

func init() {
	// 初始化python环境
	python3.Py_Initialize()
	if !python3.Py_IsInitialized() {
		fmt.Println("Error initializing the python interpreter")
		os.Exit(1)
	}
	//InsertBeforeSysPath("C:\\Users\\drelf\\AppData\\Local\\Programs\\Python\\Python37\\Lib\\site-packages")
	d2p = ImportModule(".\\gothon", "d2p")
	create_new_img = d2p.GetAttrString("create_new_img")
}

func CreateNewImg(w Mblog, cookie string) {
	// 获取函数方法
	//userInfo := d2p.GetAttrString("userInfo")
	//headers := d2p.GetAttrString("headers")

	// 设置调用的参数（一个元组）
	post := python3.PyDict_New()
	python3.PyDict_SetItemString(post, "mid", python3.PyUnicode_FromString(w.Mid))
	if w.RetweetedStatus == nil {
		python3.PyDict_SetItemString(post, "repo", python3.Py_None)
		python3.PyDict_SetItemString(post, "text", python3.PyUnicode_FromString(w.Text))
	} else {
		python3.PyDict_SetItemString(post, "repo", python3.PyUnicode_FromString(w.Text))
		python3.PyDict_SetItemString(post, "text", python3.PyUnicode_FromString("转发了 @"+w.RetweetedStatus.User.ScreenName+"："+w.RetweetedStatus.Text))
	}
	picUrls := python3.PyList_New(len(w.PicIds))
	for pos, pic := range w.PicIds {
		python3.PyList_SetItem(picUrls, pos, python3.PyUnicode_FromString(pic))
	}
	python3.PyDict_SetItemString(post, "picUrls", picUrls)
	python3.PyDict_SetItemString(post, "time", python3.PyUnicode_FromString(w.CreatedAt))
	python3.PyDict_SetItemString(post, "source", python3.PyUnicode_FromString(w.Source))

	userInfo := python3.PyDict_New()
	python3.PyDict_SetItemString(userInfo, "id", python3.PyLong_FromGoInt(w.User.ID))
	python3.PyDict_SetItemString(userInfo, "name", python3.PyUnicode_FromString(w.User.ScreenName))
	python3.PyDict_SetItemString(userInfo, "face", python3.PyUnicode_FromString(w.User.AvatarHd))
	python3.PyDict_SetItemString(userInfo, "desc", python3.PyUnicode_FromString(w.User.Description))
	python3.PyDict_SetItemString(userInfo, "follow", python3.PyLong_FromGoInt(w.User.FollowCount))
	python3.PyDict_SetItemString(userInfo, "follower", python3.PyUnicode_FromString(w.User.FollowersCountStr))

	args := python3.PyTuple_New(3)
	python3.PyTuple_SetItem(args, 0, post)
	python3.PyTuple_SetItem(args, 1, userInfo)
	python3.PyTuple_SetItem(args, 2, python3.PyUnicode_FromString(cookie))

	save := create_new_img.Call(args, python3.Py_None)
	img, _ := pythonRepr(save)
	fmt.Println(img)
	file := python3.PyTuple_New(1)
	python3.PyTuple_SetItem(file, 0, python3.PyUnicode_FromString(w.Bid+".png"))
	save.Call(file, python3.Py_None)
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
