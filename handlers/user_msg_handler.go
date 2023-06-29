package handlers

import (
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	var replyName string
	if msg.IsSendByGroup() {

		senderInGroup, err := msg.SenderInGroup()
		log.Printf("Received Group %v sender %v Text Msg : %v err %v", senderInGroup.NickName, msg.Content, err)
		replyName = senderInGroup.NickName
		//senderInGroup.Self().NickName
	} else {
		sender, err := msg.Sender()
		log.Printf("Received Group %v sender %v Text Msg : %v err %v", sender.NickName, msg.Content, err)
		replyName = ""
	}

	// 向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")
	requestText = strings.TrimSpace(strings.ReplaceAll(requestText, "@bot孙", ""))
	log.Printf("requestText  %v", requestText)
	reply, err := gtp.Completions2(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")

	if !strings.EqualFold(replyName, "") {
		reply = "@" + replyName + " " + reply
	}

	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
