package main

type Mblog struct {
	Desc  *string `json:"desc,omitempty"`
	Desc1 *string `json:"desc1,omitempty"`
	Mblog *Mblog  `json:"retweeted_status"`

	Bid        string   `json:"bid"`
	CreatedAt  string   `json:"created_at"`
	IsLongText bool     `json:"isLongText"`
	Mid        string   `json:"mid"`
	PicIds     []string `json:"pic_ids"`
	Source     string   `json:"source"`
	Text       string   `json:"text"`
	User       struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

type Result struct {
	Data struct {
		Cards []struct {
			AsyncAPI  *string `json:"async_api,omitempty"`
			CardGroup []struct {
				CardType string `json:"card_type"`
				Mblog    Mblog  `json:"mblog,omitempty"`
			} `json:"card_group,omitempty"`
			CardType     interface{} `json:"card_type"`
			CardTypeName *string     `json:"card_type_name,omitempty"`
			IsAsyn       *int        `json:"is_asyn,omitempty"`
			Itemid       string      `json:"itemid"`
			ShowType     *int        `json:"show_type,omitempty"`
		} `json:"cards"`
		PageInfo struct {
			Adhesive          int    `json:"adhesive"`
			Adid              string `json:"adid"`
			BackgroundScheme  string `json:"background_scheme"`
			CardlistHeadCards []struct {
				ChannelList []struct {
					Containerid       string  `json:"containerid"`
					DefaultAdd        int     `json:"default_add"`
					HighlightIcon     *string `json:"highlight_icon,omitempty"`
					HighlightIconDark *string `json:"highlight_icon_dark,omitempty"`
					ID                string  `json:"id"`
					Icon              *string `json:"icon,omitempty"`
					IconDark          *string `json:"icon_dark,omitempty"`
					MustShow          int     `json:"must_show"`
					Name              string  `json:"name"`
					Scheme            string  `json:"scheme"`
				} `json:"channel_list"`
				HeadType     int    `json:"head_type"`
				HeadTypeName string `json:"head_type_name"`
				MenuScheme   string `json:"menu_scheme"`
				ShowMenu     bool   `json:"show_menu"`
			} `json:"cardlist_head_cards"`
			Containerid     string   `json:"containerid"`
			ContaineridBak  string   `json:"containerid_bak"`
			Desc            string   `json:"desc"`
			DescMore        []string `json:"desc_more"`
			DetailDesc      string   `json:"detail_desc"`
			DisplayArrow    int      `json:"display_arrow"`
			FollowRelation  int      `json:"follow_relation"`
			HuatiExposure   string   `json:"huati_exposure"`
			IsObturate      int      `json:"is_obturate"`
			IsRefuse        int      `json:"is_refuse"`
			Nick            string   `json:"nick"`
			Oid             string   `json:"oid"`
			PageAttr        string   `json:"page_attr"`
			PageCommonExt   string   `json:"page_common_ext"`
			PageSize        int      `json:"page_size"`
			PageTitle       string   `json:"page_title"`
			PageTypeName    string   `json:"page_type_name"`
			PageURL         string   `json:"page_url"`
			PageView        string   `json:"page_view"`
			Portrait        string   `json:"portrait"`
			PortraitHd      string   `json:"portrait_hd"`
			PortraitSubText string   `json:"portrait_sub_text"`
			RefuseMsg       string   `json:"refuse_msg"`
			SearchScheme    string   `json:"search_scheme"`
			SharePageid     string   `json:"share_pageid"`
			ShowStyle       int      `json:"show_style"`
			SinceID         string   `json:"since_id"`
			TitleIcon       string   `json:"title_icon"`
			TitleTop        string   `json:"title_top"`
			VP              string   `json:"v_p"`
		} `json:"pageInfo"`
		Scheme      string `json:"scheme"`
		ShowAppTips int    `json:"showAppTips"`
	} `json:"data"`
	Ok int `json:"ok"`
}
type Geo struct {
	Croped bool        `json:"croped"`
	Height interface{} `json:"height"`
	Width  interface{} `json:"width"`
}
