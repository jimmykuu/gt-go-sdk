package geetest

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	VERSION      = "0.1"
	BASE_URL     = "http://api.geetest.com/get.php"
	API_SERVER   = "http://api.geetest.com/validate.php"
	API_REGISTER = "http://api.geetest.com/register.php"
)

type GeeTest struct {
	CaptchId   string
	PrivateKey string
}

func NewGeeTest(captchId, privateKey string) GeeTest {
	return GeeTest{
		PrivateKey: privateKey,
		CaptchId:   captchId,
	}
}

func (geeTest GeeTest) Challenge() string {
	return geeTest.get(fmt.Sprintf("%s?gt=%s", API_REGISTER, geeTest.CaptchId))
}

func (geeTest GeeTest) Validate(challenge, validate, seccode string) bool {
	if validate != geeTest.md5Value(geeTest.PrivateKey+"geetest"+challenge) {
		return false
	}

	values := url.Values{
		"seccode": {seccode},
		"version": {"go_" + VERSION},
	}

	backInfo := geeTest.postValue(API_SERVER, values)
	return backInfo == geeTest.md5Value(seccode)
}

func (geeTest GeeTest) md5Value(values string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(values)))
}

func (geeTest GeeTest) get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return ""
	}

	return string(body)
}

func (geeTest GeeTest) postValue(host string, values url.Values) string {
	resp, err := http.PostForm(host, values)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return ""
	}

	return string(body)
}

func (geeTest GeeTest) EmbedURL() string {
	return fmt.Sprintf("https://api.geetest.com/get.php?gt=%s&challenge=%s&product=embed", geeTest.CaptchId, geeTest.Challenge())
}

func (geeTest GeeTest) PopupURL(popupBtnId string) string {
	return fmt.Sprintf("https://api.geetest.com/get.php?gt=%s&challenge=%s&product=popup&popupbtnid=%s", geeTest.CaptchId, geeTest.Challenge(), popupBtnId)
}
