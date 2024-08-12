package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/generative-ai-go/genai"
)

type Gemini1Dot0ProModelWrapper struct {
	*genai.GenerativeModel
}

func init() {
	ioc.Registry(NewGemini1Dot0ProModelWrapper, newClient)
}
func NewGemini1Dot0ProModelWrapper(client *genai.Client) Gemini1Dot0ProModelWrapper {
	return Gemini1Dot0ProModelWrapper{
		GenerativeModel: client.GenerativeModel("gemini-1.0-pro"),
	}
}

func (s Gemini1Dot0ProModelWrapper) EphemeralChatExpectJSONResult(
	ctx context.Context,
	msg string) (map[string]interface{}, error) {
	res, err := s.GenerativeModel.StartChat().SendMessage(ctx, genai.Text(msg))
	if err != nil {
		return nil, err
	}
	return getJSONResponse(res)
}

func getJSONResponse(resp *genai.GenerateContentResponse) (map[string]interface{}, error) {
	var output strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {

			for _, part := range cand.Content.Parts {
				text := fmt.Sprintf("%s", part)
				output.WriteString(text)
			}
		}
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(output.String()), &data)
	return data, err
}
