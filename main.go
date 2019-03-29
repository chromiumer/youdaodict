package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FROM   = "en"                               //源语言
	TO     = "zh_CHS"                           //目标语言
	APPKEY = "2d70c401bde9b8c0"                 //AppKey=应用ID
	SECRET = "DHbmxeAXfajyWlR1hQJFtgWNtFpy1qCl" //应用密钥
)

type Jsons struct {
	Query string `json:"query"`
	Web   []Web  `json:"web"`
	Basic Basic  `json:"basic"`
}

type Web struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

type Basic struct {
	Examtype   []string `json:"exam_type"`
	Usphonetic string   `json:"us-phonetic"`
	Ukphonetic string   `json:"uk-phonetic"`
	Explains   []string `json:"explains"`
}

//生成加盐随机数
func GenRandNumber(n int) int {

	rand.Seed(time.Now().UnixNano())
	m := rand.Intn(n)

	return m
}

//生成签名
func genSign(appkey string, words string, salt string, secret string) string {

	h := md5.New()
	h.Write([]byte(appkey + words + salt + secret))
	cipherStr := h.Sum(nil)
	signresult := strings.ToUpper(hex.EncodeToString(cipherStr))
	return signresult
}

func main() {

	var sum string

	for i := 1; i < len(os.Args); i++ {
		sum += os.Args[i] + " "
	}

	//fmt.Printf("%v\n", sum)

	WORDS := sum

	SALT := strconv.Itoa(GenRandNumber(100))
	SIGN := genSign(APPKEY, WORDS, SALT, SECRET)

	req, _ := http.NewRequest("GET", "https://openapi.youdao.com/api", nil)
	q := req.URL.Query()
	q.Add("q", WORDS)
	q.Add("from", FROM)
	q.Add("to", TO)
	q.Add("appKey", APPKEY)
	q.Add("salt", SALT)
	q.Add("sign", SIGN)

	req.URL.RawQuery = q.Encode()
	//fmt.Println(req.URL.String())

	resp, _ := http.Get(req.URL.String())
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	m := Jsons{}
	json.Unmarshal([]byte(body), &m)

	//音标
	fmt.Printf("\n%c[1;;32m英: [%v]  美: [%v]%c[0m \n\n", 0x1B, m.Basic.Ukphonetic, m.Basic.Usphonetic, 0x1B)

	//基本释义

	for a := 0; a < len(m.Basic.Explains); a++ {
		fmt.Printf("%c[1;;32m%v%c[0m\n", 0x1B, m.Basic.Explains[a], 0x1B)
	}

	//短语
	fmt.Printf("\n%c[1;;32m%v%c[0m\n", 0x1B, "短语: ", 0x1B)
	for i := 0; i < len(m.Web); i++ {

		fmt.Printf("%c[1;;32m%v%c[0m  ", 0x1B, m.Web[i].Key, 0x1B)

		for j := 0; j < len(m.Web[i].Value); j++ {
			fmt.Printf("%c[1;;32m%v;%c[0m", 0x1B, m.Web[i].Value[j], 0x1B)
		}
		fmt.Printf("\n")

	}

	//考试级别
	fmt.Printf("\n%c[1;;32m考试级别 %v%c[0m\n\n", 0x1B, m.Basic.Examtype, 0x1B)
}
