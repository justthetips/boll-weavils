package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

//Point is a lat/long point
type Point struct {
	Latitude  string `json:"lattitude"`
	Longitude string `json:"longitude"`
}

//Place is a particular place
type Place struct {
	ID          string `json:"place_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	LatLng      Point
}

//User is the user in the application
type User struct {
	UserID          string `json:"user_id"`
	UserName        string `json:"username"`
	UserHome        string `json:"userhome"`
	UserDestination string `json:"userdestination"`
	UserRoute       string `json:"userroute"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	user := new(User)
	err := json.Unmarshal([]byte(req.Body), user)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	log.Println(user)
	id, err := createUser(user)
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/users?user_id=%s", id)},
	}, nil

}

func createUser(user *User) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("user_table"),
		Item: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(id.String()),
			},
			"username": {
				S: aws.String(user.UserName),
			},
			"userhome": {
				S: aws.String(user.UserHome),
			},
			"userdestination": {
				S: aws.String(user.UserDestination),
			},
			"userroute": {
				S: aws.String(user.UserRoute),
			},
		},
	}
	_, err = db.PutItem(input)
	return id.String(), err

}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["user_number"]
	if id == "" {
		return clientError(http.StatusBadRequest)
	}
	rt, err := getUser(id)
	if err != nil {
		return serverError(err)
	}

	js, err := json.Marshal(rt)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil

}

type dynamoUser struct {
	UserID          string `dynamodbav:"user_id"`
	UserName        string `dynamodbav:"username"`
	UserHome        string `dynamodbav:"userhome"`
	UserDestination string `dynamodbav:"userdestination"`
	Route           string `dynamodbav:"userroute"`
}

func getUser(userID string) (*User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("user_table"),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userID),
			},
		},
	}
	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	log.Println(result)
	du := dynamoUser{}
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &du)
	user := User{
		UserID:          du.UserID,
		UserName:        du.UserName,
		UserRoute:       du.Route,
		UserHome:        du.UserHome,
		UserDestination: du.UserDestination,
	}
	return &user, nil

}

type dynamoPlace struct {
	PlaceID     string `dynamodbav:"point_id"`
	Title       string `dynamodbav:"title"`
	Description string `dynamodbav:"description"`
	Latitude    string `dynamodbav:"lattitude"`
	Longitude   string `dynamodbav:"longitude"`
}

func getPlace(placeID string) *Place {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("point_table"),
		Key: map[string]*dynamodb.AttributeValue{
			"point_id": {
				S: aws.String(placeID),
			},
		},
	}
	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	dp := dynamoPlace{}
	if err != nil {
		panic(err)
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.

	err = dynamodbattribute.UnmarshalMap(result.Item, &dp)

	pt := Point{
		Latitude:  dp.Latitude,
		Longitude: dp.Longitude,
	}

	pl := &Place{
		ID:          dp.PlaceID,
		Description: dp.Description,
		Title:       dp.Title,
		LatLng:      pt,
	}
	return pl
}

func createPlace(place *Place) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("point_table"),
		Item: map[string]*dynamodb.AttributeValue{
			"place_id": {
				S: aws.String(id.String()),
			},
			"description": {
				S: aws.String(place.Description),
			},
			"title": {
				S: aws.String(place.Title),
			},
			"lattitude": {
				S: aws.String(place.LatLng.Latitude),
			},
			"longitude": {
				S: aws.String(place.LatLng.Longitude),
			},
		},
	}
	_, err = db.PutItem(input)
	return id.String(), err
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
	lambda.Start(router)
}
