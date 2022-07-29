package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var CommentUrl = "https://m.weibo.cn/comments/hotflow?id=%s&mid=%s&max_id_type=0"
var CommentReq *http.Request

type Comment struct {
	Bid          string `json:"bid"`
	CommentBadge []struct {
		Actionlog struct {
			ActCode string `json:"act_code"`
			Ext     string `json:"ext"`
		} `json:"actionlog"`
		Length  float64 `json:"length"`
		Name    string  `json:"name"`
		NewIcon bool    `json:"new_icon"`
		PicURL  string  `json:"pic_url"`
		Scheme  string  `json:"scheme"`
	} `json:"comment_badge,omitempty"`
	Comments             interface{} `json:"comments,omitempty"`
	CreatedAt            string      `json:"created_at"`
	DisableReply         int         `json:"disable_reply"`
	FloorNumber          int         `json:"floor_number"`
	ID                   string      `json:"id"`
	IsLikedByMblogAuthor bool        `json:"isLikedByMblogAuthor"`
	LikeCount            int         `json:"like_count"`
	Liked                bool        `json:"liked"`
	MaxID                int         `json:"max_id"`
	Mid                  string      `json:"mid"`
	Pic                  *struct {
		Geo   interface{} `json:"geo"`
		Large struct {
			Geo  interface{} `json:"geo"`
			Size string      `json:"size"`
			URL  string      `json:"url"`
		} `json:"large"`
		Pid  string `json:"pid"`
		Size string `json:"size"`
		URL  string `json:"url"`
	} `json:"pic,omitempty"`
	Readtimetype    string `json:"readtimetype"`
	RestrictOperate int    `json:"restrictOperate"`
	Rootid          string `json:"rootid"`
	Rootidstr       string `json:"rootidstr"`
	SafeTags        *int   `json:"safe_tags,omitempty"`
	Source          string `json:"source"`
	Text            string `json:"text"`
	TotalNumber     int    `json:"total_number"`
	User            struct {
		AvatarHd          string         `json:"avatar_hd"`
		Badge             map[string]int `json:"badge"`
		CloseBlueV        bool           `json:"close_blue_v"`
		CoverImagePhone   string         `json:"cover_image_phone"`
		Description       string         `json:"description"`
		FollowCount       int            `json:"follow_count"`
		FollowMe          bool           `json:"follow_me"`
		FollowersCount    string         `json:"followers_count"`
		FollowersCountStr string         `json:"followers_count_str"`
		Following         bool           `json:"following"`
		Gender            string         `json:"gender"`
		ID                int            `json:"id"`
		Like              bool           `json:"like"`
		LikeMe            bool           `json:"like_me"`
		Mbrank            int            `json:"mbrank"`
		Mbtype            int            `json:"mbtype"`
		ProfileImageURL   string         `json:"profile_image_url"`
		ProfileURL        string         `json:"profile_url"`
		ScreenName        string         `json:"screen_name"`
		StatusesCount     int            `json:"statuses_count"`
		Urank             int            `json:"urank"`
		Verified          bool           `json:"verified"`
		VerifiedType      int            `json:"verified_type"`
	} `json:"user"`
}

type Comments struct {
	Data struct {
		Data      []Comment `json:"data"`
		Max       int       `json:"max"`
		MaxID     int       `json:"max_id"`
		MaxIDType int       `json:"max_id_type"`
		Status    struct {
			CommentManageInfo struct {
				ApprovalCommentType   int `json:"approval_comment_type"`
				CommentPermissionType int `json:"comment_permission_type"`
				CommentSortType       int `json:"comment_sort_type"`
			} `json:"comment_manage_info"`
		} `json:"status"`
		TotalNumber int `json:"total_number"`
	} `json:"data"`
	Ok int `json:"ok"`
}

type Geo1 struct {
	Croped bool `json:"croped"`
	Height int  `json:"height"`
	Width  int  `json:"width"`
}

type Geo2 struct {
	Croped bool   `json:"croped"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

func GetComment(Bid, Mid string) string {
	url := fmt.Sprintf(CommentUrl, Bid, Mid)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var comments Comments
	err = json.Unmarshal(b, &comments)
	if checkErr(err) {
		for _, comment := range comments.Data.Data {
			if comment.User.ID == cfg.Account.UID {
				fmt.Println(comment.User.ScreenName, comment.Text)
			}
			switch coms := comment.Comments.(type) {
			case bool:
				continue
			case []Comment:
				for _, c := range coms {
					if c.User.ID == cfg.Account.UID {
						fmt.Println(c.User.ScreenName, c.Text)
					}
				}
			}
		}
		return string(b)
	} else {
		log.Error(Bid, Mid)
		return ""
	}
}
