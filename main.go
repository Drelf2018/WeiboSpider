package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type LogFormatter struct{}

// code from: https://zhuanlan.zhihu.com/p/538240645
func (m *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("15:04:05")
	msg := fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)

	b.WriteString(msg)
	return b.Bytes(), nil
}

// 微博请求头
var req *http.Request

// 用于对比微博是否编辑的集合 key: bid val: content
var MbText = make(map[string]string)

// 配置文件
var cfg Config

// 避免同时调用 python 库的等待组
var wg sync.WaitGroup

// 为了解决等待组中 goroutine 超过 2 引发恐慌(panic)而引入的计数器
// 虽然它的存在让代码更臃肿了 但是本人 golang 苦手 只能这样操作
var Tasks int

// 初始化请求会话并在请求头中加入 cookie
func init() {
	LoadContent()
	LoadConfig()
	log.SetFormatter(&LogFormatter{})
	url := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?containerid=107603%v", cfg.Account.UID)
	var err error
	if req, err = http.NewRequest("GET", url, nil); err == nil {
		req.Header.Add("cookie", cfg.Account.Cookie)
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	} else {
		panic(err)
	}
	if CommentReq, err = http.NewRequest("GET", CommentUrl, nil); err == nil {
		req.Header.Add("cookie", cfg.Account.Cookie)
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

// 这个人的微博也查(query)一下吧
func Query() {
	log.Info("开始查询")
	index := GetIndex()
	for _, card := range index.Data.Cards {
		switch card.CardType {
		// 发微博或评论
		case 9:
			go GetComment(card.Mblog.Bid, card.Mblog.Mid)
			// 如果 未记录过该微博 或者 内容有编辑 的同时满足 生成图片任务数小于 2 时
			if SavedText, ok := MbText[card.Mblog.Bid]; (!ok || SavedText != card.Mblog.Text) && Tasks < 2 {
				log.Info("更新微博 " + card.Mblog.Bid + " " + card.Mblog.Text)
				ch := make(chan string)
				Tasks += 1

				// 生成图片 生成完毕通过 ch 返回该任务微博bid
				go gothon.CreateNewImg(&wg, ch, card.Mblog, cfg.Account.Cookie)

				// 更新此微博的内容
				go func(ch chan string, NewText string) {
					bid := <-ch
					MbText[bid] = NewText
					go SaveContent()
					log.Info("完成更新 " + bid)
					Tasks -= 1
					// send_msg
				}(ch, card.Mblog.Text)
			}
		}
	}
}

func main() {
	// 间隔 10 秒的定时任务器
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	OnlineStatus := Online()
	// 提前执行一次任务 急 等不了了
	Query()
	for range ticker.C {
		go Query()
		go OnlineStatus()
	}

}
