package handlers

import (
	"context"
	"gorm.io/gorm"
	"time"

	"github.com/bitdance-panic/gobuy/app/consts"
	rpc_agent "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/clients"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/callopt"
)

var db *gorm.DB

// HandleAddToCart 这是更新商品
// @Summary 这是一段Summary
// @Description 这是一段Description
// @Accept application/json
// @Produce application/json
// @Router /product [put]
func HandleAskAgent(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var req rpc_agent.AskReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	req.UserId = int32(userID)
	resp, err := clients.AgentClient.Ask(context.Background(), &req, callopt.WithRPCTimeout(30*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, utils.H{"content": resp.Content})
}

// HandleSendMessage 处理发送消息的请求
func HandleSendMessage(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var req rpc_agent.ChatMessageReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 设置sender_id为当前用户ID
	req.SenderId = int64(userID)
	req.ReceiverId = 1
	req.Content = c.Query("content")

	_, err := clients.AgentClient.SendMessage(context.Background(), &req, callopt.WithRPCTimeout(30*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, utils.H{})
}

func HistoryMessage(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt(consts.CONTEXT_UID_KEY)
	var req rpc_agent.HistoryMessageReq
	if err := c.BindAndValidate(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 设置sender_id为当前用户ID
	req.SenderId = int64(userID)

	resp, err := clients.AgentClient.GetHistoryMessages(context.Background(), &req, callopt.WithRPCTimeout(30*time.Second))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	//messages := make([]map[string]interface{}, len(resp))
	//for i, msg := range resp {
	//	messages[i] = map[string]interface{}{
	//		"id":         msg.Id,
	//		"senderId":   msg.SenderId,
	//		"receiverId": msg.ReceiverId,
	//		"content":    msg.Content,
	//		"timestamp":  msg.Timestamp.String(),
	//	}
	//}

	utils.Success(c, utils.H{"content": resp.Messages})
}
