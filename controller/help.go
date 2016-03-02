package controller

import (
	"log"
	"net/http"

	"github.com/cleesmith/pakquery/search"
	"github.com/cleesmith/pakquery/shared/view"
)

func HelpGET(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "help/index"

	client, err := search.EsClient()
	if err != nil {
		log.Printf("controller: help#HelpGET: EsClient failed: Elasticsearch error: %s\n", err)
		ErrorES(err.Error(), w, r)
		return
	}

	info, code, err := search.EsPing(client)
	if err != nil {
		log.Printf("\nPing failed: Elasticsearch return code: %d \t err: %s\n", code, err)
		return
	}
	v.Vars["esVersion"] = info.Version.Number
	v.Vars["luceneVersion"] = info.Version.LuceneVersion
	v.Vars["maxSearchResults"] = search.EsMaxSearchResults()

	v.Render(w)
}
