package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"
	"surf-hotel-engine/models"
)

var client = resty.New()

func FindMinAndMax(hotels []models.Hotel) (min int, max int) {
	min = hotels[0].MaxPrice
	max = hotels[0].MaxPrice

	for _, value := range hotels {
		if value.MaxPrice < min {
			min = value.MaxPrice
		}
		if value.MaxPrice > max {
			max = value.MaxPrice
		}
	}
	return min, max
}

func AddPhotosToSingleHotel(hotels *models.HotelResponse) {
	var hotelsIds []string

	for _, v := range hotels.Result {
		hotelsIds = append(hotelsIds, strconv.Itoa(v.Id))
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"id": strings.Join(hotelsIds, ","),
		}).
		Get("https://yasen.hotellook.com/photos/hotel_photos")

	if err != nil {
		fmt.Println(err.Error())
	}

	var photos map[string][]int
	err = json.Unmarshal(resp.Body(), &photos)

	for i := 0; i < len(hotels.Result); i++ {
		id := strconv.Itoa(hotels.Result[i].Id)
		id1, _ := strconv.Atoi(id)

		if hotels.Result[i].Id == id1 {
			hotels.Result[i].PhotoHotel = photos[id]
		}
	}
}

func AddPhotosToHotelResponse(hotels *models.PaginationHotels) {
	var hotelsIds []string

	for _, v := range hotels.Result {
		hotelsIds = append(hotelsIds, strconv.Itoa(v.Id))
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"id": strings.Join(hotelsIds, ","),
		}).
		Get("https://yasen.hotellook.com/photos/hotel_photos")

	if err != nil {
		fmt.Println(err.Error())
	}

	var photos map[string][]int
	err = json.Unmarshal(resp.Body(), &photos)

	for i := 0; i < len(hotels.Result); i++ {
		id := strconv.Itoa(hotels.Result[i].Id)
		id1, _ := strconv.Atoi(id)

		if hotels.Result[i].Id == id1 {
			hotels.Result[i].PhotoHotel = photos[id]
		}
	}
}

func IntArrayInclude(arr []int, expected int) bool {
	for _, v := range arr {
		if v == expected {
			return true
		}
	}

	return false
}
