package main

import (
	"fmt"
	importBfp "github.com/helmutkemper/gosmImport"
	iotmakerDbInterface "github.com/helmutkemper/iotmaker.db.interface"
	iotmakerDbMongodb "github.com/helmutkemper/iotmaker.db.mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type Films struct {
	Film string
}

func main() {

	var db iotmakerDbInterface.DbFunctionsInterface
	var err error
	var data = make([]map[string]interface{}, 0)

	db = &iotmakerDbMongodb.DbFunctions{}
	err = db.Connect("mongodb://0.0.0.0:27017", "geo", []string{"point", "way"})
	if err != nil {
		log.Fatalf("db.connection.error: %v", err.Error())
	}

	importBfp.ProcessPbfFileInMemory()
}
