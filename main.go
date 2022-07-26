package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

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

// 微博请求头与 cookie 信息
var req *http.Request
var cookie = "SCF=Ao_CfNDAVoOOnPeRbWYncYwIxkt1qdMnP8GJwYAFUbdV2AwNmv18NfGuAVLGkp_al91dOQIJTvviNKRX_WnQV0I.; SUB=_2A25P0FJVDeRhGeFM7lQU8i7PyjmIHXVtO34drDV6PUJbktCOLRHBkW1NQNxRTyundxN1rIEaDb8vuDu9LX1y2PfV; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WWsM7qnLH5XXeUsRC8WX5b75NHD95QNeo-cSKz7e02fWs4DqcjPi--RiKnXiKnci--4i-zEi-2ReKzpe0nt; _T_WM=b7e06bb363bcddabb34f5e7aba1d13f8; WEIBOCN_FROM=1110006030; MLOGIN=1; XSRF-TOKEN=56db96; mweibo_short_token=b93d1a8e49; M_WEIBOCN_PARAMS=from=page_100808&luicode=10000011&lfid=1008081a127e1db26d4483eadf1d1dbe1a80c2_-_live&fid=1008081a127e1db26d4483eadf1d1dbe1a80c2_-_live&uicode=10000011"

// 用于对比微博是否编辑的集合 key: bid val: content
var MbText = make(map[string]string)

// 避免同时调用 python 库的等待组
var wg sync.WaitGroup

// 为了解决等待组中 goroutine 超过 2 引发恐慌(panic) 而引入的计数器
// 虽然它的存在让代码更臃肿了 但是本人 golang 苦手 只能这样操作
var Tasks int

// 初始化请求会话并在请求头中加入 cookie
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

// 访问微博 api 获取新微博、评论
// 吗的 开了梯子这里就爬不动了
// 我还以为代码哪里有问题了
func GetIndex() (res gothon.Result) {
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

// 这个人的微博也查(query) 一下吧
func Query() {
	index := GetIndex()
	for _, card := range index.Data.Cards {
		if c := card.CardGroup; c != nil {
			switch c[1].CardType {
			// 发微博或评论
			case "9":
				if strings.Contains(c[0].Desc, "评论") {
					// 找评论去了
				} else {
					// 如果 未记录过该微博 或者 内容有编辑 的同时满足 生成图片任务数小于 2 时
					if SavedText, ok := MbText[c[1].Mblog.Bid]; (!ok || SavedText != c[1].Mblog.Text) && Tasks < 2 {
						fmt.Println(c[1].Mblog.Bid, "start")
						ch := make(chan string)
						Tasks += 1

						// 生成图片 生成完毕通过 ch 返回该任务微博bid
						go gothon.CreateNewImg(&wg, ch, c[1].Mblog, cookie)

						// 更新此微博的内容
						go func(ch chan string, NewText string) {
							MbText[<-ch] = NewText
							Tasks -= 1
							// send_msg
						}(ch, c[1].Mblog.Text)
					}
				}
			case "30": // 在线状态
				fmt.Println(c[1].Desc1)
			}
		}
	}
}

func main() {
	// 间隔 10 秒的定时任务器
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// 提前执行一次任务 急 等不了了
	Query()
	for range ticker.C {
		Query()
	}

}
