package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/japb1998/tracking-demo/pkg/email"
	"github.com/joho/godotenv"
)

// logger
var logHandler = slog.NewTextHandler(os.Stdout, nil).WithAttrs([]slog.Attr{slog.String("pkg", "main")})
var logger = slog.New(logHandler)
var ginLambda *ginadapter.GinLambda

func main() {

	if os.Getenv("stage") == "local" {
		err := godotenv.Load()

		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("File does not exist.")
		}
	}
	r := gin.Default()
	fmt.Println("domain", os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	emailSvc := email.NewMailgunSvc(&email.MailgunOps{
		Domain: os.Getenv("MAILGUN_DOMAIN"),
		ApiKey: os.Getenv("MAILGUN_API_KEY"),
	})

	r.POST("/email", func(c *gin.Context) {
		logger.Info("sending email.")
		// we ignore the error
		u, _ := url.Parse(os.Getenv("API_URL"))
		u.Path = u.Path + "/pixel"
		query := u.Query()
		query.Add("email-id", "123")
		u.RawQuery = query.Encode()

		logger.Info("url", slog.String("url", u.String()))
		h := fmt.Sprintf(`<img src="%s"> 
		<div style="position:flex; justify-content:center; align-items:center; flex-direction: column; gap: 25px;">
			<h1>Hi, this is a test email</h1>
			<p>Thanks for opening this email</p>
		</div>`, u.String())

		e := email.NewEmail("", h, "Test Email", "<test-email>", nil, []string{"japb1998@gmail.com"}, nil)

		if err := emailSvc.Send(c.Request.Context(), e); err != nil {
			logger.Error("failed to send email", slog.String("error", err.Error()))
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte("failed to send email"))
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.WriteString("successfully send email")
	})

	if os.Getenv("stage") == "local" {
		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	} else {
		ginLambda = ginadapter.New(r)
		lambda.Start(handler)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}
