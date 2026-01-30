package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ZhipuAIResponse 是我们期望从智谱AI获取的营养信息结构体。
// 它严格对应 promptText 中要求返回的JSON格式。
type ZhipuAIResponse struct {
	Name            string  `json:"name"`             // 菜品名称
	Calories100g    float64 `json:"calories_100g"`    // 单位：千卡 (kcal/100g)
	Protein100g     float64 `json:"protein_100g"`     // 单位：克 (g/100g)
	Fat100g         float64 `json:"fat_100g"`         // 单位：克 (g/100g)
	Carbs100g       float64 `json:"carbs_100g"`       // 单位：克 (g/100g)
	EstimatedWeight int     `json:"estimated_weight"` // 估算图这一整份食物的大致重量(克)
	Reasoning       string  `json:"reasoning"`        // 简短的分析理由
}

// ZhipuAIClient 定义了与智谱AI服务交互的客户端接口。
type ZhipuAIClient interface {
	// RecognizeFood 从给定的图片输入流中识别菜品，并返回其详细营养信息。
	// ctx: 请求的上下文，用于控制请求的生命周期 (如超时)。
	// imageReader: 包含图片数据的io.Reader接口。
	RecognizeFood(ctx context.Context, imageReader io.Reader) (*ZhipuAIResponse, error)

	// AnalyzeFoodText 根据用户提供的文本描述分析食物，并返回其详细营养信息。
	// ctx: 请求的上下文。
	// text: 用户输入的食物描述。
	AnalyzeFoodText(ctx context.Context, text string) (*ZhipuAIResponse, error)
}

const (
	zhipuAPIURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	promptText  = `你是一位经验丰富的临床营养师。请分析图片中的食物。
		核心任务：
		1. 准确识别菜品名称（如果是中餐，请识别具体菜名）。
		2. 分析其食材构成、烹饪方式（如是否重油、红烧、清蒸）。
		3. **重点：** 估算该食物 **每100克 (Per 100g)** 的营养成分数据。

		请严格遵守以下 JSON 格式返回，不要包含 Markdown 格式标记（如 '''json），不要包含任何额外的解释文本：

		{
			"name": "菜品名称",
			"calories_100g": 0,    // 单位：千卡 (kcal/100g)，必须是数字
			"protein_100g": 0.0,   // 单位：克 (g/100g)，保留1位小数
			"fat_100g": 0.0,       // 单位：克 (g/100g)，保留1位小数
			"carbs_100g": 0.0,     // 单位：克 (g/100g)，保留1位小数
			"estimated_weight": 0, // 估算图这一整份食物的大致重量(克)
			"reasoning": "简短的分析理由，例如：'这是西红柿鸡蛋面，面条看起来较宽，汤汁浓郁，且有明显油光，因此每100g热量略高于清汤面。'"
		}`

	promptTextAnalysis = `你是一位经验丰富的临床营养师。请根据用户提供的食物描述进行分析。
		核心任务：
		1. 准确识别/推断菜品名称。
		2. 分析其食材构成、烹饪方式。
		3. **重点：** 估算该食物 **每100克 (Per 100g)** 的营养成分数据。

		请严格遵守以下 JSON 格式返回，不要包含 Markdown 格式标记（如 '''json），不要包含任何额外的解释文本：

		{
			"name": "菜品名称",
			"calories_100g": 0,    // 单位：千卡 (kcal/100g)，必须是数字
			"protein_100g": 0.0,   // 单位：克 (g/100g)，保留1位小数
			"fat_100g": 0.0,       // 单位：克 (g/100g)，保留1位小数
			"carbs_100g": 0.0,     // 单位：克 (g/100g)，保留1位小数
			"estimated_weight": 0, // 根据描述估算这一整份食物的大致重量(克)，如果描述中未提及分量，请给出一个常见的标准份量
			"reasoning": "简短的分析理由"
		}
		
		注意：直接返回 JSON 对象，不要输出 "好的"、"分析如下" 等任何多余字符。`
)

// zhipuAIClientImpl 是 ZhipuAIClient 接口的具体实现。
type zhipuAIClientImpl struct {
	apiKey     string       // 用于认证的智谱AI API密钥。
	httpClient *http.Client // 执行HTTP请求的客户端。
}

// NewZhipuAIClient 创建并返回一个新的智谱AI客户端实例。
// apiKey: 从配置文件中获取的智谱AI API Key。
func NewZhipuAIClient(apiKey string) ZhipuAIClient {
	return &zhipuAIClientImpl{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 设置HTTP请求的超时时间为60秒。
		},
	}
}

// RecognizeFood 实现了从图片识别食物并获取营养信息的核心逻辑。
// 它将图片数据编码后，连同prompt一起发送到智谱AI模型，并解析返回的结果。
func (c *zhipuAIClientImpl) RecognizeFood(ctx context.Context, imageReader io.Reader) (*ZhipuAIResponse, error) {
	imgBytes, err := io.ReadAll(imageReader)
	if err != nil {
		return nil, fmt.Errorf("读取图片数据失败: %w", err)
	}

	// 将图片二进制数据编码为 Base64 字符串。
	imageBase64 := base64.StdEncoding.EncodeToString(imgBytes)

	// 构建发送给智谱API的请求体。
	requestPayload := zhipuRequest{
		Model: "glm-4v", // Update to correct vision model name if needed, usually glm-4v or glm-4v-plus
		Messages: []zhipuMessage{
			{
				Role: "user",
				Content: []zhipuContent{
					{
						Type: "text",
						Text: promptText,
					},
					{
						Type: "image_url",
						ImageURL: &zhipuImageURL{
							// 图片部分是Base64编码的图片数据。
							URL: fmt.Sprintf("data:image/jpeg;base64,%s", imageBase64),
						},
					},
				},
			},
		},
		Thinking: &zhipuThinking{
			Type: "disabled",
		},
	}

	body, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, zhipuAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头，包括内容类型和认证信息。
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求到智谱API失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("智谱API返回非200状态码: %d, 响应体: %s", resp.StatusCode, string(respBody))
	}

	// 将响应体JSON反序列化
	var zhipuResp zhipuAPICompletionResponse
	if err := json.Unmarshal(respBody, &zhipuResp); err != nil {
		return nil, fmt.Errorf("反序列化智谱API响应失败: %w", err)
	}

	// 检查返回的内容是否为空。
	if len(zhipuResp.Choices) == 0 || zhipuResp.Choices[0].Message.Content == "" {
		return nil, errors.New("从智谱API接收到空内容")
	}

	content := zhipuResp.Choices[0].Message.Content
	// AI返回的内容可能被Markdown的代码块标记包裹，需要清理
	content = cleanJSONString(content)

	// 将清理后的内容反序列化到最终的ZhipuAIResponse结构体中
	var nutritionInfo ZhipuAIResponse
	if err := json.Unmarshal([]byte(content), &nutritionInfo); err != nil {
		fmt.Printf("Zhipu Unmarshal Error. Raw Content: %s\n", content)
		return nil, fmt.Errorf("AI返回数据格式错误，无法解析为营养信息")
	}

	return &nutritionInfo, nil
}

// AnalyzeFoodText 实现了根据文本描述分析食物并获取营养信息的核心逻辑。
func (c *zhipuAIClientImpl) AnalyzeFoodText(ctx context.Context, text string) (*ZhipuAIResponse, error) {
	// 构建发送给智谱API的请求体。
	// 使用 glm-4-flash 时，建议使用 system + user 的简单字符串格式，以获得最佳遵循效果。
	requestPayload := zhipuRequest{
		Model: "glm-4-flash",
		Messages: []zhipuMessage{
			{
				Role:    "system",
				Content: promptTextAnalysis,
			},
			{
				Role:    "user",
				Content: text,
			},
		},
		Thinking: &zhipuThinking{
			Type: "disabled",
		},
	}

	body, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, zhipuAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头，包括内容类型和认证信息。
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求到智谱API失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("智谱API返回非200状态码: %d, 响应体: %s", resp.StatusCode, string(respBody))
	}

	// 将响应体JSON反序列化
	var zhipuResp zhipuAPICompletionResponse
	if err := json.Unmarshal(respBody, &zhipuResp); err != nil {
		return nil, fmt.Errorf("反序列化智谱API响应失败: %w", err)
	}

	// 检查返回的内容是否为空。
	if len(zhipuResp.Choices) == 0 || zhipuResp.Choices[0].Message.Content == "" {
		return nil, errors.New("从智谱API接收到空内容")
	}

	content := zhipuResp.Choices[0].Message.Content
	// AI返回的内容可能被Markdown的代码块标记包裹，需要清理
	content = cleanJSONString(content)

	// 将清理后的内容反序列化到最终的ZhipuAIResponse结构体中
	var nutritionInfo ZhipuAIResponse
	if err := json.Unmarshal([]byte(content), &nutritionInfo); err != nil {
		fmt.Printf("Zhipu Unmarshal Error (Text). Raw Content: %s\n", content)
		return nil, fmt.Errorf("AI返回数据格式错误，无法解析为营养信息")
	}

	return &nutritionInfo, nil
}

// cleanJSONString 移除JSON字符串前后可能存在的Markdown代码块标记或无关文本。
// 它通过查找第一个 '{' 和最后一个 '}' 来提取有效的JSON部分。
func cleanJSONString(s string) string {
	s = strings.TrimSpace(s)

	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start != -1 && end != -1 && start < end {
		return s[start : end+1]
	}

	return s
}

// --- Zhipu API 内部数据结构定义 ---

// zhipuThinking 配置模型的思考过程选项。
type zhipuThinking struct {
	Type string `json:"type"` // 类型：enabled/disabled,默认为 disabled,耗时短。
}

// zhipuRequest 是发送给智谱API的请求体结构。
type zhipuRequest struct {
	Model    string         `json:"model"`              // 使用的模型名称, e.g., "glm-4.6v"。
	Messages []zhipuMessage `json:"messages"`           // 对话消息列表。
	Thinking *zhipuThinking `json:"thinking,omitempty"` // 思考过程配置。
}

// zhipuMessage 代表一次对话中的一条消息。
type zhipuMessage struct {
	Role    string      `json:"role"`    // 消息发送者的角色, e.g., "user", "system"。
	Content interface{} `json:"content"` // 消息的具体内容，可以是字符串或多部分数组。
}

// zhipuContent 代表消息中的一个内容块 (文本或图片)。
type zhipuContent struct {
	Type     string         `json:"type"`                // 内容类型, "text" 或 "image_url"。
	Text     string         `json:"text,omitempty"`      // 文本内容，当Type为"text"时使用。
	ImageURL *zhipuImageURL `json:"image_url,omitempty"` // 图片URL，当Type为"image_url"时使用。
}

// zhipuImageURL 包含图片的URL地址。
type zhipuImageURL struct {
	URL string `json:"url"` // 图片的URL，这里使用Data URI格式。
}

// zhipuAPICompletionResponse 是智谱API返回的完整响应结构。
type zhipuAPICompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"` // AI模型生成的核心内容。
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`     // 输入的token数量。
		CompletionTokens int `json:"completion_tokens"` // 输出的token数量。
		TotalTokens      int `json:"total_tokens"`      // 总的token数量。
	} `json:"usage"`
}
