package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/gif"
	"os"
	"time"

	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var logHandler = slog.NewTextHandler(os.Stdout, nil).WithAttrs([]slog.Attr{slog.String("pkg", "main")})
var logger = slog.New(logHandler)

func main() {
	ctx, c := context.WithTimeout(context.Background(), time.Second*30)
	defer c()

	lambda.StartWithOptions(HandleRequest, lambda.WithContext(ctx))
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger.Info("An Email Was Opened")
	if v, ok := event.Headers["user-agent"]; ok {
		logger.Info("From:", slog.String("user-agent", v))
	}

	// Email ID
	logger.Info("Email ID:", slog.String("email-id", event.QueryStringParameters["email-id"]))
	logger.Info("path:", slog.String("path", event.Path))

	// 1 x 1 pixel
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	var buff bytes.Buffer

	if err := gif.Encode(&buff, img, nil); err != nil {

		logger.Error("failed to encode pixel.", slog.String("error", err.Error()))

		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			IsBase64Encoded: false,
		}, nil
	}
	// reason to send image as base64
	/*
		API Gateway typically expects data in a text format, and base64 encoding is a convenient way to represent binary data
	*/
	// out 1x1 pixel as a base64 string.
	response := base64.StdEncoding.EncodeToString(buff.Bytes())

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "image/gif",
		},
		Body:            response,
		IsBase64Encoded: true,
	}, nil
}
