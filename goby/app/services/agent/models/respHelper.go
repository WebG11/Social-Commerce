package models

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func NewRespHelperModel(ctx context.Context) (*openai.ChatModel, *prompt.DefaultChatTemplate, error) {
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
	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role:    schema.System,
			Content: "你是一个回复整理者。你能根据用户的请求以及系统查询到的项的字符串及其描述,返回给前端HTML的项展示以及对这些项的中文总结描述。项标签为<ol>和<li>，li的a之前要有个<span>放序号和分隔用的点。最外层不要添加```html。描述与列表项区分开，通过<br>标签与最后一个列表项换行分隔，描述不是原文复制，只用介绍数据是什么就行。超链接要有红色样式",
		},
		&schema.Message{
			Role:    schema.User,
			Content: "用户请求为: {task}。系统查询结果为: {data},结果的描述为: {description}。展示的方式为: {showway}",
		},
	)
	return chatModel, template, nil
}
