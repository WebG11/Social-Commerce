package tools

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dao"
	chat_models "github.com/bitdance-panic/gobuy/app/services/agent/models"
	agentutils "github.com/bitdance-panic/gobuy/app/services/agent/utils"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

var (
	orderTemplate     *prompt.DefaultChatTemplate
	orderSqlGenerator *openai.ChatModel
)

func NewSearchOrdersTool() tool.BaseTool {
	InitSearchOrders()
	searchProductsTool, err := utils.InferTool("outer_search_orders", "Search for orders based on user requirements.", searchOrdersFunc)
	if err != nil {
		log.Fatal(err)
	}
	return searchProductsTool
}

type SearchOrdersParams struct {
	Prompt string `json:"sql" jsonschema:"description=User prompt"`
}

func searchOrdersFunc(ctx context.Context, params *SearchOrdersParams) (*ToolResponse, error) {
	messages, err := orderTemplate.Format(context.Background(), map[string]any{
		"task": params.Prompt,
	})
	if err != nil {
		return nil, err
	}
	sqlResult, err := orderSqlGenerator.Generate(ctx, messages)
	// 一般是没找到工具就进这里
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
		return nil, errors.New("Tools not found")
	}
	log.Println("大模型生成的SQL:", sqlResult.Content)
	sql := agentutils.CleanBlock(sqlResult.Content)
	log.Println("clean后的SQL:", sql)
	orders := make([]models.Order, 0)
	// 具体的调用逻辑
	result := tidb.DB.Raw(sql).Scan(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	var idsString strings.Builder
	// 多次追加字符串
	for i, order := range orders {
		if i != 0 {
			idsString.WriteString(", ")
		}
		idsString.WriteString(strconv.Itoa(order.ID))
		log.Printf("ID: %d", order.ID)
	}
	return &ToolResponse{
		Data:            idsString.String(),
		DataDescription: "获取的为各个订单ID,用逗号分隔",
		ShowWay:         "将各个订单ID改为<a>的超链接,href格式为 http://localhost:3000/orders/:id",
	}, nil
}

func InitSearchOrders() {
	columns, err := dao.GetOrderColumns(tidb.DB)
	if err != nil {
		log.Fatal(err.Error())
	}
	columnsString := strings.Join(columns, ", ")
	log.Printf("order表字段为: %+v", columnsString)
	orderSqlGenerator, orderTemplate, err = chat_models.NewSearchOrderSQLGenerator(context.Background(), columnsString)
	if err != nil {
		log.Panic(err.Error())
	}
}
