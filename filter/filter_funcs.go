package filter

import (
	"strings"
	"surf-hotel-engine/models"
)

// helpers

func Contains(slice []int, item int) bool {
	set := make(map[int]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// stars
func Stars0(ht models.Hotel) bool {
	return ht.Stars == 0
}

func Stars1(ht models.Hotel) bool {
	return ht.Stars == 1
}

func Stars2(ht models.Hotel) bool {
	return ht.Stars == 2
}

func Stars3(ht models.Hotel) bool {
	return ht.Stars == 3
}

func Stars4(ht models.Hotel) bool {
	return ht.Stars == 4
}

func Stars5(ht models.Hotel) bool {
	return ht.Stars == 5
}

// hotel types
func HotelType(ht models.Hotel) bool {
	//return strings.Contains(strings.ToLower(ht.Name), "отель") || strings.Contains(strings.ToLower(ht.Name), "гостиница") || strings.Contains(strings.ToLower(ht.Name), "холидей")
	return strings.Contains(strings.ToLower(ht.Name), "отель")
}

func HostelType(ht models.Hotel) bool {
	return strings.Contains(strings.ToLower(ht.Name), "хостел")
}

func ApartmentsType(ht models.Hotel) bool {
	return strings.Contains(strings.ToLower(ht.Name), "апарт")
}

func VillasType(ht models.Hotel) bool {
	return strings.Contains(strings.ToLower(ht.Name), "вилл")
}

func GuestHousesType(ht models.Hotel) bool {
	return strings.Contains(strings.ToLower(ht.Name), "гост")
}

// between filters

func PriceBetween(ht models.Hotel, fv int, sv int) bool {
	return ht.MaxPrice >= fv && ht.MaxPrice <= sv
}

func DistanceBetween(ht models.Hotel, fv float64, sv float64) bool {
	return ht.Distance >= fv && ht.Distance <= sv
}

//func RatingBetween(ht models.Hotel, fv int, sv int) bool {
//	return ht.Rating >= fv && ht.Rating <= sv
//}

// hotel options

func BreakfastIncluded(ht models.Hotel) bool {
	return ht.Rooms[0].Options.Breakfast
}

func AllInclusiveIncluded(ht models.Hotel) bool {
	return ht.Rooms[0].Options.AllInclusive
}

func RefundableIncluded(ht models.Hotel) bool {
	return ht.Rooms[0].Options.Refundable
}

// ???
func PaymentNowIncluded(ht models.Hotel) bool {
	return true
}

func PaymentInHotelIncluded(ht models.Hotel) bool {
	return true
}

//func PaymentNow(ht models.Hotel) bool {
//	return ht.Rooms[0].Options.Refundable
//}
//
//func PaymentInHotel(ht models.Hotel) bool {
//	return ht.Rooms[0].Options.
//}

// amenities
// 22  bathroom
// 133 wi-fi in room
// 11  air conditioning
// 2   hairdryer
// 4   tv
// 3   safe
// 56  car parking
// 9   restaurant
// 65  spa
// 102 private beach
// 40  gym / fitness center
// 28  pets allowed
// 13  laundry service
// 50  24h. reception
// 38  swimming pool

func BathroomIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 22)
}

func WiFiInRoomIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 133)
}

func AirConditioningIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 11)
}

func HairdryerIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 2)
}

func TvIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 4)
}

func SafeIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 3)
}

func CarParkingIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 56)
}

func RestaurantIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 9)
}

func SpaIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 65)
}

func PrivateBeachIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 102)
}

func GymFitnessCentersIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 40)
}

func PetsAllowedIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 28)
}

func LaundryIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 13)
}

func ReceptionAllDayIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 50)
}

func SwimmingPoolIncluded(ht models.Hotel) bool {
	return Contains(ht.Amenities, 38)
}
