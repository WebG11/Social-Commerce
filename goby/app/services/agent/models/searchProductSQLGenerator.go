package models

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func NewSearchProductSQLGenerator(ctx context.Context, columnsString string) (*openai.ChatModel, *prompt.DefaultChatTemplate, error) {
	temp := float32(0.7)
	llmConf := conf.GetConf().LLM
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:     llmConf.BaseURL,
		Model:       llmConf.ModelName,
		APIKey:      llmConf.ApiKey,
		Temperature: &temp,
	})
	if err != nil {
		return nil, nil, err
	}
	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role: schema.System,
			Content: `You are an SQL expert. Generate a MySQL query for the table named 'product' (wrap the table name in backticks). table 'product' with the following columns: ` +
				columnsString + `Return ONLY the raw SQL with no backticks, no Markdown, no code blocks, and no explanations.`,
		},
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。只获取is_deleted为false的商品",
		},
	)
	return chatModel, template, nil
}
