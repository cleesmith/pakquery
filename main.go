package main

import (
	"encoding/json"
	"log"

	"github.com/cleesmith/pakquery/route"
	"github.com/cleesmith/pakquery/search"
	"github.com/cleesmith/pakquery/shared/jsonconfig"
	"github.com/cleesmith/pakquery/shared/server"
	"github.com/cleesmith/pakquery/shared/session"
	"github.com/cleesmith/pakquery/shared/view"
	"github.com/cleesmith/pakquery/shared/view/plugin"
)

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	ES       search.Elasticsearch `json:"Elasticsearch"`
	Server   server.Server        `json:"Server"`
	Session  session.Session      `json:"Session"`
	Template view.Template        `json:"Template"`
	View     view.View            `json:"View"`
}

func init() {
	// Verbose logging with file name and line number
	// log.SetFlags(log.Lshortfile)
}

func main() {
	log.Println("Loading settings from config.json file")
	// let's assume this is in the working dir:
	jsonconfig.Load("config.json", config)

	// Configure elasticsearch
	search.Configure(config.ES)
	log.Println("ElasticSearch settings:")
	log.Printf("\tEsHostPort: %s\n", search.EsHostPort)
	log.Printf("\tEsIndex: %s\n", search.EsIndex)
	log.Printf("\tMaxSearchResults: %v\n", search.MaxSearchResults)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape())

	// Start listening
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
