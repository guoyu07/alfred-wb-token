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
	var title string
	var subtitle string
	if len(token) == 0 {
		title = "请输入 token"
	} else {
		result := showTokenInfo(token)
		subtitle = result
		var tokenJSONObject tokenObject
		err := json.Unmarshal([]byte(result), &tokenJSONObject)
		if err == nil {
			if tokenJSONObject.ErrorCode > 0 {
				title = fmt.Sprintf("错误代码 %v => %v\n", tokenJSONObject.ErrorCode, tokenJSONObject.Error)
			} else {
				title = fmt.Sprintf("%v 秒(%v分钟)\n", tokenJSONObject.ExpireIn, tokenJSONObject.ExpireIn/60)
			}
		} else {
			title = "查询出错啦"
		}
	}
	resp := alfred.NewResponse()
	item := alfred.ResponseItem{
		Valid:    true,
		UID:      strconv.FormatInt(int64(time.Now().Nanosecond()), 10),
		Title:    title,
		Subtitle: subtitle,
		Arg:      "",
		Icon:     "",
	}
	resp.AddItem(&item)

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
