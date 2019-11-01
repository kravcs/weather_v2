package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kravcs/gogo/cache"
	config "github.com/kravcs/gogo/config"
	"github.com/kravcs/gogo/models"
)

var Storage cache.Storage

func GetWeatherHandler(w http.ResponseWriter, r *http.Request) {

	var cities []models.CityWeather

	ch := make(chan models.CityWeather)

	vars := r.URL.Query()
	city := vars.Get("city")
	s := strings.Split(city, ",")
	baseUrl := config.Configuration.API.Endpoint + "?APPID=" + config.Configuration.API.Apikey
	duration := config.Configuration.Cache.Duration

	for _, city := range s {
		cityUrl := baseUrl + "&q=" + city
		go getContent(cityUrl, ch, duration)
	}

	for {
		result := <-ch
		cities = append(cities, result)

		if len(cities) == len(s) {
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(cities)
}

func getContent(url string, ch chan<- models.CityWeather, duration string) {
	var cityObj models.CityWeather

	content := Storage.Get(url)
	if len(content) != 0 {
		fmt.Println("Cache Hit!")
		json.Unmarshal(content, &cityObj)
		ch <- cityObj
	} else {

		err := getUrlResponse(url, &cityObj)
		if err != nil {
			log.Printf("Get url %s error: %v", url, err)
		}

		if d, err := time.ParseDuration(duration); err == nil {
			fmt.Printf("New page cached: %s for %s\n", url, duration)
			content, _ := json.Marshal(cityObj)
			Storage.Set(url, content, d)
		} else {
			fmt.Printf("Page not cached. err: %s\n", err)
		}

		ch <- cityObj
	}
}

func getUrlResponse(url string, this *models.CityWeather) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(this)
}
