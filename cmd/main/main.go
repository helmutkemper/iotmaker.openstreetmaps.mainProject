package main

import (
	"encoding/xml"
	"fmt"
	iotmakerDbInterface "github.com/helmutkemper/iotmaker.db.interface"
	iotmakerDbMongodb "github.com/helmutkemper/iotmaker.db.mongodb"
	iotmaker_geo_pbf_import "github.com/helmutkemper/iotmaker.geo.pbf.import"
	"io/ioutil"
	"log"
)

type Films struct {
	Film string
}

func main() {

	var db iotmakerDbInterface.DbFunctionsInterface
	var err error
	var nodes, ways, relations, others int

	var lat, lon float64
	importMap := iotmaker_geo_pbf_import.Import{}
	importMap.AppendNodeToFile("./bin", 1, 1.1, 1.2)
	importMap.AppendNodeToFile("./bin", 2, 2.1, 2.2)
	importMap.AppendNodeToFile("./bin", 3, 3.1, 3.2)
	importMap.AppendNodeToFile("./bin", 4, 4.1, 4.2)
	importMap.AppendNodeToFile("./bin", 5, 5.1, 5.2)
	importMap.AppendNodeToFile("./bin", 6, 6.1, 6.2)

	err, lon, lat = importMap.FindLonLatByIdInFile("./bin", 1)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("lon: %v\n", lon)
	fmt.Printf("lat: %v\n", lat)

	err, lon, lat = importMap.FindLonLatByIdInFile("./bin", 2)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("lon: %v\n", lon)
	fmt.Printf("lat: %v\n", lat)

	err, lon, lat = importMap.FindLonLatByIdInFile("./bin", 3)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("lon: %v\n", lon)
	fmt.Printf("lat: %v\n", lat)

	err, lon, lat = importMap.FindLonLatByIdInFile("./bin", 4)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("lon: %v\n", lon)
	fmt.Printf("lat: %v\n", lat)

	err, lon, lat = importMap.FindLonLatByIdInFile("./bin", 5)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("lon: %v\n", lon)
	fmt.Printf("lat: %v\n", lat)

	err, lon, lat = importMap.FindLonLatByIdInFile("./bin", 6)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("lon: %v\n", lon)
	fmt.Printf("lat: %v\n", lat)

	err, lon, lat = importMap.FindLonLatByIdInFile("./bin", 7)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("lon: %v\n", lon)
	fmt.Printf("lat: %v\n", lat)

	return
	db = &iotmakerDbMongodb.DbFunctions{}
	err = db.Connect("mongodb://0.0.0.0:27017", "geo", []string{"point", "way", "polygon"})
	if err != nil {
		log.Fatalf("db.connection.error: %v", err.Error())
	}

	//      /media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/nodes/13787735.xml

	importMap = iotmaker_geo_pbf_import.Import{}
	err, nodes, ways, relations, others = importMap.CountElements("/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/portugal-latest.osm.pbf")
	if err != nil {
		log.Fatalf("db.connection.error: %v", err.Error())
	}
	log.Printf("nodes: %v\n", nodes)
	log.Printf("ways: %v\n", ways)
	log.Printf("relations: %v\n", relations)
	log.Printf("others: %v\n", others)

	importMap.ExtractNodes("/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/portugal-latest.osm.pbf", "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/nodesJSon")
	//iotmaker_geo_pbf_import.ProcessPbfFileInMemory(db, "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/portugal-latest.osm.pbf", "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/binMap.bin")
}

type NodesTagStt struct {
	XMLName   xml.Name     `xml:"node"`
	Id        int64        `xml:"id,attr"`
	Lat       float64      `xml:"lat,attr"`
	Lon       float64      `xml:"lon,attr"`
	Version   int64        `xml:"version,attr"`
	TimeStamp string       `xml:"timestamp,attr"`
	ChangeSet string       `xml:"changeset,attr"`
	UId       int64        `xml:"uid,attr"`
	User      string       `xml:"user,attr"`
	Tag       []TagsTagStt `xml:"tag"`
}

type TagsTagStt struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
}

type OsmNodeTagStt struct {
	XMLName   xml.Name    `xml:"osm"`
	Version   string      `xml:"version,attr"`
	Generator string      `xml:"generator,attr"`
	TimeStamp string      `xml:"timestamp,attr"`
	Node      NodesTagStt `xml:"node"`
}

func test() int64 {
	var body []byte
	var err error

	fileName := "/media/kemper/c5d4fd1f-1a7e-4bdd-8124-e2ad60e187761/nodes/13787735.xml"
	nodeRemote := OsmNodeTagStt{}

	body, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("test.error: %v", err.Error())
	}
	err = xml.Unmarshal(body, &nodeRemote)

	return nodeRemote.Node.Id
}
