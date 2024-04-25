package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"subscribetool/src/pkg/utils/convert"
	"time"

	"github.com/olivere/elastic/v7"
)

type ELKConnect struct {
	URL      string
	UserName string
	Password string
	Index    string
	Stage    string
}

type ELK struct {
	client *elastic.Client
	index  string
	stage  string
}

var currentDate *time.Time

const DateFormat string = "2006-01-02"

const (
	mappings = ` 
	{
		"settings": {
		  "number_of_shards": 2,
		  "number_of_replicas": 1,
		  "index.mapping.total_fields.limit":1000000
		},
		"mappings": {
		  "properties": {
			"@timestamp": {
			  "type": "date"
			},
			"fields": {
			  "properties": {
				"endpoint": {
				  "type": "text",
				  "fields": {
					"keyword": {
					  "type": "keyword",
					  "ignore_above": 256
					}
				  }
				},
				"actionName": {
				  "type": "text",
				  "fields": {
					"keyword": {
					  "type": "keyword",
					  "ignore_above": 256
					}
				  }
				},
				"input": {
				  "type": "text",
				  "fields": {
					"keyword": {
					  "type": "keyword",
					  "ignore_above": 256
					}
				  }
				},
				"output": {
				  "type": "text",
				  "fields": {
					"keyword": {
					  "type": "keyword",
					  "ignore_above": 256
					}
				  }
				}
			  }
			},
			"level": {
			  "type": "text",
			  "fields": {
				"keyword": {
				  "type": "keyword",
				  "ignore_above": 256
				}
			  }
			},
			"message": {
			  "type": "text",
			  "fields": {
				"keyword": {
				  "type": "keyword",
				  "ignore_above": 256
				}
			  }
			}
		  }
		}
	  }`
)

type Content struct {
	Timestamp  time.Time `json:"@timestamp"`
	Endpoint   string    `json:"endpoint"`
	ActionName string    `json:"actionName"`
	Input      string    `json:"input"`
	Output     string    `json:"output"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
}

func NewELK(config ELKConnect) (*ELK, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(config.URL),
		elastic.SetBasicAuth(config.UserName, config.Password),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}

	return &ELK{
		client: client,
		index:  config.Index,
		stage:  config.Stage,
	}, nil
}

func (e *ELK) Error(endpoint, actionName string, request, message interface{}) {
	nowDate := time.Now().Format(DateFormat)
	var stage string
	if e.stage == "dev" {
		stage = "-dev"
	} else if e.stage == "uat" {
		stage = "-uat"
	}
	index := e.index + stage + "-" + nowDate
	ctx := context.Background()
	exist, existError := e.client.IndexExists(index).Do(ctx)
	if existError != nil {
		fmt.Println("ELK:existError", existError)
		return
	}

	if !exist {
		_, createErr := e.client.CreateIndex(index).Do(ctx)
		if createErr != nil {
			fmt.Println("ELK:createErr", createErr)
			return
		}
	}

	strRequest := ""
	if request != nil {
		jsonRequest, _ := json.Marshal(request)
		strRequest = string(jsonRequest)
	}

	strMessage := ""
	if message != nil {
		jsonMessage, _ := json.Marshal(message)
		strMessage = string(jsonMessage)
	}

	content := Content{
		Timestamp:  time.Now(),
		Endpoint:   endpoint,
		ActionName: actionName,
		Level:      "Error",
		Message:    strMessage,
		Input:      strRequest,
		Output:     "",
	}

	jsonContent, _ := json.Marshal(content)
	strContent := string(jsonContent)
	e.client.Index().
		Index(index).
		BodyJson(strContent).
		Do(ctx)

	if currentDate == nil ||
		(currentDate != nil && currentDate.Format(DateFormat) != nowDate) {
		err := e.FlushIndex()
		if err != nil {
			return
		}
		currentDate = convert.ValueToTimePointer(time.Now())
	}
}

func (e *ELK) FlushIndex() error {
	previousDate := time.Now().AddDate(0, 0, -90)
	formatDate := previousDate.Format(DateFormat)
	index := e.index + "-" + formatDate
	ctx := context.Background()

	_, existError := e.client.IndexExists(index).Do(ctx)
	if existError != nil {
		return existError
	}
	_, deleteErr := e.client.DeleteIndex(index).Do(ctx)

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
