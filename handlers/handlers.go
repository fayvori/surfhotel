package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/zheeeng/pagination"
	"go.mongodb.org/mongo-driver/bson"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"surf-hotel-engine/constants"
	"surf-hotel-engine/databases"
	"surf-hotel-engine/filter"
	"surf-hotel-engine/models"
	"surf-hotel-engine/utils"
	"sort"
)

var client = resty.New()

// resty client
func getRestyClient() *resty.Client {
	client := resty.New()

	//client.SetTimeout(5 * time.Second)
	//retryAfter := func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
	//	return 2, errors.New("quota exceeded")
	//}

	client.
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				fmt.Println(r.StatusCode())
				return r.StatusCode() == http.StatusConflict
			},
		).
		//SetRetryAfter(retryAfter).
		SetRetryCount(30)

	return client
}

func GetSearchId(ctx *gin.Context) {
	// url params
	q := ctx.Request.URL.Query()

	// get md5 hash
	//1.  Token
	//2.  Marker
	//3.  adultsCount=2
	//4.  checkIn=2021-12-10
	//5.  checkOut=2021-12-13
	//6.  childAge1=10
	//7.  childrenCount=1
	//8.  currency=USD
	//9.  customerIP=109.252.191.186
	//10. cityId
	//11. lang=ru
	//12. waitForResult=0

	cityId, _ := strconv.Atoi(q.Get("cityId"))

	hash := Md5HotelSearchHash(fmt.Sprintf("%s:%s:%s:%s:%s:%d:%s:%s:%s:%s",
		constants.TOKEN,
		constants.MARKER,
		q.Get("adultsCount"),
		q.Get("checkIn"),
		q.Get("checkOut"),
		cityId,
		q.Get("currency"),
		constants.CUSTOMER_IP,
		q.Get("lang"),
		q.Get("waitForResult"),
	))

	// initialize resty client
	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("%s/start.json?cityId=%d&checkIn=%s&checkOut=%s&adultsCount=%s&customerIP=%s&lang=%s&currency=%s&waitForResult=%s&marker=%s&signature=%s",
			constants.HOTELLOOK_ADDR,
			cityId,
			q.Get("checkIn"),
			q.Get("checkOut"),
			q.Get("adultsCount"),
			constants.CUSTOMER_IP,
			q.Get("lang"),
			q.Get("currency"),
			q.Get("waitForResult"),
			constants.MARKER,
			hash,
		))

	if err != nil {
		fmt.Println(err.Error())
	}

	var response models.HotelResponse

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		fmt.Println(err.Error())
	}

	filteredHotels := filter.FilterHotelResponse(&response, ctx.Request.Body, q.Get("search"))

	sort.SliceStable(filteredHotels, func (i, j int) bool {
		return filteredHotels[i].Price < filteredHotels[j].Price
	})

	sort.SliceStable(filteredHotels, func (i, j int) bool {
		return filteredHotels[i].Stars > filteredHotels[j].Stars
	})

	// http://localhost:4000/hotels/api/v1/search_id?checkIn=2022-07-20&checkOut=2022-07-30&adultsCount=1&currency=USD&lang=ru&cityId=12153&customerIP=185.15.63.18&waitForResult=1
	pgPages := fmt.Sprintf("http://localhost:4000/api/v1/search_id?sortBy=price&sortAsc=0&roomsCount=2&offset=0&marker=327233&page=%s&page_size=200", q.Get("page"))

	pg := pagination.DefaultPagination()
	pgt := pg.Parse(pgPages)
	paginatedData := pgt.WrapWithTruncate(TrunctableHotels(filteredHotels), len(response.Result))

	responseBody, _ := json.Marshal(paginatedData)
	var paginateResponse models.PaginationHotels
	err = json.Unmarshal(responseBody, &paginateResponse)

	utils.AddPhotosToHotelResponse(&paginateResponse)

	// utils for app filters
	paginateResponse.TotalFiltered = len(filteredHotels)

	// get min and max values for filters
	if len(response.Result) > 0 {
		min, max := utils.FindMinAndMax(response.Result)
		paginateResponse.MinPrice = min
		paginateResponse.MaxPrice = max
	}

	ctx.JSON(http.StatusOK, paginateResponse)
}

func SingleHotelSearchId(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	hash := Md5HotelSearchHash(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s:%s:%s",
		constants.TOKEN,
		constants.MARKER,
		q.Get("adultsCount"),
		q.Get("checkIn"),
		q.Get("checkOut"),
		q.Get("currency"),
		constants.CUSTOMER_IP,
		q.Get("hotelId"),
		q.Get("lang"),
		q.Get("waitForResult"),
	))

	resp, err := getRestyClient().R().
		EnableTrace().
		Post(fmt.Sprintf("%s/start.json?hotelId=%s&checkIn=%s&checkOut=%s&adultsCount=%s&customerIP=%s&lang=%s&currency=%s&waitForResult=%s&marker=%s&signature=%s",
			constants.HOTELLOOK_ADDR,
			q.Get("hotelId"),
			q.Get("checkIn"),
			q.Get("checkOut"),
			q.Get("adultsCount"),
			constants.CUSTOMER_IP,
			q.Get("lang"),
			q.Get("currency"),
			q.Get("waitForResult"),
			constants.MARKER,
			hash,
		))

	if err != nil {
		fmt.Println(err.Error())
	}

	var response models.HotelResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		fmt.Println(err.Error())
	}

	utils.AddPhotosToSingleHotel(&response)

	ctx.JSON(http.StatusOK, response)
}

type TrunctableHotels []models.Hotel

func (tb TrunctableHotels) Slice(startIndex, endIndex int) pagination.Truncatable {
	return tb[startIndex:endIndex]
}

func (tb TrunctableHotels) Len() int {
	return len(tb)
}

func SearchHotel(ctx *gin.Context) {
	// query params
	q := ctx.Request.URL.Query()

	limit := "2000"

	// generate md5 hash
	hash := Md5HotelSearchHash(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
		constants.TOKEN,
		constants.MARKER,
		limit,
		"0",
		q.Get("roomsCount"),
		q.Get("searchId"),
		"0",
		"popularity",
	))

	// response model
	var response models.HotelResponse

	if q.Get("searchId") == "0" {
		return
	}

	// resty client request to api
	resp, err := getRestyClient().R().
		SetQueryParams(map[string]string{
			"searchId":   q.Get("searchId"),
			"limit":      limit,
			"sortBy":     "popularity",
			"sortAsc":    "0",
			"roomsCount": q.Get("roomsCount"),
			"offset":     "0",
			"marker":     constants.MARKER,
			"signature":  hash,
		}).
		Get(fmt.Sprintf("%s/getResult.json", constants.HOTELLOOK_ADDR))

	// check error
	if err != nil {
		fmt.Println(err)
	}

	// unmarshal bytes of body to response model
	err = json.Unmarshal(resp.Body(), &response)

	// check unmarshal error
	if err != nil {
		fmt.Println(err.Error())
	}

	// filter response
	filteredHotels := filter.FilterHotelResponse(&response, ctx.Request.Body, q.Get("search"))

	sort.Slice(filteredHotels, func (i, j int) bool {
		return filteredHotels[i].Rating > filteredHotels[j].Rating
	})

	// pagination config
	pg := pagination.DefaultPagination()
	pgt := pg.Parse(fmt.Sprintf("http://localhost:8080/api/v1/hotel_search?sortBy=price&sortAsc=0&roomsCount=2&offset=0&marker=327233&searchId=4558050&limit=2000&page=1&page_size=1000"))
	paginatedData := pgt.WrapWithTruncate(TrunctableHotels(filteredHotels), len(response.Result))

	responseBody, _ := json.Marshal(paginatedData)
	var paginateResponse models.PaginationHotels
	err = json.Unmarshal(responseBody, &paginateResponse)

	// check pagination error
	if err != nil {
		fmt.Println(err.Error())
	}

	//adding photos
	utils.AddPhotosToHotelResponse(&paginateResponse)

	// utils for app filters
	paginateResponse.TotalFiltered = len(filteredHotels)

	// get min and max values for filters
	if len(response.Result) > 0 {
		min, max := utils.FindMinAndMax(response.Result)
		paginateResponse.MinPrice = min
		paginateResponse.MaxPrice = max
	}

	// return response as json
	if len(response.Result) > 0 {
		//ctx.JSON(http.StatusOK, string(resp.Body()))
		ctx.JSON(http.StatusOK, paginateResponse)
	} else {
		fmt.Printf("resp body if its nil %v\n", string(resp.Body()))
		ctx.JSON(http.StatusConflict, nil)
	}

	fmt.Println("url", resp.Request.URL)
	fmt.Println("length of filteredHotels", len(filteredHotels))
}

func HotelsDb(ctx *gin.Context) {
	var currencies models.Currencies

	filePath, _ := filepath.Abs("currencies_rates.json")
	f, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Println(err.Error())
	}

	err = json.Unmarshal(f, &currencies)

	if err != nil {
		log.Println(err.Error())
	}

	q := ctx.Request.URL.Query()

	if len(q.Get("cityId")) == 0 {
		return
	}

	options := options2.Find()
	options.
		SetLimit(150)

	currency := q.Get("currency")
	nights, _ := strconv.Atoi(q.Get("nights"))

	cur, err := databases.MongoClient().Find(databases.Ctx, bson.M{"cityid": q.Get("cityId"), "lang": q.Get("lang")}, options)
	if err != nil {
		log.Printf(err.Error())
	}

	var hotels models.HotelResponse
	if err = cur.All(databases.Ctx, &hotels.Result); err != nil {
		log.Println(err.Error())
	}

	// change currency rate
	for _, v := range hotels.Result {
		priceWithNightsInUSD := float64(v.Price * nights)
		v.Rooms[0].Total = int(priceWithNightsInUSD * currencies.Rates[currency])
		v.Rooms[0].Price = int(float64(v.Rooms[0].Price) * currencies.Rates[currency])

		if v.PhotoHotel == nil || len(v.PhotoHotel) == 0 {
			// set default photo
			v.PhotoHotel = []int{5162091993}
		}
	}

	ctx.JSON(http.StatusOK, hotels)
}

func UpdateCurrencies(ctx *gin.Context) {
	resp, err := getRestyClient().R().
		EnableTrace().
		SetHeader("apikey", constants.CURRENCY_API_KEY).
		SetQueryParam("base", constants.CURRENCY_BASE).
		Get(constants.CURRENCY_API_URL)

	if err != nil {
		log.Println(err.Error())
	}

	if err != nil {
		log.Println(err.Error())
	}

	f, err := os.Create("currencies_rates.json")

	if err != nil {
		log.Println(err.Error())
	}

	defer f.Close()

	f.Write(resp.Body())

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
