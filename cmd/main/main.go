package main

import (
	iotmakerDbInterface "github.com/helmutkemper/iotmaker.db.interface"
	iotmakerDbMongodb "github.com/helmutkemper/iotmaker.db.mongodb"
	iotmaker_geo_pbf_import "github.com/helmutkemper/iotmaker.geo.pbf.import"
	"log"
)

type Films struct {
	Film string
}

func main() {

	var db iotmakerDbInterface.DbFunctionsInterface
	var err error

	db = &iotmakerDbMongodb.DbFunctions{}
	err = db.Connect("mongodb://0.0.0.0:27017", "geo", []string{"point", "way", "polygon"})
	if err != nil {
		log.Fatalf("db.connection.error: %v", err.Error())
	}

	iotmaker_geo_pbf_import.ProcessPbfFileInMemory(db, "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/portugal-latest.osm.pbf", "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/binMap.bin")
}
