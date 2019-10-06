package main

import (
	"encoding/json"
	"log"
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

//Point is a lat/long point
type Point struct {
	Latitude  string
	Longitude string
}

//Place is a particular place
type Place struct {
	ID          string
	Title       string
	Description string
	LatLng      Point
}

//Route is a collection of places
type Route struct {
	ID    string
	Name  string
	Stops []*Place
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["route_number"]
	if id == "" {
		return clientError(http.StatusBadRequest)
	}
	rt, err := getRoute(id)
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

type dynamoRoute struct {
	RouteID   string `dynamodbav:"route_id"`
	RouteName string `dynamodbav:"route_name"`
	Points    []string
}

type dynamoPlace struct {
	PlaceID     string `dynamodbav:"point_id"`
	Title       string `dynamodbav:"title"`
	Description string `dynamodbav:"description"`
	Latitude    string `dynamodbav:"lattitude"`
	Longitude   string `dynamodbav:"longitude"`
}

func getRoute(id string) (*Route, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("route_table"),
		Key: map[string]*dynamodb.AttributeValue{
			"route_id": {
				S: aws.String(id),
			},
		},
	}
	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	if err != nil {
		panic(err)
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.
	rt := new(dynamoRoute)
	err = dynamodbattribute.UnmarshalMap(result.Item, rt)

	route := &Route{
		ID:   rt.RouteID,
		Name: rt.RouteName,
	}

	var pnts []*Place

	for _, point := range rt.Points {
		p, _ := getPlace(point)
		pnts = append(pnts, p)
	}
	route.Stops = pnts
	return route, nil

}

func getPlace(id string) (*Place, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("point_table"),
		Key: map[string]*dynamodb.AttributeValue{
			"point_id": {
				S: aws.String(id),
			},
		},
	}
	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	dp := dynamoPlace{}
	if err != nil {
		return nil, err
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
	return pl, nil
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
