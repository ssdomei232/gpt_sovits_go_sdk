# gpt-sovits-go-sdk

```go
package main

import (
 "context"
 "fmt"
 "os"
)

func main() {
 // 创建客户端实例
 client := NewClient("http://127.0.0.1:9880")

 ctx := context.Background()

 // 示例1: 使用POST请求进行TTS转换
 ttsReq := TTSRequest{
  Text:            "你好，这是一个测试文本",     // 需要合成的文本
  TextLang:        "zh",                      // 文本语言
  RefAudioPath:    "reference_audio.wav",     // 参考音频路径
  PromptText:      "这是一个示例提示文本",      // 提示文本
  PromptLang:      "zh",                      // 提示文本语言
  TextSplitMethod: "cut5",                    // 文本分割方法
  BatchSize:       1,                         // 批次大小
  MediaType:       "wav",                     // 输出音频类型
  StreamingMode:   false,                     // 是否流式输出
 }

 resp, err := client.TTS(ctx, ttsReq)
 if err != nil {
  fmt.Printf("TTS请求失败: %v\n", err)
  return
 }

 if resp.StatusCode == 200 {
  // 保存音频文件
  err = os.WriteFile("output.wav", resp.AudioData, 0644)
  if err != nil {
   fmt.Printf("保存音频文件失败: %v\n", err)
   return
  }
  fmt.Println("音频生成成功并已保存为 output.wav")
 } else {
  fmt.Printf("TTS请求失败，状态码: %d\n", resp.StatusCode)
 }

 // 示例2: 使用简化接口
 simpleResp, err := client.TTSSimple(ctx, 
  "这是简化接口的测试文本", 
  "zh", 
  "reference_audio.wav", 
  "zh", 
  "提示文本")
 if err != nil {
  fmt.Printf("简化TTS请求失败: %v\n", err)
  return
 }

 if simpleResp.StatusCode == 200 {
  err = os.WriteFile("simple_output.wav", simpleResp.AudioData, 0644)
  if err != nil {
   fmt.Printf("保存简单音频文件失败: %v\n", err)
   return
  }
  fmt.Println("简单音频生成成功并已保存为 simple_output.wav")
 }

 // 示例3: 控制命令
 err = client.Control(ctx, "restart") // 或 "exit"
 if err != nil {
  fmt.Printf("控制命令执行失败: %v\n", err)
 } else {
  fmt.Println("控制命令执行成功")
 }

 // 示例4: 更新模型权重
 err = client.SetGPTWeights(ctx, "GPT_SoVITS/pretrained_models/custom_gpt_model.ckpt")
 if err != nil {
  fmt.Printf("更新GPT权重失败: %v\n", err)
 } else {
  fmt.Println("GPT权重更新成功")
 }

 err = client.SetSoVITSWeights(ctx, "GPT_SoVITS/pretrained_models/custom_sovits_model.pth")
 if err != nil {
  fmt.Printf("更新SoVITS权重失败: %v\n", err)
 } else {
  fmt.Println("SoVITS权重更新成功")
 }

 // 示例5: 使用GET接口设置GPT权重
 err = client.SetGPTWeightsWithGet(ctx, "GPT_SoVITS/pretrained_models/custom_gpt_model.ckpt")
 if err != nil {
  fmt.Printf("通过GET接口更新GPT权重失败: %v\n", err)
 } else {
  fmt.Println("通过GET接口GPT权重更新成功")
 }
}
```
