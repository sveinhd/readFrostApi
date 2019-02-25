package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var s = flag.String("s", "SN50540", "ids")

//https://mholt.github.io/json-to-go/
type frostData struct {
	Context          string    `json:"@context"`
	Type             string    `json:"@type"`
	APIVersion       string    `json:"apiVersion"`
	License          string    `json:"license"`
	CreatedAt        time.Time `json:"createdAt"`
	QueryTime        float64   `json:"queryTime"`
	CurrentItemCount int       `json:"currentItemCount"`
	ItemsPerPage     int       `json:"itemsPerPage"`
	Offset           int       `json:"offset"`
	TotalItemCount   int       `json:"totalItemCount"`
	CurrentLink      string    `json:"currentLink"`
	Data             []struct {
		Type        string `json:"@type"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ShortName   string `json:"shortName"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		WmoID       int    `json:"wmoId"`
		Geometry    struct {
			Type        string    `json:"@type"`
			Coordinates []float64 `json:"coordinates"`
			Nearest     bool      `json:"nearest"`
		} `json:"geometry"`
		Masl           int      `json:"masl"`
		ValidFrom      string   `json:"validFrom"`
		County         string   `json:"county"`
		CountyID       int      `json:"countyId"`
		Municipality   string   `json:"municipality"`
		MunicipalityID int      `json:"municipalityId"`
		StationHolders []string `json:"stationHolders"`
		ExternalIds    []string `json:"externalIds"`
		WigosID        string   `json:"wigosId"`
	} `json:"data"`
}

func main() {

	flag.Parse()

	//	fmt.Println("Reading Environment Variable")
	var clientid string
	var idIsSet bool

	clientid, idIsSet = os.LookupEnv("CLIENTID")
	if !idIsSet {
		log.Fatal("CLIENTIDI not set")
	}

	client := &http.Client{}
	url := "https://frost.met.no/sources/v0.jsonld?ids=" + *s
	req, err := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(clientid, "")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var d frostData
	jsonErr := json.Unmarshal(jsonBody, &d)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for i := 0; i < len(d.Data); i++ {

		fmt.Println("ID: " + d.Data[i].ID)
		fmt.Println("Name: " + d.Data[i].Name)
		fmt.Println("Geometry Type: " + d.Data[i].Geometry.Type)

		fmt.Printf("Lat: %f Long: %f \n", d.Data[i].Geometry.Coordinates[1], d.Data[i].Geometry.Coordinates[0])

		fmt.Println("Municipality: " + d.Data[i].Municipality)
		fmt.Println("Country: " + d.Data[i].Country)

		fmt.Print("ExternalIds: ")
		for j := 0; j < len(d.Data[i].ExternalIds); j++ {
			fmt.Print(d.Data[i].ExternalIds[j] + ", ")
		}
		fmt.Println()
	}
	fmt.Println()
}
