package gothon

import (
	"fmt"
	"os"
	"sync"

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

var create_new_img *python3.PyObject

// 初始化python环境
func init() {
	python3.Py_Initialize()
	if !python3.Py_IsInitialized() {
		fmt.Println("Error initializing the python interpreter")
		os.Exit(1)
	}
	d2p := ImportModule(".\\gothon", "d2p")
	create_new_img = d2p.GetAttrString("create_new_img")
}

// code from: https://blog.csdn.net/weixin_43988498/article/details/120631081
func CreateNewImg(wg *sync.WaitGroup, ch chan string, w Mblog, cookie string) {
	defer func() { ch <- w.Bid }()
	defer wg.Done()
	wg.Wait()
	wg.Add(1)

	// 设置调用的参数 args(post, userInfo, cookie)
	post := python3.PyDict_New()
	SetItemString(post, "mid", w.Mid)
	SetItemString(post, "bid", w.Bid)
	if w.RetweetedStatus == nil {
		python3.PyDict_SetItemString(post, "repo", python3.Py_None)
		SetItemString(post, "text", w.Text)
	} else {
		SetItemString(post, "repo", w.Text)
		SetItemString(post, "text", "转发了 @"+w.RetweetedStatus.User.ScreenName+"："+w.RetweetedStatus.Text)
	}
	picUrls := python3.PyList_New(len(w.PicIds))
	for pos, pic := range w.PicIds {
		python3.PyList_SetItem(picUrls, pos, python3.PyUnicode_FromString(pic))
	}
	python3.PyDict_SetItemString(post, "picUrls", picUrls)
	SetItemString(post, "time", w.CreatedAt)
	SetItemString(post, "source", w.Source)

	userInfo := python3.PyDict_New()
	python3.PyDict_SetItemString(userInfo, "id", python3.PyLong_FromGoInt(w.User.ID))
	SetItemString(userInfo, "name", w.User.ScreenName)
	SetItemString(userInfo, "face", w.User.AvatarHd)
	SetItemString(userInfo, "desc", w.User.Description)
	python3.PyDict_SetItemString(userInfo, "follow", python3.PyLong_FromGoInt(w.User.FollowCount))
	SetItemString(userInfo, "follower", w.User.FollowersCountStr)

	args := python3.PyTuple_New(3)
	python3.PyTuple_SetItem(args, 0, post)
	python3.PyTuple_SetItem(args, 1, userInfo)
	python3.PyTuple_SetItem(args, 2, python3.PyUnicode_FromString(cookie))

	save := create_new_img.Call(args, python3.Py_None)
	file := python3.PyTuple_New(1)
	python3.PyTuple_SetItem(file, 0, python3.PyUnicode_FromString(w.Bid+".png"))
	save.Call(file, python3.Py_None)
}

// ImportModule
// @Description: 导入一个包
func ImportModule(dir, name string) *python3.PyObject {
	sysModule := python3.PyImport_ImportModule("sys")                 // import sys
	path := sysModule.GetAttrString("path")                           // path = sys.path
	python3.PyList_Insert(path, 0, python3.PyUnicode_FromString(dir)) // path.insert(0, dir)
	return python3.PyImport_ImportModule(name)                        // return __import__(name)
}

// 为字典中添加项
func SetItemString(p *python3.PyObject, key, val string) {
	python3.PyDict_SetItemString(p, key, python3.PyUnicode_FromString(val))
}
