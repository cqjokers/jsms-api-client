package jsms_api_client

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/golang/glog"
	"encoding/base64"
	"encoding/json"
)

const URL = "https://api.sms.jpush.cn/v1/"

type jsms struct {
	appKey       string
	masterSecret string
}

func New(appKey, masterSecret string) (*jsms) {
	smsClient := &jsms{
		appKey:       appKey,
		masterSecret: masterSecret,
	}
	return smsClient
}

//发送文本验证码短信
func (this *jsms) SendCode(mobile, temp_id string) []byte {
	url := URL + "codes"
	body := make(map[string]interface{})
	body["mobile"] = mobile
	body["temp_id"] = temp_id
	resp := this.request("post", url, body)
	fmt.Println(resp)
	return resp
}

//发送语音验证码短信
func (this *jsms) SendVoiceCode(mobile string) []byte{
	url := URL + "voice_codes"
	body := make(map[string]interface{})
	body["mobile"] = mobile
	return this.request("post",url,body)
}

//验证码验证
func (this *jsms) CheckCode(msg_id, code string) []byte {
	url := URL + "codes/" + msg_id + "/valid"
	body := make(map[string]interface{})
	body["ode"] = code
	return this.request("POST", url, body)
}

//发送单条模板短信、单条定时短信
func (this *jsms) SendMessage(mobile, tempId string, temp_para map[string]interface{}, time string) []byte {
	path := "messages"
	body := make(map[string]interface{})
	body["mobile"] = mobile
	body["temp_id"] = tempId

	if (temp_para != nil && len(temp_para) > 0) {
		body["temp_para"] = temp_para
	}
	if (time != "") {
		path = "schedule"
		body["send_time"] = time
	}
	url := URL + path
	return this.request("POST", url, body)
}


func (this *jsms) request(method, url string, body map[string]interface{}) []byte {
	bodyByte, _ := json.Marshal(body)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(bodyByte)))
	defer glog.Flush()
	if err != nil {
		glog.Error("[jsms]-请求出错，" + err.Error())
		panic(err)
	}
	auth := this.appKey + ":" + this.masterSecret
	req.Header.Add("Authorization", "Basic "+string(base64.StdEncoding.EncodeToString([]byte(auth))))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody
}
