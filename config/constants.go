package config

import "github.com/aws/aws-sdk-go-v2/aws"

var (
	StationTableName = aws.String("stations")
	TrainTableName   = aws.String("trains") 
)
