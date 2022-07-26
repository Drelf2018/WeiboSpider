package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Drelf2018/WeiboSpider/gothon"
	log "github.com/sirupsen/logrus"
)

func checkErr(err error) bool {
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

var req *http.Request
var cookie = ""

func init() {
	url := "https://m.weibo.cn/api/container/getIndex?from=page_100808&mod[]=TAB%3Ffrom%3Dpage_100808&mod[]=TAB&containerid=1008081a127e1db26d4483eadf1d1dbe1a80c2_-_live"
	var err error
	if req, err = http.NewRequest("GET", url, nil); err == nil {
		req.Header.Add("cookie", cookie)
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	} else {
		panic(err)
	}
}

func getIndex() (res gothon.Result) {
	if resp, err := (&http.Client{}).Do(req); checkErr(err) {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); checkErr(err) {
			if err := json.Unmarshal(body, &res); checkErr(err) {
				return
			}
		}
	}
	return
}

func query() {
	index := getIndex()
	// var res gothon.Result
	// json.Unmarshal(js, &res)
	for _, card := range index.Data.Cards {
		if c := card.CardGroup; c != nil {
			switch c[1].CardType {
			case "9": // 发微博或评论
				if strings.Contains(c[0].Desc, "评论") {
					// 找评论去了
				} else {
					gothon.CreateNewImg(c[1].Mblog, cookie)
				}
			case "30": // 在线状态
				fmt.Println(c[1].Desc1)
			}
		}
	}
}

func main() {
	query()
}
