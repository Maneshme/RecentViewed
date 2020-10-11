package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	es "recentViewed/elasticsearch"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Recently Viewed Items API!")
	log.Println("Endpoint Hit: homePage")
}

func getRecentViewed(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: return all recently viewed items")
	queryParams := r.URL.Query()
	userName := queryParams["userName"][0]
	from, err := strconv.Atoi(queryParams["from"][0])
	limit, lerr := strconv.Atoi(queryParams["limit"][0])
	if err != nil || lerr != nil {
		fmt.Fprintf(w, "Invalid limit or from")
	}
	client := es.CreateClient()
	err = client.ConnectClient()
	if err != nil {
		fmt.Fprintf(w, "Error creating elastic search client")
	}
	productsViewed, err := client.GetAllDataByLimit(userName, from, limit)
	if err != nil {
		fmt.Fprintf(w, "Error getting data from elastic search")
	}
	client.CloseConnection(userName)
	json.NewEncoder(w).Encode(productsViewed)
}

func postRecentViewed(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: Create a recently viewed item in elasticsearch")
	queryParams := r.URL.Query()
	userName := queryParams["userName"][0]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var product es.ProductViewed
	json.Unmarshal(reqBody, &product)
	client := es.CreateClient()
	err := client.ConnectClient()
	if err != nil {
		fmt.Fprintf(w, "Error creating elastic search client")
	}
	err = client.PostData(userName, product)
	if err != nil {
		fmt.Fprintf(w, "Error posting data from elastic search")
	}
	client.CloseConnection(userName)
	json.NewEncoder(w).Encode(product)
}
