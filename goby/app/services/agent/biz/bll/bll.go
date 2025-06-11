package bll

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitdance-panic/gobuy/app/models"
	rpc_agent "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dao"
	chat_models "github.com/bitdance-panic/gobuy/app/services/agent/models"
	"github.com/bitdance-panic/gobuy/app/services/agent/tools"
	"github.com/bitdance-panic/gobuy/app/services/agent/utils"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"
	"log"
)

var agent compose.Runnable[map[string]any, string]

var respHelper *openai.ChatModel
var respHelperTemplate *prompt.DefaultChatTemplate

func initAgent() {
	ctx := context.Background()
	tools := []tool.BaseTool{
		tools.NewSearchProductsTool(),
		tools.NewSearchOrdersTool(),
	}
	utilSelector, template, err := chat_models.NewUtilSelectorModel(ctx, &tools)
	if err != nil {
		log.Fatal(err)
	}
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatal(err)
	}
	chain := compose.NewChain[map[string]any, string]()
	chain.
		AppendChatTemplate(template).
		AppendChatModel(utilSelector).
		AppendToolsNode(toolsNode).
		AppendLambda(compose.InvokableLambda(func(_ context.Context, results []*schema.Message) (string, error) {
			return results[0].Content, nil
		}))
	agent, err = chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func initHelper() {
	var err error
	respHelper, respHelperTemplate, err = chat_models.NewRespHelperModel(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func Init() {
	initAgent()
	initHelper()
}

func Ask(ctx context.Context, req *rpc_agent.AskReq) (resp *rpc_agent.AskResp, err error) {
	variables := map[string]any{
		"task":   req.UserPrompt,
		"userID": req.UserId,
	}
	// 调用大模型
	invokeResp, err := agent.Invoke(ctx, variables)
	// 一般是没找到工具就进这里
	if err != nil {
		return nil, err
	}
	var toolResp tools.ToolResponse
	err = json.Unmarshal([]byte(invokeResp), &toolResp)
	if err != nil {
		return nil, err
	}
	messages, err := respHelperTemplate.Format(ctx, map[string]any{
		"task":        variables["task"],
		"data":        toolResp.Data,
		"description": toolResp.DataDescription,
		"showway":     toolResp.ShowWay,
	})
	if err != nil {
		return nil, err
	}
	// TODO stream
	msgs, err := respHelper.Generate(ctx, messages)
	if err != nil {
		return nil, err
	}
	msgs.Content = utils.CleanBlock(msgs.Content)
	fmt.Println("大模型生成的结果:", msgs.Content)
	return &rpc_agent.AskResp{
		Content: msgs.Content,
	}, nil
}

// SendMessage 发送消息
func SendMessage(ctx context.Context, req *rpc_agent.ChatMessageReq) (*rpc_agent.ChatMessageResp, error) {
	// 从请求对象中提取 senderID, receiverID 和 content
	senderID := req.GetSenderId()
	receiverID := req.GetReceiverId()
	content := req.GetContent()

	// 调用 DAO 层的方法保存消息到数据库
	err := dao.SendMessage(tidb.DB, senderID, receiverID, content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetHistoryMessages 获取历史消息
func GetHistoryMessages(ctx context.Context, req *rpc_agent.HistoryMessageReq) ([]models.ChatMessage, error) {
	// 从请求对象中提取 senderID, receiverID 和 limit
	senderID := req.GetSenderId()
	receiverID := 1
	limit := 50

	// 调用 DAO 层的方法获取历史消息
	messages, err := dao.GetHistoryMessages(tidb.DB, senderID, int64(receiverID), limit)
	if err != nil {
		klog.Errorf("Failed to fetch history messages: %v", err)
		return nil, err
	}

	return messages, nil
}
