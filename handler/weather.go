package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/kravcs/weather_v2/cache"
	"github.com/kravcs/weather_v2/model"
)

type Data struct {
	Result model.CityWeather
	Err    error
}

type WeatherHandler struct {
	APIEnpoint    string
	APIKey        string
	CacheDuration int
	Storage       cache.Storage
}

// GetWeatherHandler provides weather forecast for one or several cities
func (wh WeatherHandler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	log.Println("handler started")
	defer log.Println("handler ended")

	cityParam := r.FormValue("city")
	if cityParam == "" {
		return NewStatusError(http.StatusBadRequest, "Query parameter \"city\" is missing")
	}
	citiesToForecast := strings.Split(cityParam, ",")

	cityChannel := make(chan Data, len(citiesToForecast))

	for _, city := range citiesToForecast {
		go wh.getForecast(ctx, city, cityChannel)
	}

	var cities []model.CityWeather

	for {
		select {
		case data := <-cityChannel:
			if data.Err != nil {
				return NewStatusError(http.StatusBadRequest, data.Err.Error())
			}

			cities = append(cities, data.Result)

			if len(cities) == len(citiesToForecast) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-type", "application/json")
				w.Write([]byte("Some weather output"))

				return nil
			}

		case <-ctx.Done():
			log.Println(ctx.Err().Error())
			return NewStatusError(http.StatusBadRequest, ctx.Err().Error())
		}
	}
}

func (wh WeatherHandler) getForecast(ctx context.Context, city string, cityChannel chan<- Data) {

	select {
	case <-ctx.Done():
		fmt.Println("(fail): Something happen")
		return
	default:
		var cityObj model.CityWeather

		// try to take data from chache
		cityForecast := wh.Storage.Get(city)
		if len(cityForecast) != 0 {
			json.Unmarshal(cityForecast, &cityObj)
			cityChannel <- Data{Result: cityObj}
			return
		}

		// try to take data from 3rd api
		err := wh.search(ctx, city, &cityObj)
		if err != nil {
			cityChannel <- Data{Err: err}
			return
		}

		// try to save content in cache
		content, _ := json.Marshal(cityObj)
		wh.Storage.Set(city, content, time.Duration(wh.CacheDuration))
		cityChannel <- Data{Result: cityObj}
		return
	}
}

func (wh WeatherHandler) search(ctx context.Context, city string, this *model.CityWeather) error {

	url := wh.APIEnpoint + "?APPID=" + wh.APIKey + "&q=" + city

	res, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "Error while GET url response from 3rd api")
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("Error StatusCode from 3rd api: " + string(body))
	}

	err = json.Unmarshal(body, this)
	if err != nil {
		return errors.Wrap(err, "Error response body from 3rd api")
	}

	// waiting for cancelation
	select {
	case <-ctx.Done():
		return nil
	}
}
