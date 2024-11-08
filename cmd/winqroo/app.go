package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"winqroo/config"
	"winqroo/routes"

	"github.com/aws/aws-sdk-go-v2/aws"
	configForAWS "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	// "go.uber.org/zap"
)

type App struct {
	config       *config.Config
	router       http.Handler
	ddb          *dynamodb.Client
	sesClient    *ses.Client
	lambdaClient *lambda.Client
	// logger       *zap.Logger
}

func NewApp(config *config.Config) *App {
	awsConfig, err := configForAWS.LoadDefaultConfig(context.TODO(),
		// Uncomment the below options for dev/local environment
		// func(options *configForAWS.LoadOptions) error {
		// 	options.Region = config.AwsConfig.Region
		// 	options.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		// 		config.AwsConfig.AccessKey,
		// 		config.AwsConfig.SecretKey,
		// 		"",
		// 	))
		// 	return nil
		// },
	)
	if err != nil {
		fmt.Errorf("unable to load SDK config, %v", err)
		return nil
	}

	app := &App{
		config:       config,
		ddb:          dynamodb.NewFromConfig(awsConfig),
		sesClient:    ses.NewFromConfig(awsConfig),
		lambdaClient: lambda.NewFromConfig(awsConfig),
	}

	// Load all the registered routes.
	app.loadRoutes()

	return app
}

func (a *App) loadRoutes() {
	a.router = routes.NewRoutes(a.ddb, a.sesClient)
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", a.config.ServerPort),
		Handler: a.router,
	}

	// Check DynamoDB Connectivity
	_, err := a.ddb.ListTables(ctx, &dynamodb.ListTablesInput{Limit: aws.Int32(1)})
	if err != nil {
		return fmt.Errorf("failed to connect to dynamodb: %w", err)
	}

	// scripts.CreateDynamoDBTables(a.config) // For Development, init once!

	fmt.Println("Starting server at port: ", a.config.ServerPort)

	ch := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}

		close(ch)
	}()

	select {
	case err = <-ch:
		return err

	case <-ctx.Done():
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(ctxWithTimeout)
	}
}
