package steamapi

import "testing"

func TestGetAppPrice(t *testing.T) {
	//Arrange
	testTable := []struct {
		appid    int
		cc       string
		expected *PriceOverview
	}{
		{
			620,
			"ua",
			&PriceOverview{16900, 16900, 0, "UAH", false},
		},
		{
			570,
			"ua",
			nil,
		},
		{
			271590,
			"ua",
			nil,
		},
		{
			216938,
			"ua",
			nil,
		},
	}

	//Act

	//Assert
	for _, v := range testTable {
		price, err := GetAppPrice(v.appid, v.cc)
		if err != nil {
			t.Errorf("Incorrect")
		}
		if v.expected == nil && price == nil {
			continue
		}
		if *price != *v.expected {
			t.Errorf("Error")
		}
	}
}

func TestGetAppsPrice(t *testing.T) {
	testTable := []struct {
		inputApps      []int
		expectedPrices []*PriceOverview
	}{
		{
			[]int{620, 570, 271590, 216938, 1798520},
			[]*PriceOverview{
				{16900, 16900, 0, "UAH", false},
				nil,
				nil,
				nil,
				{4000, 4000, 0, "UAH", false},
			},
		},
	}

	zeroLengthTest := []int{}

	for i, v := range testTable {
		prices, err := GetAppsPrice(&v.inputApps, "ua")
		if err != nil {
			t.Errorf("Error while getting apps price. Test number %d", i)
		}
		for j, k := range *prices {
			if k == nil && v.expectedPrices[j] == nil {
				continue
			}
			if *k != *v.expectedPrices[j] {
				t.Errorf("Incorrect price, %d", j)
			}
		}
	}

	_, err := GetAppsPrice(&zeroLengthTest, "ua")
	if err == nil {
		t.Errorf("No error in zero length test")
	}
}

func TestAppType(t *testing.T) {
	//Arrange
	testTable := []struct {
		appType AppType
		strType string
	}{
		{
			Game,
			"game",
		},
		{
			Dlc,
			"dlc",
		},
		{
			Music,
			"music",
		},
		{
			Demo,
			"demo",
		},
		{
			Advertising,
			"advertising",
		},
		{
			Mod,
			"mod",
		},
		{
			Video,
			"video",
		},
		{
			Unknown,
			"unknown",
		},
	}

	unknownTest := struct {
		appType AppType
		strType string
	}{
		Unknown,
		"sdfawdsfh",
	}
	//ASSERT

	for _, v := range testTable {
		if v.appType.String() != v.strType {
			t.Errorf("Error %s != %s", v.appType.String(), v.strType)
		}
		if StringToAppType(v.strType) != v.appType {
			t.Errorf("Error while converting to appType %d != %d", StringToAppType(v.strType), v.appType)
		}
	}
	if StringToAppType(unknownTest.strType) != unknownTest.appType {
		t.Errorf("Error in unknown test, %d != %d", StringToAppType(unknownTest.strType), unknownTest.appType)
	}
}

func TestGetAppInfo(t *testing.T) {
	testTable := []struct {
		appid       int
		cc          string
		expected    AppInfo
		expectedErr error
	}{
		{
			620,
			"ua",
			AppInfo{
				620,
				"Portal 2",
				Game,
				[]int{323180},
				[]int{7877, 204528, 8187},
				"Portal 2 draws from the award-winning formula of innovative gameplay, story, and music that earned the original Portal over 70 industry accolades and created a cult following.The single-player portion of Portal 2 introduces a cast of dynamic new characters, a host of fresh puzzle elements, and a much larger set of devious test chambers. Players will explore never-before-seen areas of the Aperture Science Labs and be reunited with GLaDOS, the occasionally murderous computer companion who guided them through the original game.The game’s two-player cooperative mode features its own entirely separate campaign with a unique story, test chambers, and two new player characters. This new mode forces players to reconsider everything they thought they knew about portals. Success will require them to not just act cooperatively, but to think cooperatively.Product FeaturesExtensive single player: Featuring next generation gameplay and a wildly-engrossing story.Complete two-person co-op: Multiplayer game featuring its own dedicated story, characters, and gameplay.Advanced physics: Allows for the creation of a whole new range of interesting challenges, producing a much larger but not harder game.Original music.Massive sequel: The original Portal was named 2007&#39;s Game of the Year by over 30 publications worldwide. Editing Tools: Portal 2 editing tools will be included.",
				"https://cdn.akamai.steamstatic.com/steam/apps/620/header.jpg?t=1610490805",
				PriceOverview{16900, 16900, 0, "UAH", false},
				[]string{"Action", "Adventure"},
			},
			nil,
		},
	}

	//ASSERT
	for i, v := range testTable {
		appinfo, err := GetAppInfo(v.appid, v.cc)
		if err != v.expectedErr {
			t.Errorf("Errors dont match. %d", i)
		}
		if !appinfo.Equal(v.expected) {
			t.Errorf("Wrong appinfo %d", i)
		}
	}
}
