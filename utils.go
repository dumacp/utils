package utils

import (
	"encoding/json"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/yanatan16/itertools"
)

type ByNumericalFilename []os.FileInfo

func (nf ByNumericalFilename) Len() int      { return len(nf) }
func (nf ByNumericalFilename) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilename) Less(i, j int) bool {

	pathA := nf[i].Name()
	pathB := nf[j].Name()

	subA := pathA[0:strings.LastIndex(pathA, ".")]
	subB := pathB[0:strings.LastIndex(pathB, ".")]

	if subA == subB {
		pathA_ := nf[i].Name()
		pathB_ := nf[j].Name()
		a, err1 := strconv.ParseInt(pathA_[strings.LastIndex(pathA_, ".")+1:len(pathA_)], 10, 64)
		b, err2 := strconv.ParseInt(pathB_[strings.LastIndex(pathB_, ".")+1:len(pathB_)], 10, 64)

		if err1 == nil && err2 == nil {
			return a < b
		}
	}
	return pathA < pathB
}

func SortFileInfo(files []os.FileInfo) {
	sort.Sort(ByNumericalFilename(files))
}

func CloseChannels(value int, quits ...chan int) {
	for _, quit := range quits {
		CloseChannel(value, quit)
	}
}

func CloseChannel(value int, quit chan int) {
	select {
	case quit <- value:
		//fmt.Printf("channel value: %v\n", value)
	default:
		close(quit)
	}
}

func FinishChannel(chs ...chan interface{}) {
	for _, ch := range chs {
		for _ = range ch {
			//fmt.Printf("valueChannel %v: %v\n",i, v)
		}
	}
}

func NewTee(it itertools.Iter, n, lenBuffer int) []itertools.Iter {

	if lenBuffer > 30 {
		panic("maxim lenBuffer allow is 30")
	}

	iters := make([]itertools.Iter, n)
	for i := 0; i < n; i++ {
		iters[i] = make(itertools.Iter, lenBuffer)
	}
	go func() {
		for newval := range it {
			for _, v := range iters {
				v <- newval
			}
		}
	}()
	return iters
}

//:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
//:::                                                                         :::
//:::  This routine calculates the distance between two points (given the     :::
//:::  latitude/longitude of those points). It is being used to calculate     :::
//:::  the distance between two locations using GeoDataSource (TM) prodducts  :::
//:::                                                                         :::
//:::  Definitions:                                                           :::
//:::    South latitudes are negative, east longitudes are positive           :::
//:::                                                                         :::
//:::  Passed to function:                                                    :::
//:::    lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)  :::
//:::    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)  :::
//:::    unit = the unit you desire for results                               :::
//:::           where: 'M' is statute miles (default)                         :::
//:::                  'K' is kilometers                                      :::
//:::                  'N' is nautical miles                                  :::
//:::                                                                         :::
//:::  Worldwide cities and other features databases with latitude longitude  :::
//:::  are available at https://www.geodatasource.com                         :::
//:::                                                                         :::
//:::  For enquiries, please contact sales@geodatasource.com                  :::
//:::                                                                         :::
//:::  Official Web site: https://www.geodatasource.com                       :::
//:::                                                                         :::
//:::               GeoDataSource.com (C) All Rights Reserved 2018            :::
//:::                                                                         :::
//:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
func Distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func PrettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "	")
	if err != nil {
		return ""
	}
	return string(b)
}
