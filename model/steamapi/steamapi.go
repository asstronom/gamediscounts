package steamapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/tidwall/gjson"
)

var (
	tagsPolicy *bluemonday.Policy = bluemonday.StrictPolicy()
)

type AppType int64

var (
	Game        AppType = 0
	Dlc         AppType = 1
	Music       AppType = 2
	Demo        AppType = 3
	Advertising AppType = 4
	Mod         AppType = 5
	Video       AppType = 6
	Unknown     AppType = 7
)

func (appType AppType) String() string {
	switch appType {
	case Game:
		return "game"
	case Dlc:
		return "dlc"
	case Music:
		return "music"
	case Demo:
		return "demo"
	case Advertising:
		return "advertising"
	case Mod:
		return "mod"
	case Video:
		return "video"
	case Unknown:
		return "unknown"
	}
	return "unknown"
}

func StringToAppType(s string) AppType {
	switch s {
	case "game":
		return Game
	case "dlc":
		return Dlc
	case "music":
		return Music
	case "demo":
		return Demo
	case "advertising":
		return Advertising
	case "mod":
		return Mod
	case "video":
		return Video
	case "unknown":
		return Unknown
	}
	return Unknown
}

type PriceOverview struct {
	Initial          float64
	Final            float64
	Discount_percent float64
	Currency         string
	IsFree           bool
}

type Package struct {
	PackageId  int
	Name       string
	AppIds     []int
	Price      PriceOverview
	Individual float64
	PageImage  string
	SmallLogo  string
}

type AppInfo struct {
	Appid       int
	Name        string
	Type        AppType
	DLC         []int
	Packages    []int
	Description string
	HeaderImage string
	Price       PriceOverview
	Genres      []string
}

func (a AppInfo) Equal(b AppInfo) bool {
	if a.Appid == b.Appid &&
		a.Name == b.Name &&
		a.Type == b.Type &&
		a.Description == b.Description &&
		a.HeaderImage == b.HeaderImage &&
		a.Price == b.Price {
		for i, v := range a.DLC {
			if v != b.DLC[i] {
				return false
			}
		}
		for i, v := range a.Packages {
			if v != b.Packages[i] {
				return false
			}
		}
		for i, v := range a.Genres {
			if v != b.Genres[i] {
				return false
			}
		}
	}
	return true
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

func extractPriceOverview(body *[]byte, appid int) (*PriceOverview, error) {
	path := fmt.Sprintf("%d.success", appid)
	if !gjson.Get(string(*body), path).Bool() {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid appid: %d", appid))
	}

	path = fmt.Sprintf("%d.data", appid)
	if len(gjson.Get(string(*body), path).Array()) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("Game may be free or has different pay methods. Appid: %d", appid))
	}

	pricepath := fmt.Sprintf("%d.data.price_overview", appid)
	res := gjson.Get(string(*body), pricepath)

	return &PriceOverview{res.Get("initial").Value().(float64),
		res.Get("final").Value().(float64),
		res.Get("discount_percent").Value().(float64),
		res.Get("currency").Value().(string),
		false}, nil

}

func GetAppPrice(appid int, cc string) (*PriceOverview, error) {
	steamapiLink := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%d&cc=%s&filters=price_overview", appid, cc)
	resp, err := http.Get(steamapiLink)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	res, err := extractPriceOverview(&body, appid)

	if err != nil {
		if err.Error() != fmt.Sprintf("Invalid appid: %d", appid) &&
			err.Error() != fmt.Sprintf("Game may be free or has different pay methods. Appid: %d", appid) {
			return nil, err
		}
	}

	return res, nil
}

func GetAppsPrice(appids *[]int, cc string) (*[]*PriceOverview, error) {
	if len(*appids) == 0 {
		return nil, fmt.Errorf("appids slice empty")
	}
	var steamapps string
	steamapps += strconv.Itoa((*appids)[0])
	for i := 1; i < len(*appids); i++ {
		steamapps += "," + strconv.Itoa((*appids)[i])
	}
	steamapiLink := fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%s&cc=%s&filters=price_overview", steamapps, cc)
	//fmt.Println(steamapiLink)
	//fmt.Println(appids)
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
		//fmt.Println(i, temp)
		if err != nil {
			if err.Error() != fmt.Sprintf("Invalid appid: %d", (*appids)[i]) &&
				err.Error() != fmt.Sprintf("Game may be free or has different pay methods. Appid: %d", (*appids)[i]) {
				return nil, err
			}
		}
		result = append(result, temp)
	}

	return &result, nil
}

func GetFeaturedCategories(cc string) ([]int, []PriceOverview, error) {
	steamapiLink := fmt.Sprintf("https://store.steampowered.com/api/featuredcategories?cc=%s", cc)
	resp, err := http.Get(steamapiLink)
	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	items := gjson.Get(string(body), "specials.items").Array()

	resultIDs := []int{}
	resultOverviews := []PriceOverview{}

	for i := 0; i < len(items); i++ {
		resultIDs = append(resultIDs, int(items[i].Get("id").Value().(float64)))
		resultOverviews = append(resultOverviews,
			PriceOverview{items[i].Get("original_price").Value().(float64),
				items[i].Get("final_price").Value().(float64),
				items[i].Get("discount_percent").Value().(float64),
				items[i].Get("currency").Value().(string),
				false})
	}
	return resultIDs, resultOverviews, nil
}

func GetAppInfo(appid int, cc string) (AppInfo, error) {
	var result AppInfo
	steamQuery := fmt.Sprintf(`http://store.steampowered.com/api/appdetails?appids=%d&cc=%s`, appid, cc)
	resp, err := http.Get(steamQuery)
	if err != nil {
		return AppInfo{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AppInfo{}, err
	}
	path := fmt.Sprintf("%d.success", appid)
	if !gjson.Get(string(body), path).Bool() {
		return AppInfo{}, fmt.Errorf("invalid appid")
	}
	path = fmt.Sprintf("%d.data", appid)
	infoJson := gjson.Get(string(body), path)
	if len(infoJson.Array()) == 0 {
		return AppInfo{}, fmt.Errorf("game may be free or has different pay methods")
	}
	result.Type = StringToAppType(infoJson.Get("type").String())
	result.Name = infoJson.Get("name").String()
	result.Appid, err = strconv.Atoi(infoJson.Get("steam_appid").String())
	if err != nil {
		return AppInfo{}, err
	}
	if infoJson.Get("is_free").Bool() {
		result.Price.IsFree = true
	} else {
		result.Price = PriceOverview{
			infoJson.Get("price_overview.initial").Float(),
			infoJson.Get("price_overview.final").Float(),
			infoJson.Get("price_overview.discount_percent").Float(),
			infoJson.Get("price_overview.currency").String(),
			false,
		}
	}
	for _, r := range infoJson.Get("dlc").Array() {
		result.DLC = append(result.DLC, int(r.Int()))
	}
	for _, r := range infoJson.Get("packages").Array() {
		result.Packages = append(result.Packages, int(r.Int()))
	}
	for _, r := range infoJson.Get("genres").Array() {
		result.Genres = append(result.Genres, r.Get("description").String())
	}
	result.Description = tagsPolicy.Sanitize(infoJson.Get("detailed_description").String())
	result.Description = strings.ReplaceAll(result.Description, "\t", "")
	result.HeaderImage = infoJson.Get("header_image").String()
	return result, nil
}

func GetPackageInfo(packageid int, cc string) (Package, error) {
	var result Package
	steamQuery := fmt.Sprintf(`http://store.steampowered.com/api/packagedetails?packageids=%d&cc=%s`, packageid, cc)
	resp, err := http.Get(steamQuery)
	if err != nil {
		return Package{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Package{}, err
	}

	path := fmt.Sprintf("%d.success", packageid)
	if !gjson.Get(string(body), path).Bool() {
		return Package{}, fmt.Errorf(fmt.Sprintf("Invalid appid: %d", packageid))
	}
	path = fmt.Sprintf("%d.data", packageid)
	infoJson := gjson.Get(string(body), path)
	if len(infoJson.Array()) == 0 {
		return Package{}, fmt.Errorf(fmt.Sprintf("Game may be free or has different pay methods. Appid: %d", packageid))
	}
	result.PackageId = packageid
	result.Name = infoJson.Get("name").String()
	for _, r := range infoJson.Get("apps").Array() {
		result.AppIds = append(result.AppIds, int(r.Get("id").Int()))
	}
	result.Price = PriceOverview{
		infoJson.Get("price.initial").Float(),
		infoJson.Get("price.final").Float(),
		infoJson.Get("price.discount_percent").Float(),
		infoJson.Get("price.currency").String(),
		infoJson.Get("price.initial").Float() == 0,
	}
	result.Individual = infoJson.Get("price.individual").Float()
	result.PageImage = infoJson.Get("page_image").String()
	result.SmallLogo = infoJson.Get("small_logo").String()
	return result, nil
}
