package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"encoding/json"
	"time"

	"strconv"

	alfred "github.com/ruedap/go-alfred"
)

const (
	// http://open.weibo.com/wiki/Oauth2/get_token_info
	URL_TOKEN_INFO_POST string = "https://api.weibo.com/oauth2/get_token_info"
)

type tokenObject struct {
	UID      int         `json:"uid"`
	Appkey   string      `json:"appkey"`
	Scope    interface{} `json:"scope"`
	CreateAt int         `json:"create_at"`
	ExpireIn int         `json:"expire_in"`

	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Request   string `json:"request"`
}

func showTokenInfo(accessToken string) string {
	form := url.Values{}
	form.Add("access_token", accessToken)

	req, err := http.NewRequest("POST", URL_TOKEN_INFO_POST, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	return string(responseBody)
}

func main() {
	token := ""
	if len(os.Args) > 1 {
		for _, q := range os.Args[1:] {
			token += q
		}
	}
	resp := alfred.NewResponse()

	if len(token) == 0 {
		title := "请输入 token"
		addItem(resp, title)

	} else {
		result := showTokenInfo(token)
		var tokenJSONObject tokenObject
		err := json.Unmarshal([]byte(result), &tokenJSONObject)
		if err == nil {
			if tokenJSONObject.ErrorCode > 0 {
				title := fmt.Sprintf("错误代码: %v", tokenJSONObject.ErrorCode)
				addItem(resp, title)
				addItem(resp, "错误信息:"+tokenJSONObject.Error)
			} else {
				title := fmt.Sprintf("%v 秒(%v分钟)\n", tokenJSONObject.ExpireIn, tokenJSONObject.ExpireIn/60)
				addItem(resp, title)
				t := time.Unix(0, int64(tokenJSONObject.ExpireIn+tokenJSONObject.CreateAt)*1000*int64(time.Millisecond))
				addItem(resp, "到期时间: "+fmt.Sprintf("%v\n", t.Format("2006-01-02 15:04:05")))
				addItem(resp, "UID: "+fmt.Sprintf("%v\n", tokenJSONObject.UID))
				addItem(resp, "AppKey: "+fmt.Sprintf("%v\n", tokenJSONObject.Appkey))
			}
		} else {
			title := "查询出错啦"
			addItem(resp, title)
		}
	}

	xml, err := resp.ToXML()
	if err != nil {
		title := fmt.Sprintf("Error: %v", err.Error())
		subtitle := "出错啦"
		arg := title
		errXML := alfred.ErrorXML(title, subtitle, arg)
		fmt.Println(errXML)
		return
	}
	fmt.Println(xml)
}

func addItem(response *alfred.Response, title string) {
	item := alfred.ResponseItem{
		Valid:    true,
		UID:      strconv.FormatInt(int64(time.Now().Nanosecond()), 10),
		Title:    title,
		Subtitle: "",
		Arg:      "",
		Icon:     "",
	}
	response.AddItem(&item)
}
