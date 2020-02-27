package main

import (
	"fmt"
	iotmakerDbInterface "github.com/helmutkemper/iotmaker.db.interface"
	iotmakerDbMongodb "github.com/helmutkemper/iotmaker.db.mongodb"
	iotmaker_geo_osm "github.com/helmutkemper/iotmaker.geo.osm"
	iotmaker_geo_pbf_import "github.com/helmutkemper/iotmaker.geo.pbf.import"
	"log"
	"time"
)

var db iotmakerDbInterface.DbFunctionsInterface

func main() {
	var err error

	importMap := iotmaker_geo_pbf_import.Import{}
	importMap.DontFindDuplicatedId = true

	db = &iotmakerDbMongodb.DbFunctions{}
	err = db.Connect("mongodb://0.0.0.0:27017", "geo", []string{"point", "way", "polygon", "surrounding", "surroundingRight", "surroundingLeft"})
	if err != nil {
		log.Fatalf("db.connection.error: %v", err.Error())
	}

	start := time.Now()
	dirPath := "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187763"

	importMap = iotmaker_geo_pbf_import.Import{}
	err = importMap.SetDirFromBinaryFilesCache(dirPath + "/testBin/")
	if err != nil {
		panic(err)
	}

	err = importMap.SetMapFilePath(dirPath + "/portugal-latest.osm.pbf")
	if err != nil {
		panic(err)
	}

	err = importMap.CountElements()
	if err != nil {
		panic(err)
	}

	//err = importMap.ExtractNodesToBinaryFilesDir()
	//if err != nil {
	//  panic( err )
	//}

	//err = importMap.FindAllNodesForTest()
	//if err != nil {
	//  panic( err )
	//}

	err = importMap.ProcessWaysFromMapFile(processWayFunctionPointer)
	if err != nil {
		panic(err)
	}

	fmt.Printf("duração: %v\n", time.Since(start))
	fmt.Println("terminou ok")

}

var dis = iotmaker_geo_osm.DistanceStt{}
var disMin = iotmaker_geo_osm.DistanceStt{}

func init() {
	dis.SetMeters(100)
	disMin.SetMeters(50)
}

func processWayFunctionPointer(wayConverted iotmaker_geo_pbf_import.WayConverted) {
	var err error

	wayToDb := iotmaker_geo_osm.WayStt{}
	polygonSurroundingToDb := iotmaker_geo_osm.PolygonStt{}
	polygonSurroundingRightToDb := iotmaker_geo_osm.PolygonStt{}
	polygonSurroundingLeftToDb := iotmaker_geo_osm.PolygonStt{}

	for k := range wayConverted.Node {
		err = wayToDb.AddLngLatDegrees(wayConverted.Node[k][0], wayConverted.Node[k][1])
		if err != nil {
			panic(err)
		}
	}

	for key, value := range wayConverted.Tags {
		wayToDb.AddTag(key, value)
	}

	wayToDb.SetId(wayConverted.ID)
	err = wayToDb.Init()
	if err != nil {
		panic(err)
	}

	wayToDb.MakeGeoJSonFeature()
	err = db.Insert("way", wayToDb)
	if err != nil {
		panic(err)
	}

	return
	if len(wayToDb.Loc) < 3 {
		return
	}

	err, polygonSurroundingToDb = wayToDb.MakePolygonSurroundings(dis, disMin)
	if err != nil {
		panic(err)
	}

	polygonSurroundingToDb.MakeGeoJSonFeature()
	err = db.Insert("surrounding", wayToDb)
	if err != nil {
		panic(err)
	}

	err, polygonSurroundingLeftToDb = wayToDb.MakePolygonSurroundingsLeft(dis, disMin)
	if err != nil {
		panic(err)
	}

	polygonSurroundingLeftToDb.MakeGeoJSonFeature()
	err = db.Insert("surroundingLeft", wayToDb)
	if err != nil {
		panic(err)
	}

	err, polygonSurroundingRightToDb = wayToDb.MakePolygonSurroundingsRight(dis, disMin)
	if err != nil {
		panic(err)
	}

	polygonSurroundingRightToDb.MakeGeoJSonFeature()
	err = db.Insert("surroundingRight", wayToDb)
	if err != nil {
		panic(err)
	}

	//fazer:
	//visible

}
