package main

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/bll"
	"gorm.io/gorm"
	"time"

	rpc_agent "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent"
)

var db *gorm.DB

type AgentServiceImpl struct{}

func (*AgentServiceImpl) Ask(ctx context.Context, req *rpc_agent.AskReq) (resp *rpc_agent.AskResp, err error) {
	return bll.Ask(ctx, req)
}

func (*AgentServiceImpl) SendMessage(ctx context.Context, req *rpc_agent.ChatMessageReq) (resp *rpc_agent.ChatMessageResp, err error) {
	return bll.SendMessage(ctx, req)
}

func (*AgentServiceImpl) GetHistoryMessages(ctx context.Context, req *rpc_agent.HistoryMessageReq) (resp *rpc_agent.HistoryMessageResp, err error) {
	// 调用 bll 层的方法获取历史消息
	messages, err := bll.GetHistoryMessages(ctx, &rpc_agent.HistoryMessageReq{
		SenderId: req.SenderId,
	})
	if err != nil {
		return nil, err
	}

	// 将 []models.ChatMessage 转换为 *agent.HistoryMessageResp
	rpcResp := &rpc_agent.HistoryMessageResp{
		Messages: make([]*rpc_agent.ChatMessage, len(messages)),
	}
	for i, msg := range messages {
		rpcResp.Messages[i] = &rpc_agent.ChatMessage{
			Id:         msg.ID,
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			Timestamp:  msg.Timestamp.Format(time.RFC3339),
		}
	}

	return rpcResp, nil
}
