package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trussworks/ecr-scan/pkg/ecrscan"
	"go.uber.org/zap"
)

var logger *zap.Logger

// This function is for establishing our session with AWS.
func makeECRClient(region, profile string) *ecr.ECR {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: aws.String(region),
		},
	}))
	ecrClient := ecr.New(sess)
	return ecrClient
}

func evaluateImage() (ecrscan.Report, error) {
	evaluator := ecrscan.Evaluator{
		MaxScanAge: viper.GetInt("maxScanAge"),
		Logger:     logger,
		ECRClient:  makeECRClient(viper.GetString("region"), viper.GetString("profile")),
	}
	target := ecrscan.Target{
		Repository: viper.GetString("repository"),
		ImageTag:   viper.GetString("tag"),
	}
	return evaluator.Evaluate(&target)
}

// HandleRequest is the root command's Lambda handler
func HandleRequest(ctx context.Context, target ecrscan.Target) (ecrscan.Report, error) {
	viper.Set("repository", target.Repository)
	viper.Set("tag", target.ImageTag)
	result, err := evaluateImage()
	if err != nil {
		logger.Error("Error evaluating target image", zap.Error(err))
	}
	logger.Info("Scan result", zap.Any("Report", result))
	return result, err
}

var rootCmd = &cobra.Command{
	Use:   "ecr-scan",
	Short: "ecr-scan is an application for analyzing ECR scan findings",
	Long:  "ecr-scan is an application for analyzing ECR scan findings",
	Run: func(cmd *cobra.Command, args []string) {
		// declare error so that call to zap.NewProduction() will use
		// logger declared above
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("could not initialize zap logger: %v", err)
		}

		if viper.GetBool("lambda") {
			lambda.Start(HandleRequest)
		} else {
			result, aerr := evaluateImage()
			if aerr != nil {
				logger.Error("Error evaluating target image", zap.Error(aerr))
			}
			logger.Info("Scan result", zap.Any("Report", result))
		}

		err = logger.Sync()
		if err != nil {
			log.Fatal("could not sync logger")
		}
	},
}

// Execute root command execute function
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()
	if err := viper.BindEnv("repository", "ECR_REPOSITORY"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("tag", "IMAGE_TAG"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("lambda", "LAMBDA"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("maxScanAge", "MAX_SCAN_AGE"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("profile", "AWS_PROFILE"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("region", "AWS_REGION"); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().StringP("repository", "r", "", "ECR repository where the image is located")
	rootCmd.Flags().StringP("tag", "t", "", "Image tag to retrieve findings for")
	rootCmd.Flags().Bool("lambda", false, "Run as Lambda function")
	rootCmd.Flags().IntP("maxScanAge", "m", 24, "Maximum allowed age for image scan (hours)")
	rootCmd.Flags().String("profile", "", "The AWS profile to use")
	rootCmd.Flags().String("region", "", "The AWS region to use")
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("lambda", false)
	viper.SetDefault("maxScanAge", 24)
}
