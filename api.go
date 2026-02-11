package gpt_sovits_go_sdk

// 提供了对GPT-SoVITS API的完整封装

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 代表 GPT-SoVITS API 客户端
type Client struct {
	BaseURL    string       // API基础URL
	HTTPClient *http.Client // HTTP客户端
}

// TTSRequest 代表 TTS 请求载荷
type TTSRequest struct {
	Text          string `json:"text"`           // str.(必填) 需要合成的文本
	TextLang      string `json:"text_lang"`      // str.(必填) 待合成文本的语言
	RefAudioPath  string `json:"ref_audio_path"` // str.(必填) 参考音频路径
	PromptText    string `json:"prompt_text"`    // str.(可选) 参考音频的提示文本
	PromptLang    string `json:"prompt_lang"`    // str.(必填) 参考音频提示文本的语言
	MediaType     string `json:"media_type"`     // str. 输出音频媒体类型，支持 "wav", "raw", "ogg", "aac"
	StreamingMode bool   `json:"streaming_mode"` // bool. 是否返回流式响应
}

// TTSResponse 代表 TTS 响应
type TTSResponse struct {
	StatusCode int    // HTTP状态码
	AudioData  []byte // 音频数据
	Error      error  // 错误信息
}

// ControlRequest 代表控制请求载荷
type ControlRequest struct {
	Command string `json:"command"` // "restart" 或 "exit"
}

// SetWeightsRequest 代表模型权重更新请求
type SetWeightsRequest struct {
	WeightsPath string `json:"weights_path"` // 权重文件路径
}

// NewClient 创建一个新的 GPT-SoVITS API 客户端
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second, // 根据需要调整超时时间
		},
	}
}

// TTS 发送文本转语音请求并返回音频响应
func (c *Client) TTS(ctx context.Context, req TTSRequest) (*TTSResponse, error) {
	// 将请求序列化为JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return &TTSResponse{Error: fmt.Errorf("请求序列化失败: %w", err)}, nil
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/tts", c.BaseURL)
	// 创建带上下文的HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &TTSResponse{Error: fmt.Errorf("创建请求失败: %w", err)}, nil
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return &TTSResponse{Error: fmt.Errorf("请求失败: %w", err)}, nil
	}
	defer resp.Body.Close()

	// 读取响应体
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return &TTSResponse{Error: fmt.Errorf("读取响应体失败: %w", err)}, nil
	}

	return &TTSResponse{
		StatusCode: resp.StatusCode,
		AudioData:  audioData,
	}, nil
}

// TTSSimple 提供简化的 TTS GET 接口
func (c *Client) TTSSimple(ctx context.Context, text, textLang, refAudioPath, promptLang, promptText string) (*TTSResponse, error) {
	// 构建请求对象
	req := TTSRequest{
		Text:          text,
		TextLang:      textLang,
		RefAudioPath:  refAudioPath,
		PromptLang:    promptLang,
		PromptText:    promptText,
		MediaType:     "wav",
		StreamingMode: false,
	}

	return c.TTS(ctx, req)
}

// Control 向服务器发送控制命令
func (c *Client) Control(ctx context.Context, command string) error {
	// 创建控制请求对象
	controlReq := ControlRequest{Command: command}

	// 序列化请求
	jsonData, err := json.Marshal(controlReq)
	if err != nil {
		return fmt.Errorf("控制请求序列化失败: %w", err)
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/control", c.BaseURL)
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建控制请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("控制请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("控制请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SetGPTWeights 更新 GPT 模型权重
func (c *Client) SetGPTWeights(ctx context.Context, weightsPath string) error {
	// 创建权重请求对象
	weightsReq := SetWeightsRequest{WeightsPath: weightsPath}

	// 序列化请求
	jsonData, err := json.Marshal(weightsReq)
	if err != nil {
		return fmt.Errorf("权重请求序列化失败: %w", err)
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/set_gpt_weights", c.BaseURL)
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建权重请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("设置GPT权重请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("设置GPT权重失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SetSoVITSWeights 更新 SoVITS 模型权重
func (c *Client) SetSoVITSWeights(ctx context.Context, weightsPath string) error {
	// 创建权重请求对象
	weightsReq := SetWeightsRequest{WeightsPath: weightsPath}

	// 序列化请求
	jsonData, err := json.Marshal(weightsReq)
	if err != nil {
		return fmt.Errorf("权重请求序列化失败: %w", err)
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/set_sovits_weights", c.BaseURL)
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建权重请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("设置SoVITS权重请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("设置SoVITS权重失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetTTSWithURLParams 使用URL参数提供 GET 接口
func (c *Client) GetTTSWithURLParams(ctx context.Context, params map[string]string) (*TTSResponse, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s/tts", c.BaseURL)

	// 构建查询参数
	query := "?"
	for key, value := range params {
		query += fmt.Sprintf("%s=%s&", key, value)
	}

	// 移除末尾的 &
	if len(query) > 1 {
		query = query[:len(query)-1]
	}

	url += query

	// 创建GET请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &TTSResponse{Error: fmt.Errorf("创建请求失败: %w", err)}, nil
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return &TTSResponse{Error: fmt.Errorf("请求失败: %w", err)}, nil
	}
	defer resp.Body.Close()

	// 读取响应体
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return &TTSResponse{Error: fmt.Errorf("读取响应体失败: %w", err)}, nil
	}

	return &TTSResponse{
		StatusCode: resp.StatusCode,
		AudioData:  audioData,
	}, nil
}

// ControlWithGet 提供控制命令的 GET 接口
func (c *Client) ControlWithGet(ctx context.Context, command string) error {
	// 构建请求URL
	url := fmt.Sprintf("%s/control?command=%s", c.BaseURL, command)

	// 创建GET请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("创建控制请求失败: %w", err)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("控制请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("控制请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SetGPTWeightsWithGet 提供设置 GPT 权重的 GET 接口
func (c *Client) SetGPTWeightsWithGet(ctx context.Context, weightsPath string) error {
	// 构建请求URL
	url := fmt.Sprintf("%s/set_gpt_weights?weights_path=%s", c.BaseURL, weightsPath)

	// 创建GET请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("创建权重请求失败: %w", err)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("设置GPT权重请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("设置GPT权重失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SetSoVITSWeightsWithGet 提供设置 SoVITS 权重的 GET 接口
func (c *Client) SetSoVITSWeightsWithGet(ctx context.Context, weightsPath string) error {
	// 构建请求URL
	url := fmt.Sprintf("%s/set_sovits_weights?weights_path=%s", c.BaseURL, weightsPath)

	// 创建GET请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("创建权重请求失败: %w", err)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("设置SoVITS权重请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("设置SoVITS权重失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
