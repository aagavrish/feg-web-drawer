package utils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DBname   string `json:"dbname"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type ServerConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type Node struct {
	ID int16
	X  float64
	Y  float64
}

type Element struct {
	ID      int16
	NodesID []int16
}

func LoadConfig(path string, config interface{}) error {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonFile, &config); err != nil {
		return err
	}

	return nil
}

func CreateDataSourceName(config DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Hostname, config.Port, config.DBname)
}

func ConnectToDatabase(dsn string) (*sql.DB, error) {
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

func LoadSqlFiles(dir string) (map[string]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	mapSqlFiles := make(map[string]string, len(files))
	for _, file := range files {
		fileContent, err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		mapSqlFiles[file.Name()] = string(fileContent)
	}

	return mapSqlFiles, nil
}

func GetNodes(database *sql.DB, sql string) ([]Node, error) {
	rows, err := database.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []Node
	for rows.Next() {
		var node Node

		if err := rows.Scan(&node.ID, &node.X, &node.Y); err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func GetElements(database *sql.DB, sql string) ([]Element, error) {
	rows, err := database.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements []Element
	for rows.Next() {
		var element Element
		var n1, n2, n3 int16

		if err := rows.Scan(&element.ID, &n1, &n2, &n3); err != nil {
			return nil, err
		}

		element.NodesID = []int16{n1, n2, n3}
		elements = append(elements, element)
	}

	return elements, nil
}

func drawElements(canvas *svg.SVG, elements []Element, nodes []Node, scale int) {
	for _, element := range elements {
		for i := 0; i < len(element.NodesID); i++ {
			startNode := nodes[element.NodesID[i]-1]
			endNode := nodes[element.NodesID[(i+1)%len(element.NodesID)]-1]
			canvas.Line(int(startNode.X)*scale, -int(startNode.Y)*scale, int(endNode.X)*scale, -int(endNode.Y)*scale, `style="stroke:grey;stroke-width:3"`)
		}
	}
}

func drawNodes(canvas *svg.SVG, nodes []Node, scale int) {
	for i, node := range nodes {
		canvas.Circle(int(node.X)*scale, -int(node.Y)*scale, 8, `fill="lightgrey" stroke="black" stroke-width="2"`)
		canvas.Text(int(node.X)*scale, -int(node.Y)*scale+4, fmt.Sprintf("%d", i+1), `text-anchor="middle" fill="black" font-size="10" font-weight="bold"`)
	}
}

func CreateSVG(nodes []Node, elements []Element, buf *bytes.Buffer) error {
	const (
		width      = 1000
		height     = 600
		translateX = 500
		translateY = 300
		scale      = 5
	)

	canvas := svg.New(buf)
	canvas.Start(width, height)
	canvas.Translate(translateX, translateY)

	drawElements(canvas, elements, nodes, scale)
	drawNodes(canvas, nodes, scale)

	canvas.Text(0, height/3, fmt.Sprintf("Number of nodes: %d", len(nodes)), `text-anchor="middle" font-family="Arial" font-size="14" font-weight="bold"`)
	canvas.Text(0, height/3+20, fmt.Sprintf("Number of elements: %d", len(elements)), `text-anchor="middle" font-family="Arial" font-size="14" font-weight="bold"`)

	canvas.Gend()
	canvas.End()

	return nil
}
