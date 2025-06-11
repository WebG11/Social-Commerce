package models

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func NewUtilSelectorModel(ctx context.Context, tools *[]tool.BaseTool) (*openai.ChatModel, *prompt.DefaultChatTemplate, error) {
	llmConf := conf.GetConf().LLM
	temp := float32(0.7)
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL:     llmConf.BaseURL,
		Model:       llmConf.ModelName,
		APIKey:      llmConf.ApiKey,
		Temperature: &temp,
	})
	if err != nil {
		return nil, nil, err
	}
	toolInfos := make([]*schema.ToolInfo, 0, len(*tools))
	for _, tool := range *tools {
		info, err := tool.Info(ctx)
		if err != nil {
			return nil, nil, err
		}
		toolInfos = append(toolInfos, info)
	}
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		return nil, nil, err
	}
	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role: schema.System,
			Content: `你是一个工具挑选者，你的任务是挑选工具，并将用户的输入完全按照原样传入工具中，
			不要添加任何额外的内容或评论。不要对输入进行解释、总结或修改。`,
		},
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。当前请求用户的id为{userID}",
		},
	)
	return chatModel, template, nil
}
