/*
Purpose:
- simple route statistics

Description:
- Calculates route statistics based on CSV routing data.

Releases:
- v1.0.0 - 2023/10/10 : initial release

Author:
- Klaus Tockloth

Copyright:
- © 2023 | Klaus Tockloth

Contact:
- freizeitkarte@googlemail.com

Remarks:
- Lint: golangci-lint run --no-config --enable gocritic
- Vulnerability detection: govulncheck ./...
- Test data (lengerich-steinbrüche): https://brouter.m11n.de/#map=14/52.1719/7.9318/standard,route-quality&lonlats=7.88421,52.179888;7.884904,52.180073;7.885889,52.181287;7.885197,52.180718;7.884926,52.18092;7.885337,52.181416;7.884946,52.182789;7.885997,52.183609;7.887607,52.183746;7.888595,52.183346;7.889509,52.184044;7.889864,52.185578;7.890424,52.185351;7.891487,52.186053;7.892523,52.18589;7.899235,52.182669;7.900188,52.182353;7.90117,52.18263;7.905138,52.181593;7.906022,52.180869;7.905771,52.180141;7.906038,52.179949;7.908716,52.180496;7.911386,52.179825;7.913069,52.179048;7.91216,52.178538;7.912681,52.178236;7.912516,52.177917;7.913187,52.178198;7.916564,52.176681;7.920394,52.175886;7.921036,52.175039;7.921824,52.17474;7.923159,52.174881;7.925341,52.173813;7.925177,52.173561;7.925785,52.172656;7.92984,52.172907;7.930692,52.171658;7.935655,52.170951;7.937017,52.170223;7.937456,52.16975;7.937014,52.169751;7.93847,52.168717;7.93857,52.169035;7.9381,52.169295;7.938503,52.170088;7.939615,52.170242;7.940585,52.170225;7.942239,52.169479;7.942616,52.168375;7.949555,52.167805;7.951484,52.168723;7.952913,52.168718;7.955668,52.168064;7.955489,52.167181;7.95588,52.167999;7.963219,52.166236;7.964282,52.165471;7.966418,52.165469;7.970746,52.163448;7.976291,52.16273;7.977878,52.162156;7.979201,52.16069;7.980222,52.160259;7.980075,52.159624;7.98069,52.158905;7.98017,52.15912;7.980198,52.156633;7.979482,52.156377;7.979488,52.155477;7.97661,52.155373;7.973866,52.156027;7.971959,52.155964;7.971679,52.155144;7.97045,52.155475;7.967401,52.155224;7.962469,52.156639;7.962313,52.156369;7.962983,52.155664;7.962373,52.156022;7.962585,52.155777;7.963006,52.155901;7.96255,52.156329;7.962659,52.157021;7.96153,52.158427;7.962099,52.159529;7.96005,52.160486;7.952759,52.159701;7.952659,52.158977;7.954137,52.15794;7.95393,52.157189;7.951126,52.157492;7.95108,52.158089;7.950518,52.158317;7.948177,52.158499;7.94912,52.159507;7.948428,52.160077;7.948934,52.160524;7.948651,52.160963;7.946573,52.16044;7.945924,52.159633;7.943457,52.160153;7.942472,52.159526;7.930824,52.160787;7.929119,52.16049;7.928997,52.162235;7.928324,52.162831;7.923977,52.163539;7.921555,52.164608;7.922918,52.166126;7.923372,52.168819;7.923135,52.169329;7.924065,52.170274;7.923827,52.171375;7.924375,52.172225;7.924063,52.171893;7.923512,52.172148;7.92234,52.174746;7.92107,52.175071;7.920478,52.175873;7.916557,52.176664;7.91308,52.178317;7.912595,52.178218;7.912343,52.178547;7.913089,52.179105;7.911624,52.179839;7.908929,52.180605;7.906896,52.180503;7.905973,52.180014;7.906133,52.180894;7.905692,52.181365;7.901211,52.182444;7.902225,52.182551;7.900926,52.183162;7.901275,52.183697;7.899233,52.184778;7.896837,52.185472;7.895259,52.186558;7.892808,52.186769;7.888485,52.188377;7.886216,52.188457;7.884505,52.187887;7.882725,52.18675;7.883557,52.185586;7.883988,52.183828;7.885025,52.183074;7.884803,52.182768;7.885113,52.1816;7.88479,52.180959;7.885164,52.180577;7.884225,52.180003&profile=dummy
- Inefficiency: uses multiple loops (instead of one loop) for simplicity
- Deepest point on land: Dead Sea approx. -430 m

Links:
- https://brouter.m11n.de
- https://bikerouter.de
*/

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// general program info
var (
	progName    = os.Args[0]
	progVersion = "v1.0.0"
	progDate    = "2023/10/10"
	progPurpose = "simple route statistics"
	progInfo    = "Calculates route statistics based on CSV routing data."
)

// CSV record logical field names
const (
	Longitude = iota
	Latitude
	Elevation
	Distance
	CostPerKm
	ElevCost
	TurnCost
	NodeCost
	InitialCost
	WayTags
	NodeTags
	Time
	Energy
)

// statistics
var highways = make(map[string]int)
var surfaces = make(map[string]int)
var smoothness = make(map[string]int)
var mtbScales = make(map[string]int)
var sacScales = make(map[string]int)
var highwaysUnpaved int
var highwaysPaved int
var metersDownhill int
var metersUphill int
var altitudeMetersMin int
var altitudeMetersMax int

// constants
const deepestPointOnLand = -430

/*
main starts this program.
*/
func main() {
	fmt.Printf("\nProgram:\n")
	fmt.Printf("  Name    : %s\n", progName)
	fmt.Printf("  Release : %s - %s\n", progVersion, progDate)
	fmt.Printf("  Purpose : %s\n", progPurpose)
	fmt.Printf("  Info    : %s\n\n", progInfo)

	// file argument required
	if len(os.Args) != 2 {
		printUsage()
	}

	filename := os.Args[1]
	fi, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("error <%v> at os.Stat(), file = %v", err, filename)
	}

	fmt.Printf("Input file:\n")
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	fmt.Printf("name = %s\n", filename)
	fmt.Printf("size = %d byte\n", fi.Size())
	fmt.Printf("time = %s\n", fi.ModTime().Format(time.RFC3339))
	fmt.Printf("%s\n", strings.Repeat("-", 47))

	records, err := csvReadAll(filename)
	if err != nil {
		log.Fatalf("error <%v> at csvReadAll(), file = %v", err, filename)
	}

	// build highways map
	// ------------------
	for index, record := range records {
		if index == 0 {
			continue // skip csv header
		}
		distance, _ := strconv.Atoi(record[Distance])
		waytags := strings.Split(record[WayTags], " ")
		highwayIdentifier := getHighwayIdentifier(waytags)
		a1 := highways[highwayIdentifier]
		highways[highwayIdentifier] = a1 + distance
	}
	fmt.Printf("\nHighway types:\n")
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	keys := make([]string, 0, len(highways))
	distanceTotal := 0
	for key := range highways {
		keys = append(keys, key)
		distanceTotal += highways[key]
	}
	sort.Strings(keys)
	for _, key := range keys {
		percent := float64(highways[key]) / float64(distanceTotal) * 100.0
		fmt.Printf("%-28s  %5.1f %%  %6d m\n", key, percent, highways[key])
	}
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	fmt.Printf("%-28s  %5.1f %%  %6d m\n", "total", 100.0, distanceTotal)

	// calculate altitude meters
	// -------------------------
	lastElevation := math.MaxInt
	for index, record := range records {
		if index == 0 {
			continue // skip csv header
		}
		elevation, _ := strconv.Atoi(record[Elevation])
		if elevation < deepestPointOnLand {
			// elevation value probably incorrect
			continue
		}
		if lastElevation == math.MaxInt {
			lastElevation = elevation
			altitudeMetersMin = elevation
			altitudeMetersMax = elevation
			continue
		}
		if elevation > lastElevation {
			metersUphill += elevation - lastElevation
		} else {
			metersDownhill += elevation - lastElevation
		}
		if elevation < altitudeMetersMin {
			altitudeMetersMin = elevation
		}
		if elevation > altitudeMetersMax {
			altitudeMetersMax = elevation
		}
		lastElevation = elevation
	}
	fmt.Printf("\nAltitude meters (approximately):\n")
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	fmt.Printf("%-38s %+6d m\n", "downhill", metersDownhill)
	fmt.Printf("%-38s %+6d m\n", "uphill", metersUphill)
	fmt.Printf("%-38s %6d m\n", "height min", altitudeMetersMin)
	fmt.Printf("%-38s %6d m\n", "height max", altitudeMetersMax)
	fmt.Printf("%-38s %6d m\n", "height difference", altitudeMetersMax-altitudeMetersMin)
	fmt.Printf("%s\n", strings.Repeat("-", 47))

	// build surfaces map
	// ------------------
	for index, record := range records {
		if index == 0 {
			continue // skip csv header
		}
		distance, _ := strconv.Atoi(record[Distance])
		waytags := strings.Split(record[WayTags], " ")
		surface := ""
		for _, waytag := range waytags {
			if strings.HasPrefix(waytag, "surface") {
				surface = waytag
			}
		}
		if surface == "" {
			// no surface tag found, use highway type specific default value
			highwayIdentifier := getHighwayIdentifier(waytags)
			if isPaved(waytags, highwayIdentifier) {
				surface = "default=paved"
			} else {
				surface = "default=unpaved"
			}
		}
		a1 := surfaces[surface]
		surfaces[surface] = a1 + distance
	}
	fmt.Printf("\nHighway surface types:\n")
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	distanceTotal = 0
	keys = make([]string, 0, len(surfaces))
	for key := range surfaces {
		keys = append(keys, key)
		distanceTotal += surfaces[key]
	}
	sort.Strings(keys)
	for _, key := range keys {
		percent := float64(surfaces[key]) / float64(distanceTotal) * 100.0
		fmt.Printf("%-28s  %5.1f %%  %6d m\n", key, percent, surfaces[key])
	}
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	fmt.Printf("%-28s  %5.1f %%  %6d m\n", "total", 100.0, distanceTotal)

	// calculate paved, unpaved
	// ------------------------
	for index, record := range records {
		if index == 0 {
			continue // skip csv header
		}
		distance, _ := strconv.Atoi(record[Distance])
		waytags := strings.Split(record[WayTags], " ")
		highwayIdentifier := getHighwayIdentifier(waytags)
		if isPaved(waytags, highwayIdentifier) {
			highwaysPaved += distance
		} else {
			highwaysUnpaved += distance
		}
	}
	fmt.Printf("\nPaved and unpaved surfaces:\n")
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	distanceTotal = highwaysUnpaved + highwaysPaved
	percent := float64(highwaysPaved) / float64(distanceTotal) * 100.0
	fmt.Printf("%-28s  %5.1f %%  %6d m\n", "paved", percent, highwaysPaved)
	percent = float64(highwaysUnpaved) / float64(distanceTotal) * 100.0
	fmt.Printf("%-28s  %5.1f %%  %6d m\n", "unpaved", percent, highwaysUnpaved)
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	fmt.Printf("%-28s  %5.1f %%  %6d m\n", "total", 100.0, distanceTotal)

	// build smoothness map
	// --------------------
	for index, record := range records {
		if index == 0 {
			continue // skip csv header
		}
		distance, _ := strconv.Atoi(record[Distance])
		waytags := strings.Split(record[WayTags], " ")
		smoothnessIdentifier := "default=unknown"
		for _, waytag := range waytags {
			if strings.HasPrefix(waytag, "smoothness") {
				smoothnessIdentifier = waytag
			}
		}
		a1 := smoothness[smoothnessIdentifier]
		smoothness[smoothnessIdentifier] = a1 + distance
	}
	fmt.Printf("\nSmoothness types:\n")
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	keys = make([]string, 0, len(smoothness))
	distanceTotal = 0
	for key := range smoothness {
		keys = append(keys, key)
		distanceTotal += smoothness[key]
	}
	sort.Strings(keys)
	for _, key := range keys {
		percent := float64(smoothness[key]) / float64(distanceTotal) * 100.0
		fmt.Printf("%-28s  %5.1f %%  %6d m\n", key, percent, smoothness[key])
	}
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	fmt.Printf("%-28s  %5.1f %%  %6d m\n", "total", 100.0, distanceTotal)

	// build mtb difficulties (mtb:scale) map
	// --------------------------------------
	for index, record := range records {
		if index == 0 {
			continue // skip csv header
		}
		distance, _ := strconv.Atoi(record[Distance])
		waytags := strings.Split(record[WayTags], " ")
		mtbScale := "default=unknown"
		for _, waytag := range waytags {
			if strings.HasPrefix(waytag, "mtb:scale") {
				mtbScale = waytag
			}
		}
		a1 := mtbScales[mtbScale]
		mtbScales[mtbScale] = a1 + distance
	}
	fmt.Printf("\nMountain bike difficulties (mtb:scale):\n")
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	keys = make([]string, 0, len(mtbScales))
	distanceTotal = 0
	for key := range mtbScales {
		keys = append(keys, key)
		distanceTotal += mtbScales[key]
	}
	sort.Strings(keys)
	for _, key := range keys {
		percent := float64(mtbScales[key]) / float64(distanceTotal) * 100.0
		fmt.Printf("%-28s  %5.1f %%  %6d m\n", key, percent, mtbScales[key])
	}
	fmt.Printf("%s\n", strings.Repeat("-", 47))
	fmt.Printf("%-28s  %5.1f %%  %6d m\n", "total", 100.0, distanceTotal)

	// build hiking difficulties (sac_scale) map
	// -----------------------------------------
	for index, record := range records {
		if index == 0 {
			continue // skip csv header
		}
		distance, _ := strconv.Atoi(record[Distance])
		waytags := strings.Split(record[WayTags], " ")
		sacScale := "default=unknown"
		for _, waytag := range waytags {
			if strings.HasPrefix(waytag, "sac_scale") {
				sacScale = waytag
			}
		}
		a1 := sacScales[sacScale]
		sacScales[sacScale] = a1 + distance
	}
	fmt.Printf("\nHiking difficulties (sac_scale):\n")
	fmt.Printf("%s\n", strings.Repeat("-", 57))
	keys = make([]string, 0, len(sacScales))
	distanceTotal = 0
	for key := range sacScales {
		keys = append(keys, key)
		distanceTotal += sacScales[key]
	}
	sort.Strings(keys)
	for _, key := range keys {
		percent := float64(sacScales[key]) / float64(distanceTotal) * 100.0
		fmt.Printf("%-38s  %5.1f %%  %6d m\n", key, percent, sacScales[key])
	}
	fmt.Printf("%s\n", strings.Repeat("-", 57))
	fmt.Printf("%-38s  %5.1f %%  %6d m\n", "total", 100.0, distanceTotal)

	fmt.Printf("\n")
	os.Exit(0)
}

/*
csvReadAll reads all records from csv file.
*/
func csvReadAll(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // delimiter: tab instead of comma
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

/*
isPaved determines if a highway is paved. See https://wiki.openstreetmap.org/wiki/Key:surface (as of 09/2023).
*/
func isPaved(waytags []string, highway string) bool {
	prefix := "surface="
	for _, waytag := range waytags {
		if strings.HasPrefix(waytag, prefix) {
			surface, _ := strings.CutPrefix(waytag, prefix)
			switch surface {
			case "paved", "asphalt", "chipseal", "concrete",
				"concrete:lanessurface", "concrete:plates", "paving_stonessurface", "sett",
				"unhewn_cobblestone", "cobblestone", "cobblestone:flattened", "brick",
				"metal", "wood", "stepping_stones", "rubber":
				return true // paved
			case "unpaved", "compacted", "fine_gravel", "gravel",
				"shells", "rock", "pebblestone", "ground",
				"dirt", "earth", "grass", "grass_paver",
				"metal_grid", "mud", "sand", "woodchips",
				"snow", "ice", "salt":
				return false // unpaved
			}
		}
	}
	// no surface tag found, use default for highway type
	switch highway {
	case "highway=track", "highway=track (grade2)", "highway=track (grade3)", "highway=track (grade4)", "highway=track (grade5)",
		"highway=service (grade2)", "highway=service (grade3)", "highway=service (grade4)", "highway=service (grade5)",
		"highway=path":
		return false // unpaved
	default:
		return true // paved
	}
}

/*
getHighwayIdentifier gets highway identifier from list of waytags.
*/
func getHighwayIdentifier(waytags []string) string {
	highwayIdentifier := ""
	for _, waytag := range waytags {
		if strings.HasPrefix(waytag, "highway") {
			highwayIdentifier = waytag
			if waytag == "highway=track" {
				for _, waytag := range waytags { // find tracktype
					prefix := "tracktype="
					if strings.HasPrefix(waytag, prefix) {
						grade, _ := strings.CutPrefix(waytag, prefix)
						highwayIdentifier = highwayIdentifier + " (" + grade + ")"
						break
					}
				}
			}
		}
	}
	return highwayIdentifier
}

/*
printUsage prints the usage of this program.
*/
func printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s csvfile\n", progName)
	fmt.Printf("\nExample:\n")
	fmt.Printf("  %s riesenbeck.csv\n", progName)
	fmt.Printf("\nArguments:\n")
	fmt.Printf("  csvfile = name of csv file (from brouter.m11n.de)\n")
	fmt.Printf("\n")
	os.Exit(1)
}
