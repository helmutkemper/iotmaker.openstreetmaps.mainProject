package main

import (
	"fmt"
	iotmaker_db_geo_mongodb "github.com/helmutkemper/iotmaker.db.geo.mongodb"
	"github.com/helmutkemper/iotmaker.db.geo.mongodb/factoryGeoDbMongoDb"
	iotmaker_geo_osm "github.com/helmutkemper/iotmaker.geo.osm"
	iotmaker_geo_pbf_import "github.com/helmutkemper/iotmaker.geo.pbf.import"
	"github.com/helmutkemper/osmpbf"
	"github.com/helmutkemper/util"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sync"
	"time"
)

var db iotmaker_db_geo_mongodb.DbFunctionsFromMapInterface

/*
nodes: 5.765.970.256
ways: 639.271.137
relations: 7.514.422
*/
//{"tag.admin_level": "2"}
func main() {
	var err error

	importMap := iotmaker_geo_pbf_import.Import{}
	importMap.DontFindDuplicatedId = true

	err, db = factoryGeoDbMongoDb.NewConnection("mongodb://0.0.0.0:27017", "globo")
	if err != nil {
		log.Fatalf("db.connection.error: %v", err.Error())
	}

	start := time.Now()
	dirPath := "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187763"
	dirPath = "/home/kemper/osm"

	importMap = iotmaker_geo_pbf_import.Import{}
	err = importMap.SetDirFromBinaryFilesCache(dirPath + "/testBin/")
	//if err != nil {
	//	panic(err)
	//}

	//err = importMap.SetMapFilePath(dirPath + "/portugal-latest.osm.pbf")
	err = importMap.SetMapFilePath(dirPath + "/planet-200210.osm.pbf")
	if err != nil {
		panic(err)
	}

	//err = importMap.CountElements()
	//if err != nil {
	//	panic(err)
	//}

	//err = importMap.ExtractNodesToBinaryFilesDir()
	//if err != nil {
	//  panic( err )
	//}

	//err = importMap.FindAllNodesForTest()
	//if err != nil {
	//  panic( err )
	//}

	//err = importMap.ProcessWaysFromMapFile(functionToDecideWhetherTheWayShouldBeProcessedOrNot, processWayFunctionPointer)
	//if err != nil {
	//	panic(err)
	//}

	var step int64 = 500000
	var i int64 = 0
	var wg sync.WaitGroup

	for i = 0; i < 639271137; i += step {

		var nodesStep int64 = 5765970256 / 10
		var iNode int64 = 0

		for iNode = 0; iNode < 5765970256; iNode += nodesStep {

			wg.Add(1)
			go func() {
				startCounter := iNode
				endCounter := startCounter + nodesStep
				err = importMap.GetAllNodesFromMap(getAllNodesToPopulateWays, startCounter, endCounter)
				if err != nil {
					panic(err)
				}
				fmt.Printf("GetAllNodesFromMap().duração: %v\n", time.Since(start))
				wg.Done()
			}()

		}
		wg.Wait()
		err = importMap.GetAllWaysFromMap(getAllWaysAndPutIntoDb, i, i+step)
		if err != nil {
			panic(err)
		}
		fmt.Printf("GetAllWaysFromMap().duração: %v\n", time.Since(start))
	}

	fmt.Printf("duração: %v\n", time.Since(start))
	fmt.Println("terminou ok")

}

var dis = iotmaker_geo_osm.DistanceStt{}
var disMin = iotmaker_geo_osm.DistanceStt{}

func init() {
	dis.SetMeters(50)
	disMin.SetMeters(25)
}

func getAllNodesToPopulateWays(node osmpbf.Node) int64 {
	// lon
	//-9.525146484375
	//-6.075439453125

	// lat
	//42.27730877423709
	//36.78289206199065

	//if node.Lon < -9.52514648 || node.Lon > -6.07543945 {
	//	return
	//}

	//if node.Lat < 36.78289206 || node.Lat > 42.27730877 {
	//	return
	//}

	dataTmpWay := make([]osmpbf.Way, 0)
	var dataWay []iotmaker_geo_osm.WayStt

	node.Lon = util.Round(node.Lon, 0.5, 8.0)
	node.Lat = util.Round(node.Lat, 0.5, 8.0)

	err, count := db.WayTmpCount(bson.M{})
	if err != nil {
		panic(err)
	}

	if count == 0 {
		return count
	}

	err = db.WayTmpFind(bson.M{"nodeids": node.ID}, &dataTmpWay)
	if err != nil {
		panic(err)
	}

	for wayKey := range dataTmpWay {
		wayId := dataTmpWay[wayKey].ID
		wayNodesList := dataTmpWay[wayKey].NodeIDs

		err = db.WayToPopulateFind(bson.M{"id": wayId}, &dataWay)
		if err != nil {
			panic(err)
		}

		if len(dataWay) == 0 {
			continue
		}

		dataWayId := dataWay[0].Id

		for nodeKey, nodeId := range wayNodesList {
			if nodeId == node.ID {
				dataWay[0].Loc[nodeKey] = [2]float64{node.Lon, node.Lat}
				dataWay[0].Rad[nodeKey] = [2]float64{iotmaker_geo_osm.DegreesToRadians(node.Lon), iotmaker_geo_osm.DegreesToRadians(node.Lat)}

				err = db.WayToPopulateUpdateLocations(dataWayId, int64(nodeKey), dataWay[0].Loc[nodeKey], dataWay[0].Rad[nodeKey])
				if err != nil {
					panic(err)
				}

				pass := true
				for _, pointValue := range dataWay[0].Loc {
					if pointValue[0] == 0.0 && pointValue[1] == 0.0 {
						pass = false
						break
					}
				}

				if pass == true {
					err = db.WayToPopulateDeleteByOsmId(dataWayId)
					if err != nil {
						panic(err)
					}

					err = db.WayTmpDeleteByOsmId(dataWayId)
					if err != nil {
						panic(err)
					}

					err = dataWay[0].Init()
					if err != nil {
						panic(err)
					}

					dataWay[0].MakeGeoJSonFeature()
					err, _ = dataWay[0].MakeMD5()
					if err != nil {
						panic(err)
					}

					err = db.WayInsert(dataWay[0])
					if err != nil {
						panic(err)
					}

					//polygonSurroundingToDb := iotmaker_geo_osm.PolygonStt{}
					//polygonSurroundingRightToDb := iotmaker_geo_osm.PolygonStt{}
					//polygonSurroundingLeftToDb := iotmaker_geo_osm.PolygonStt{}
					//
					//if len(dataWay[0].Loc) < 3 {
					//	return count
					//}
					//
					//err, polygonSurroundingToDb = dataWay[0].MakePolygonSurroundings(dis, disMin)
					//if err != nil {
					//	panic(err)
					//}
					//
					//err = polygonSurroundingToDb.Init()
					//if err != nil {
					//	panic(err)
					//}
					//
					//polygonSurroundingToDb.MakeGeoJSonFeature()
					//err, _ = polygonSurroundingToDb.MakeMD5()
					//if err != nil {
					//	panic(err)
					//}
					//
					//err = db.SurroundingInsert(polygonSurroundingToDb)
					//if err != nil {
					//	panic(err)
					//}
					//
					//err, polygonSurroundingLeftToDb = dataWay[0].MakePolygonSurroundingsLeft(dis, disMin)
					//if err != nil {
					//	panic(err)
					//}
					//
					//err = polygonSurroundingLeftToDb.Init()
					//if err != nil {
					//	panic(err)
					//}
					//
					//polygonSurroundingLeftToDb.MakeGeoJSonFeature()
					//err, _ = polygonSurroundingLeftToDb.MakeMD5()
					//if err != nil {
					//	panic(err)
					//}
					//
					//err = db.SurroundingLeftInsert(polygonSurroundingLeftToDb)
					//if err != nil {
					//	panic(err)
					//}
					//
					//err, polygonSurroundingRightToDb = dataWay[0].MakePolygonSurroundingsRight(dis, disMin)
					//if err != nil {
					//	panic(err)
					//}
					//
					//err = polygonSurroundingRightToDb.Init()
					//if err != nil {
					//	panic(err)
					//}
					//
					//polygonSurroundingRightToDb.MakeGeoJSonFeature()
					//
					//err, _ = polygonSurroundingRightToDb.MakeMD5()
					//if err != nil {
					//	panic(err)
					//}
					//
					//err = db.SurroundingRightInsert(polygonSurroundingRightToDb)
					//if err != nil {
					//	panic(err)
					//}
				}
				break
			}
		}

	}

	return count
}

func getAllWaysAndPutIntoDb(way osmpbf.Way) {

	err, count := db.WayCount(bson.M{"id": way.ID})
	if err != nil {
		panic(err)
	}

	if count != 0 {
		return
	}

	if way.Info.Visible == false {
		return
	}

	err = db.WayTmpInsert(way)
	if err != nil {
		panic(err)
	}

	wayToDb := iotmaker_geo_osm.WayStt{}
	for key, value := range way.Tags {
		wayToDb.AddTag(key, value)
	}

	totalNodes := len(way.NodeIDs)
	fakeData := make([][2]float64, totalNodes)

	wayToDb.SetId(way.ID)
	wayToDb.Visible = way.Info.Visible
	wayToDb.Loc = fakeData
	wayToDb.Rad = fakeData

	err = db.WayToPopulateInsert(wayToDb)
	if err != nil {
		panic(err)
	}
}

func functionToDecideWhetherTheWayShouldBeProcessedOrNot(id int64) bool {
	err, found := db.WayCount(bson.M{"id": id})
	if err != nil {
		panic(err)
	}

	if found == 0 {
		return true
	}

	return false
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
	err = db.WayInsert(wayToDb)
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
	err = db.SurroundingInsert(wayToDb)
	if err != nil {
		panic(err)
	}

	err, polygonSurroundingLeftToDb = wayToDb.MakePolygonSurroundingsLeft(dis, disMin)
	if err != nil {
		panic(err)
	}

	polygonSurroundingLeftToDb.MakeGeoJSonFeature()
	err = db.SurroundingLeftInsert(wayToDb)
	if err != nil {
		panic(err)
	}

	err, polygonSurroundingRightToDb = wayToDb.MakePolygonSurroundingsRight(dis, disMin)
	if err != nil {
		panic(err)
	}

	polygonSurroundingRightToDb.MakeGeoJSonFeature()
	err = db.SurroundingRightInsert(wayToDb)
	if err != nil {
		panic(err)
	}

	//fazer:
	//visible

}
