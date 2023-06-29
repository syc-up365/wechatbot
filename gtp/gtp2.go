package gtp

import (
	"bytes"
	"encoding/json"
	"github.com/869413421/wechatbot/config"
	"io/ioutil"
	"log"
	"net/http"
)

const BASEURL2 = "https://api.openai.com/v1/chat/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody2 struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Choices []ChoicesBody          `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoicesBody struct {
	Index        int               `json:"index"`
	Message      map[string]string `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

type ChoiceItem2 struct {
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody2 struct {
	Model    string                    `json:"model"`
	Messages [1]map[string]interface{} `json:"messages"`
	//TopP        int     `json:"top_p"`
	//N           int     `json:"n"`
	//Stream      bool    `json:"stream"`
	////LogProbs    string  `json:"logprobs"`
	//Stop string `json:"stop"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/chat/completions \
//-H "Content-Type: application/json" \
//-H "Authorization: Bearer $OPENAI_API_KEY" \
//-d '{
//"model": "gpt-3.5-turbo",
//"messages": [{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": "Hello!"}]
//}'

func Completions2(msg string) (string, error) {
	var msgMap [1]map[string]interface{}

	//var system map[string]interface{} /*创建集合 */
	//system = make(map[string]interface{}, 2)
	//system["role"] = "system"
	//system["content"] = "You are a helpful assistant."
	//msgMap[0] = system

	var user map[string]interface{} /*创建集合 */
	user = make(map[string]interface{}, 2)
	user["role"] = "user"
	user["content"] = msg
	msgMap[0] = user

	requestBody := ChatGPTRequestBody2{
		Model:    "gpt-3.5-turbo",
		Messages: msgMap,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL2+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody2{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string

	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v.Message["content"]
			break
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
