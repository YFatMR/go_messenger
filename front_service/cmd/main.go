package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Info struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

//func printFeature(client pb. , point *pb.Point) {
//log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//defer cancel()
//feature, err := client.GetFeature(ctx, point)
//if err != nil {
//	log.Fatalf("client.GetFeature failed: %v", err)
//}
//log.Println(feature)
//}

func getInfo(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id := params["id"]
	info := Info{id, "Some info"}
	json.NewEncoder(responseWriter).Encode(info)
}

func main() {
	fmt.Println("Service started")
	router := mux.NewRouter()
	router.HandleFunc("/info/{id}", getInfo).Methods("GET")
	http.ListenAndServe(":8000", router)
	fmt.Println("Service finished")
}
