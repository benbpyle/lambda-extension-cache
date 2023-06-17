package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

type Model struct {
	Id       string `json:"id"`
	FieldOne string `json:"fieldOne"`
	FieldTwo string `json:"fieldTwo"`
}

func getModel(id string) (*Model, error) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:4000/%s", id), nil)
	c := &http.Client{}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, error := c.Do(request)
	if error != nil {
		return nil, error
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		logrus.Debug("Item not found by key")
		return nil, nil
	}

	resBody, _ := ioutil.ReadAll(response.Body)

	var model Model
	err := json.Unmarshal(resBody, &model)

	return &model, err
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	m, err := getModel(event.PathParameters["id"])
	status := 200
	var body string

	if err != nil {
		b, _ := json.Marshal(err)
		body = string(b)
		status = 404
	} else {
		b, _ := json.Marshal(m)
		body = string(b)
	}

	e := events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Method":      "OPTIONS,POST,GET,PUT,DELETE",
		},
	}

	return e, nil
}

func main() {
	lambda.Start(handler)
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
