package chat

import (
	"archetype/app/shared/infrastructure/gemini"
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type IChat interface {
	SendMessage(ctx context.Context, domain interface{}) (interface{}, error)
}

type chatStruct struct {
	model gemini.Gemini1Dot0ProModelWrapper
}

func init() {
	ioc.Registry(NewChat, gemini.NewGemini1Dot0ProModelWrapper)
}
func NewChat(model gemini.Gemini1Dot0ProModelWrapper) IChat {
	return chatStruct{
		model: model,
	}
}

func (s chatStruct) SendMessage(ctx context.Context, domain interface{}) (interface{}, error) {
	// design your custom prompt and chat here using gemini-1.0-pro model
	return s.model.EphemeralChatExpectJSONResult(ctx, `Write a JSON object on a single line with the following properties:

	* A random number between 1 and 100 for the "age" property.
	* A random string of 10 characters for the "name" property (avoid special characters if possible). 
	* A random boolean value for the "active" property.
	* A list of 3 random numbers between 1 and 50 for the "scores" property.
	
	Here's a basic example format:
	{ "age": 52, "name": "JohnDoe", "active": true, "scores": [12, 45, 23] }`)
}

/*
func Instance() IChat {
	return ioc.Get[IChat](NewChat)
}
*/
