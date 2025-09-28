package sinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"loggingLibGo/message"
	"net/http"
)

type ElasticSearch struct {
	URL    string
	Index  string
	client *http.Client
}

type SearchQuery struct {
	Level     string
	Namespace string
	StartDate string
	EndDate   string
}

func NewElasticSink(url, index string) *ElasticSearch {
	return &ElasticSearch{
		URL:    url,
		Index:  index,
		client: &http.Client{},
	}
}

func (e *ElasticSearch) Write(msg message.Message) error {
	//fmt.Println("ElasticSink.Write() called", msg)
	// key value pairs for elastic search
	doc := map[string]any{
		"timestamp": msg.Timestamp,
		"level":     msg.Level.String(),
		"namespace": msg.Namespace,
		"content":   msg.Content,
	}

	// convert to json
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// create the request
	endPoint := fmt.Sprintf("%s/%s/_doc", e.URL, e.Index)

	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		fmt.Println("Error :", err)
		return err
	}

	defer resp.Body.Close() // close response body

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return fmt.Errorf("elastic search returned status code %d", resp.StatusCode)
	}

	return nil
}

func (e *ElasticSearch) Search(query SearchQuery) ([]map[string]any, error) {
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
	}

	mustClauses := esQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})

	if query.Level != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"level": query.Level,
			},
		})
	}

	if query.Namespace != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"namespace": query.Namespace,
			},
		})
	}

	if query.StartDate != "" && query.EndDate != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": query.StartDate,
					"lte": query.EndDate,
				},
			},
		})
	}

	esQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustClauses

	jsonData, err := json.Marshal(esQuery)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/_search", e.URL, e.Index)
	req, err := http.NewRequest("GET", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	hitsArray := result["hits"].(map[string]interface{})["hits"].([]interface{})

	var logs []map[string]interface{}

	for _, hit := range hitsArray {
		src := hit.(map[string]interface{})["_source"]
		logs = append(logs, src.(map[string]interface{}))
	}

	return logs, nil

}
