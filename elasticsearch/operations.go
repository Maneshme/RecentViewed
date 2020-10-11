package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	elastic "github.com/olivere/elastic/v7"
)

type ElasticClient struct {
	esClient *elastic.Client
}

type ProductViewed struct {
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	ViewedAt    string `json:"viewedAt"`
}

var (
	viewedMapping = `{
		"mappings": {
			"properties": { 
				"productId": {
					"type": "keyword"
			  	},
			  	"productName":{
					"type":"keyword"
				},
				"viewedAt":{
					"type":"keyword"
				}
			}
		}
	  }`
)

func CreateClient() ElasticClient {
	return ElasticClient{}
}

func (client *ElasticClient) ConnectClient() error {
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		log.Printf("Error: Connecting to elastic search %s", err.Error())
		return err
	}
	client.esClient = esClient
	return nil
}

func (client *ElasticClient) PostData(indexPrefix string, viewed ProductViewed) error {
	ctx := context.Background()
	dataJSON, err := json.Marshal(viewed)
	if err != nil {
		log.Printf("Error: Marshalling viewed data %s", err.Error())
		return err
	}
	js := string(dataJSON)
	indexName := fmt.Sprintf("%s_viewed", strings.ToLower(indexPrefix))
	exists, err := client.esClient.IndexExists(indexName).Do(ctx)
	if err != nil || !exists {
		// Create a new index.
		createIndex, err := client.esClient.CreateIndex(indexName).Body(viewedMapping).Do(ctx)
		if err != nil {
			log.Printf("Error: Creating Index %s", err.Error())
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	_, err = client.esClient.Index().
		Index(indexName).
		BodyJson(js).
		Do(ctx)

	if err != nil {
		log.Printf("Error: Posting data to elastic search %s", err.Error())
		return err
	}
	return nil
}

func (client *ElasticClient) GetAllDataByLimit(indexPrefix string, from, limit int) ([]ProductViewed, error) {
	indexName := fmt.Sprintf("%s_viewed", strings.ToLower(indexPrefix))
	var productsViewed []ProductViewed
	searchResult, err := client.esClient.Search().
		Index(indexName).
		Query(elastic.NewMatchAllQuery()).
		Sort("viewedAt", true).
		From(from).Size(limit).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		// Handle error
		return []ProductViewed{}, err
	}
	if hits := searchResult.TotalHits(); hits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var product ProductViewed
			err := json.Unmarshal(hit.Source, &product)
			if err != nil {
				return []ProductViewed{}, err
			}
			productsViewed = append(productsViewed, product)
		}
	}

	return productsViewed, nil
}

func (client *ElasticClient) CloseConnection(indexPrefix string) {
	indexName := fmt.Sprintf("%s_viewed", strings.ToLower(indexPrefix))
	_ = client.esClient.CloseIndex(indexName)
}
