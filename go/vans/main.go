package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

//State is the current state of the van
type State struct {
	VanID     string `json:"van_id"`
	Latitude  float64
	Longitude float64
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String("vantable"),
	}
	result, err := db.Scan(params)
	if err != nil {
		return serverError(err)
	}
	states := []State{}
	newStates := []State{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &states)
	for _, state := range states {
		newState := State{
			VanID:     state.VanID,
			Latitude:  state.Latitude + (rand.Float64()-.5)*.1,
			Longitude: state.Longitude + (rand.Float64()-.5)*.1,
		}
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":lt": {
					N: aws.String(fmt.Sprintf("%f", newState.Latitude)),
				},
				":lg": {
					N: aws.String(fmt.Sprintf("%f", newState.Longitude)),
				},
			},
			TableName: aws.String("vantable"),
			Key: map[string]*dynamodb.AttributeValue{
				"van_id": {
					S: aws.String(newState.VanID),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set Latitude = :lt, Longitude = :lg"),
		}
		_, err := db.UpdateItem(input)
		if err != nil {
			return serverError(err)
		}
		newStates = append(newStates, newState)

	}
	js, err := json.Marshal(newStates)
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(show)
}
