package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"surf-hotel-engine/models"
)

type Filters struct {
	Price struct {
		Included bool `json:"included"`
		MinValue int  `json:"minValue"`
		MaxValue int  `json:"maxValue"`
	} `json:"price"`
	Distance struct {
		Included bool    `json:"included"`
		MinValue float64 `json:"minValue"`
		MaxValue float64 `json:"maxValue"`
	} `json:"distance"`
	Stars struct {
		Included bool `json:"included"`
		Values   struct {
			Star1 bool `json:"star1"`
			Star2 bool `json:"star2"`
			Star3 bool `json:"star3"`
			Star4 bool `json:"star4"`
			Star5 bool `json:"star5"`
		} `json:"values"`
	} `json:"stars"`
	HousingType struct {
		Included bool   `json:"included"`
		Lang     string `json:"lang"`
		Values   struct {
			Hotels      bool `json:"hotels"`
			Apartments  bool `json:"apartments"`
			Hostels     bool `json:"hostels"`
			Villas      bool `json:"villas"`
			GuestHouses bool `json:"guestHouses"`
		} `json:"values"`
	} `json:"housingType"`
	FoodAndPayment struct {
		Breakfast      bool `json:"breakfast"`
		AllInclusive   bool `json:"allInclusive"`
		PaymentNow     bool `json:"paymentNow"`
		PaymentInHotel bool `json:"paymentInHotel"`
		FreeCancel     bool `json:"freeCancel"`
	} `json:"foodAndPayment"`
	RoomOptions struct {
		Bathroom        bool `json:"bathroom"`
		WiFi            bool `json:"wiFi"`
		AirConditioning bool `json:"airConditioning"`
		Hairdryer       bool `json:"hairdryer"`
		Tv              bool `json:"tv"`
		Safe            bool `json:"safe"`
	} `json:"roomOptions"`
	HotelOptions struct {
		Parking         bool `json:"parking"`
		Restaurant      bool `json:"restaurant"`
		SwimmingPool    bool `json:"swimmingPool"`
		Spa             bool `json:"spa"`
		PrivateBeach    bool `json:"privateBeach"`
		Fitness         bool `json:"fitness"`
		PetsAllowed     bool `json:"petsAllowed"`
		Safe            bool `json:"safe"`
		Laundry         bool `json:"laundry"`
		AllDayReception bool `json:"allDayReception"`
	} `json:"hotelOptions"`
}

func FilterHotelResponse(hotels *models.HotelResponse, bytesBody io.ReadCloser, search string) []models.Hotel {
	var rawBody, err = ioutil.ReadAll(bytesBody)
	if err != nil {
		fmt.Println(err.Error())
	}

	var filteredHotels []models.Hotel
	var filterList Filters
	fmt.Printf("filters from app --> %v\n", string(rawBody))

	err = json.Unmarshal(rawBody, &filterList)

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, v := range hotels.Result {
		if len(search) > 3 {
			if !strings.Contains(strings.ToLower(v.Name), strings.ToLower(search)) {
				continue
			}
		}

		if filterList.RoomOptions.Bathroom {
			if !BathroomIncluded(v) {
				continue
			}
		}

		if filterList.RoomOptions.Tv {
			if !TvIncluded(v) {
				continue
			}
		}

		if filterList.RoomOptions.Safe {
			if !SafeIncluded(v) {
				continue
			}
		}

		if filterList.RoomOptions.AirConditioning {
			if !AirConditioningIncluded(v) {
				continue
			}
		}

		if filterList.RoomOptions.WiFi {
			if !WiFiInRoomIncluded(v) {
				continue
			}
		}

		if filterList.RoomOptions.Hairdryer {
			if !HairdryerIncluded(v) {
				continue
			}
		}

		// hotel filter
		if filterList.HotelOptions.Parking {
			if !CarParkingIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.Restaurant {
			if !RestaurantIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.SwimmingPool {
			if !SwimmingPoolIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.Spa {
			if !SpaIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.PrivateBeach {
			if !PrivateBeachIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.Fitness {
			if !GymFitnessCentersIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.PetsAllowed {
			if !PetsAllowedIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.Laundry {
			if !LaundryIncluded(v) {
				continue
			}
		}

		if filterList.HotelOptions.AllDayReception {
			if !ReceptionAllDayIncluded(v) {
				continue
			}
		}

		// food and payment filter
		if filterList.FoodAndPayment.Breakfast {
			if !BreakfastIncluded(v) {
				continue
			}
		}

		if filterList.FoodAndPayment.AllInclusive {
			if !AllInclusiveIncluded(v) {
				continue
			}
		}

		if filterList.FoodAndPayment.FreeCancel {
			if !RefundableIncluded(v) {
				continue
			}
		}

		// price and distance filter
		if filterList.Price.Included {
			if !PriceBetween(v, filterList.Price.MinValue, filterList.Price.MaxValue) {
				continue
			}
		}

		if filterList.Distance.Included {
			if !DistanceBetween(v, filterList.Distance.MinValue, filterList.Distance.MaxValue) {
				continue
			}
		}

		if filterList.Stars.Included {
			if (Stars1(v) || Stars0(v)) && !filterList.Stars.Values.Star1 {
				continue
			}

			if Stars2(v) && !filterList.Stars.Values.Star2 {
				continue
			}

			if Stars3(v) && !filterList.Stars.Values.Star3 {
				continue
			}

			if Stars4(v) && !filterList.Stars.Values.Star4 {
				continue
			}

			if Stars5(v) && !filterList.Stars.Values.Star5 {
				continue
			}
		}

		if filterList.HousingType.Included {
			if !filterList.HousingType.Values.Hotels && HotelType(v) {
				continue
			}

			if !filterList.HousingType.Values.Hostels && HostelType(v) {
				continue
			}

			if !filterList.HousingType.Values.Apartments && ApartmentsType(v) {
				continue
			}

			if !filterList.HousingType.Values.Villas && VillasType(v) {
				continue
			}

			if !filterList.HousingType.Values.GuestHouses && GuestHousesType(v) {
				continue
			}
		}

		filteredHotels = append(filteredHotels, v)
	}

	//filteredHotelsMap := funk.ToMap(filteredHotels, "Id")
	//fmt.Println(filteredHotelsMap)
	//
	//for k, v := range filteredHotelsMap.(map[int]models.Hotel) {
	//if filterList.HousingType.Hotels {
	//	if HotelType(v) {
	//		continue
	//	} else {
	//		delete(filteredHotelsMap.(map[int]models.Hotel), k)
	//	}
	//}
	//
	//if filterList.HousingType.Hostels {
	//	if HostelType(v) {
	//		continue
	//	} else {
	//		delete(filteredHotelsMap.(map[int]models.Hotel), k)
	//	}
	//}
	//
	//if filterList.HousingType.Apartments {
	//	if !ApartmentsType(v) {
	//		delete(filteredHotelsMap.(map[int]models.Hotel), k)
	//		break
	//	}
	//}
	//
	//switch filter := filterList.HousingType; {
	//case filter.Hotels:
	//	if !HotelType(v) {
	//		delete(filteredHotelsMap.(map[int]models.Hotel), k)
	//	}
	//	continue
	//case filter.Apartments:
	//	if !ApartmentsType(v) {
	//		delete(filteredHotelsMap.(map[int]models.Hotel), k)
	//	}
	//	continue
	//case filter.Hostels:
	//	if !HostelType(v) {
	//		delete(filteredHotelsMap.(map[int]models.Hotel), k)
	//	}
	//}
	//}
	//
	//fmt.Println("map length", len(filteredHotelsMap.(map[int]models.Hotel)))
	//
	//fmt.Println(filteredHotelsMap.(map[int]models.Hotel))
	//
	//for i, v := range hotels.Result {
	//	if filterList.HousingType.Hotels {
	//		if HotelType(v) {
	//			if filterList.Stars.Star1 {
	//				if !Stars1(v) {
	//					remove(hotels.Result, i)
	//				}
	//			}
	//			if filterList.Stars.Star2 {
	//				if !Stars2(v) {
	//					remove(hotels.Result, i)
	//				}
	//			}
	//			if filterList.Stars.Star3 {
	//				if !Stars3(v) {
	//					remove(hotels.Result, i)
	//				}
	//			}
	//			if filterList.Stars.Star4 {
	//				if !Stars4(v) {
	//					remove(hotels.Result, i)
	//				}
	//			}
	//			if filterList.Stars.Star5 {
	//				if !Stars5(v) {
	//					remove(hotels.Result, i)
	//				}
	//			}
	//		}
	//	}
	//}
	//
	//for _, v := range hotels.Result {
	//	// housing filter
	//
	//	if filterList.HousingType.Hotels {
	//		if HotelType(v) {
	//			filteredHotels = append(filteredHotels, v)
	//		}
	//	}
	//
	//	if filterList.HousingType.Hostels {
	//		if HostelType(v) {
	//			filteredHotels = append(filteredHotels, v)
	//		}
	//	}
	//
	//	if filterList.HousingType.Apartments {
	//		if ApartmentsType(v) {
	//			filteredHotels = append(filteredHotels, v)
	//		}
	//	}
	//
	//	// stars filter
	//
	//	//if filterList.Stars.Star3 {
	//	//	if Stars3(v) {
	//	//		filteredHotels = append(filteredHotels, v)
	//	//	}
	//	//}
	//	//
	//	//if filterList.Stars.Star4 {
	//	//	if Stars4(v) {
	//	//		filteredHotels = append(filteredHotels, v)
	//	//	}
	//	//}
	//}
	//
	//if len(filteredHotels) > 0 {
	//	for i := 0; i < len(filteredHotels); i++ {
	//		if filterList.RoomOptions.Bathroom {
	//			if !BathroomIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.RoomOptions.Tv {
	//			if !TvIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.RoomOptions.Safe {
	//			if !SafeIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.RoomOptions.AirConditioning {
	//			if !AirConditioningIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.RoomOptions.WiFi {
	//			if !WiFiInRoomIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.RoomOptions.Hairdryer {
	//			if !HairdryerIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		// hotel filter
	//		if filterList.HotelOptions.Parking {
	//			if !CarParkingIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Restaurant {
	//			if !RestaurantIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.SwimmingPool {
	//			if !SwimmingPoolIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Spa {
	//			if !SpaIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.PrivateBeach {
	//			if !PrivateBeachIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Fitness {
	//			if !GymFitnessCentersIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.PetsAllowed {
	//			if !PetsAllowedIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Laundry {
	//			if !LaundryIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.HotelOptions.AllDayReception {
	//			if !ReceptionAllDayIncluded(filteredHotels[i]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		// food and payment filter
	//
	//		// price and distance filter
	//		if filterList.Price.Included {
	//			if !PriceBetween(filteredHotels[i], filterList.Price.Values[0], filterList.Price.Values[1]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//
	//		if filterList.Distance.Included {
	//			if !DistanceBetween(filteredHotels[i], filterList.Distance.Values[0], filterList.Distance.Values[1]) {
	//				remove(filteredHotels, i)
	//			}
	//		}
	//	}
	//} else {
	//	for _, v := range hotels.Result {
	//		if filterList.RoomOptions.Bathroom {
	//			if !BathroomIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.RoomOptions.Tv {
	//			if !TvIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.RoomOptions.Safe {
	//			if !SafeIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.RoomOptions.AirConditioning {
	//			if !AirConditioningIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.RoomOptions.WiFi {
	//			if !WiFiInRoomIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.RoomOptions.Hairdryer {
	//			if !HairdryerIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		// hotel filter
	//		if filterList.HotelOptions.Parking {
	//			if !CarParkingIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Restaurant {
	//			if !RestaurantIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.SwimmingPool {
	//			if !SwimmingPoolIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Spa {
	//			if !SpaIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.PrivateBeach {
	//			if !PrivateBeachIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Fitness {
	//			if !GymFitnessCentersIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.PetsAllowed {
	//			if !PetsAllowedIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.Laundry {
	//			if !LaundryIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.HotelOptions.AllDayReception {
	//			if !ReceptionAllDayIncluded(v) {
	//				continue
	//			}
	//		}
	//
	//		// food and payment filter
	//
	//		// price and distance filter
	//		if filterList.Price.Included {
	//			if !PriceBetween(v, filterList.Price.Values[0], filterList.Price.Values[1]) {
	//				continue
	//			}
	//		}
	//
	//		if filterList.Distance.Included {
	//			if !DistanceBetween(v, filterList.Distance.Values[0], filterList.Distance.Values[1]) {
	//				continue
	//			}
	//		}
	//	}
	//}
	//
	//for _, v := range hotels.Result {
	//	if filterList.RoomOptions.Bathroom {
	//		if !BathroomIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.RoomOptions.Tv {
	//		if !TvIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.RoomOptions.Safe {
	//		if !SafeIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.RoomOptions.AirConditioning {
	//		if !AirConditioningIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.RoomOptions.WiFi {
	//		if !WiFiInRoomIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.RoomOptions.Hairdryer {
	//		if !HairdryerIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	// hotel filter
	//	if filterList.HotelOptions.Parking {
	//		if !CarParkingIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.Restaurant {
	//		if !RestaurantIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.SwimmingPool {
	//		if !SwimmingPoolIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.Spa {
	//		if !SpaIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.PrivateBeach {
	//		if !PrivateBeachIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.Fitness {
	//		if !GymFitnessCentersIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.PetsAllowed {
	//		if !PetsAllowedIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.Laundry {
	//		if !LaundryIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.HotelOptions.AllDayReception {
	//		if !ReceptionAllDayIncluded(v) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	// food and payment filter
	//
	//	// price and distance filter
	//	if filterList.Price.Included {
	//		if !PriceBetween(v, filterList.Price.Values[0], filterList.Price.Values[1]) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//
	//	if filterList.Distance.Included {
	//		if !DistanceBetween(v, filterList.Distance.Values[0], filterList.Distance.Values[1]) {
	//			filteredHotels = RemoveByHotelId(filteredHotels, v.Id)
	//		}
	//	}
	//}
	//
	//for k, v := range filteredHotels {
	//	if filterList.Stars.Included {
	//if utils.IntArrayInclude(filterList.Stars.Values, 1) {
	//	if !Stars1(v) {
	//		filteredHotels[k] = filteredHotels[len(filteredHotels)-1]
	//		filteredHotels[len(filteredHotels)-1] = models.Hotel{}
	//		filteredHotels = filteredHotels[:len(filteredHotels)-1]
	//	}
	//}
	//
	//if utils.IntArrayInclude(filterList.Stars.Values, 2) {
	//	if !Stars2(v) {
	//		filteredHotels[k] = filteredHotels[len(filteredHotels)-1]
	//		filteredHotels[len(filteredHotels)-1] = models.Hotel{}
	//		filteredHotels = filteredHotels[:len(filteredHotels)-1]
	//	}
	//}
	//
	//if utils.IntArrayInclude(filterList.Stars.Values, 3) {
	//	if !Stars3(v) {
	//		filteredHotels[k] = filteredHotels[len(filteredHotels)-1]
	//		filteredHotels[len(filteredHotels)-1] = models.Hotel{}
	//		filteredHotels = filteredHotels[:len(filteredHotels)-1]
	//	}
	//}
	//
	//if utils.IntArrayInclude(filterList.Stars.Values, 4) {
	//	if !Stars4(v) {
	//		filteredHotels[k] = filteredHotels[len(filteredHotels)-1]
	//		filteredHotels[len(filteredHotels)-1] = models.Hotel{}
	//		filteredHotels = filteredHotels[:len(filteredHotels)-1]
	//	}
	//}
	//
	//if utils.IntArrayInclude(filterList.Stars.Values, 5) {
	//	if !Stars5(v) {
	//		// remove element
	//		filteredHotels[k] = filteredHotels[len(filteredHotels)-1]
	//		filteredHotels[len(filteredHotels)-1] = models.Hotel{}
	//		filteredHotels = filteredHotels[:len(filteredHotels)-1]
	//	}
	//}
	//}
	//}

	return filteredHotels
}
