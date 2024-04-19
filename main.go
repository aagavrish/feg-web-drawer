package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"web_drawer/utils"
)

const (
	databaseConfigPath = "configs/database.json"
	serverConfigPath   = "configs/server.json"
	sqlPath            = "sql/"
)

var (
	dbConfig     utils.DatabaseConfig
	serverConfig utils.ServerConfig
	sqlFiles     map[string]string
	dsn          string
)

type PageData struct {
	SVG template.HTML
}

func handler(w http.ResponseWriter, r *http.Request) {
	database, err := utils.ConnectToDatabase(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	nodes, err := utils.GetNodes(database, sqlFiles["select_nodes.sql"])
	if err != nil {
		log.Fatal(err)
	}

	elements, err := utils.GetElements(database, sqlFiles["select_elements.sql"])
	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer

	if err := utils.CreateSVG(nodes, elements, &buffer); err != nil {
		log.Fatal(err)
	}

	pd := PageData{
		SVG: template.HTML(buffer.String()),
	}

	tmp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmp.Execute(w, pd)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error

	if err = utils.LoadConfig(databaseConfigPath, &dbConfig); err != nil {
		log.Fatal(err)
	}

	if err = utils.LoadConfig(serverConfigPath, &serverConfig); err != nil {
		log.Fatal(err)
	}

	dsn = utils.CreateDataSourceName(dbConfig)
	sqlFiles, err = utils.LoadSqlFiles(sqlPath)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", serverConfig.Hostname, serverConfig.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
