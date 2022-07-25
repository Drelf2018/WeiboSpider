package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

func init() {
	url := "https://m.weibo.cn/api/container/getIndex?from=page_100808&mod[]=TAB%3Ffrom%3Dpage_100808&mod[]=TAB&containerid=1008081a127e1db26d4483eadf1d1dbe1a80c2_-_live"
	var err error
	if req, err = http.NewRequest("GET", url, nil); err == nil {
		req.Header.Add("cookie", "SCF=Ao_CfNDAVoOOnPeRbWYncYwIxkt1qdMnP8GJwYAFUbdV2AwNmv18NfGuAVLGkp_al91dOQIJTvviNKRX_WnQV0I.; SUB=_2A25P0FJVDeRhGeFM7lQU8i7PyjmIHXVtO34drDV6PUJbktCOLRHBkW1NQNxRTyundxN1rIEaDb8vuDu9LX1y2PfV; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WWsM7qnLH5XXeUsRC8WX5b75NHD95QNeo-cSKz7e02fWs4DqcjPi--RiKnXiKnci--4i-zEi-2ReKzpe0nt; _T_WM=b7e06bb363bcddabb34f5e7aba1d13f8; WEIBOCN_FROM=1110006030; MLOGIN=1; XSRF-TOKEN=56db96; mweibo_short_token=b93d1a8e49; M_WEIBOCN_PARAMS=from=page_100808&luicode=10000011&lfid=1008081a127e1db26d4483eadf1d1dbe1a80c2_-_live&fid=1008081a127e1db26d4483eadf1d1dbe1a80c2_-_live&uicode=10000011")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	} else {
		panic(err)
	}
}

func getIndex() (res Result) {
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

func get_created_time(ctime string) string {
	return ctime
}

func handle(w Cardgroup) {
	retweeted_status := w.Mblog.RetweetedStatus
	is_long := w.Mblog.IsLongText
	weibo_id := w.Mblog.ID
	mid := w.Mblog.Mid
	// 	# 获取用户简介
	user_desc := w.Mblog.User.Description

	weibo_url := "https://m.weibo.cn/detail/" + weibo_id
	weibo_avatar := w.Mblog.User.AvatarHd
	// 	content = ['[CQ:image,file='+weibo_avatar+']']
	// 	content.append('\n')
	created_time := get_created_time(w.Mblog.CreatedAt)
	// 	# 保存要插入在content_list中的位置
	// 	content_index = len(content_list)
	// 	# 查询楼中楼
	// 	comment_created_time = await GetWeiboComment(weibo_id, mid, headers, uid, content_list, wbindex, weibo_url)
	// 	if now_comment_time < comment_created_time:
	// 		now_comment_time = comment_created_time
	// 	if not (last_weibo_time < created_time): # 不是新微博
	// 		continue
	if ((retweeted_status) && (retweet_id:= retweeted_status.ID; retweet_id)) {
		is_long_retweet = retweeted_status.get('isLongText')
		if is_long:
			weibo = get_long_weibo(weibo_id, headers, False) # 捕捉对象的长微博不截断
			if not weibo:
				weibo = parse_weibo(w['mblog'])
		else:
			weibo = parse_weibo(w['mblog'])
		if is_long_retweet:
			retweet = get_long_weibo(retweet_id, headers, True) # 转发的长微博截断
			if not retweet:
				retweet = parse_weibo(retweeted_status)
		else:
			retweet = parse_weibo(retweeted_status)
	}
		
}

func query() {
	index := getIndex()
	for _, card := range index.Data.Cards {
		if c := card.CardGroup; c != nil {
			switch c[1].CardType {
			case "9": // 发微博或评论
				if strings.Contains(*c[0].Desc, "评论") {
					//
				} else {
					fmt.Println(c[1].Mblog.Text)
				}
				handle(c[1])
			case "30": // 在线状态
				fmt.Println(*c[1].Desc1)
			}
		}
	}
}

func main() {
	query()
}
