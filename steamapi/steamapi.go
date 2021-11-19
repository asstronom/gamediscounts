package steamapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

type PriceOverview struct {
	Initial          float64
	Final            float64
	Discount_percent float64
	Currency         string
	Formatted        string
}

func GetAppList() []gjson.Result {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	res := gjson.Get(string(body), "applist.apps").Array()
	return res
}

func extractPriceOverview(body *[]byte, appid int) (PriceOverview, error) {
	path := fmt.Sprintf("%d.success", appid)
	if gjson.Get(string(*body), path).Value() == "false" {
		return PriceOverview{}, fmt.Errorf(fmt.Sprintf("Invalid appid: %d", appid))
	}
	path = fmt.Sprintf("%d.data", appid)
	if len(gjson.Get(string(*body), path).Array()) == 0 {
		return PriceOverview{}, fmt.Errorf(fmt.Sprintf("Invalid appid: %d", appid))
	}

	pricepath := fmt.Sprintf("%d.data.price_overview", appid)
	res := gjson.Get(string(*body), pricepath)
	return PriceOverview{res.Get("initial").Value().(float64),
		res.Get("final").Value().(float64),
		res.Get("discount_percent").Value().(float64),
		res.Get("currency").Value().(string),
		res.Get("final_formatted").Value().(string)}, nil

}

func GetAppPrice(appid int, cc string) (PriceOverview, error) {
	steamapiLink := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%d&cc=%s&filters=price_overview", appid, cc)
	resp, err := http.Get(steamapiLink)
	if err != nil {
		return PriceOverview{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return PriceOverview{}, err
	}
	path := fmt.Sprintf("%d.success", appid)
	if !gjson.Get(string(body), path).Bool() {
		return PriceOverview{}, fmt.Errorf(fmt.Sprintf("Invalid appid: %d", appid))
	}
	path = fmt.Sprintf("%d.data", appid)
	if len(gjson.Get(string(body), path).Array()) == 0 {
		return PriceOverview{}, fmt.Errorf(fmt.Sprintf("Invalid appid: %d", appid))
	}

	fmt.Println("Hello")

	pricepath := fmt.Sprintf("%d.data.price_overview", appid)
	res := gjson.Get(string(body), pricepath)
	return PriceOverview{res.Get("initial").Value().(float64),
		res.Get("final").Value().(float64),
		res.Get("discount_percent").Value().(float64),
		res.Get("currency").Value().(string),
		res.Get("final_formatted").Value().(string)}, nil
}

func GetAppsPrice(appids *[]int, cc string) (*[]*PriceOverview, error) {
	var steamapps string
	steamapps += strconv.Itoa((*appids)[0])
	for i := 1; i < len(*appids); i++ {
		steamapps += "," + strconv.Itoa((*appids)[i])
	}
	steamapiLink := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%s&cc=%s&filters=price_overview", steamapps, cc)
	fmt.Println(steamapiLink)
	fmt.Println(appids)
	resp, err := http.Get(steamapiLink)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := []*PriceOverview{}

	for i := 0; i < len(*appids); i++ {
		temp, err := extractPriceOverview(&body, (*appids)[i])
		fmt.Println(i, temp)
		if err != nil {
			if err.Error() != fmt.Sprintf("Invalid appid: %d", (*appids)[i]) {
				return nil, err
			}
		}
		result = append(result, &temp)
	}

	return &result, nil
}
