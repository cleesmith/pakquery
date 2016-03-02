package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cleesmith/pakquery/search"
	"github.com/cleesmith/pakquery/shared/view"
)

type EventU2Record struct {
	RecordType               string    `json:"record_type"`
	Timestamp                time.Time `json:"@timestamp"`
	IndexedAt                time.Time `json:"indexed_at"`
	SensorId                 int64     `json:"sensor_id"`
	SensorHostname           string    `json:"sensor_hostname"`
	SensorInterface          string    `json:"sensor_interface"`
	SensorType               string    `json:"sensor_type"`
	EventId                  int64     `json:"event_id"`
	EventSecond              int64     `json:"event_second"`
	InputType                string    `json:"input_type"`
	Source                   string    `json:"source"`
	SourceOffset             int64     `json:"source_offset"`
	EventMicrosecond         int64     `json:"event_microsecond"`
	ClassificationId         int64     `json:"classification_id,omitempty"`
	Priority                 int64     `json:"priority"`
	GeneratorId              int64     `json:"generator_id"`
	SignatureId              int64     `json:"signature_id"`
	SrcIP                    string    `json:"src_ip"`
	SrcIPv6                  string    `json:"src_ipv6,omitempty"`
	SPort                    int64     `json:"sport"`
	SrcCountryCode           string    `json:"src_country_code,omitempty"`
	SrcLocation              string    `json:"src_location,omitempty"`
	DstIP                    string    `json:"dst_ip"`
	DstIPv6                  string    `json:"dst_ipv6,omitempty"`
	DPort                    int64     `json:"dport"`
	DstCountryCode           string    `json:"dst_country_code,omitempty"`
	DstLocation              string    `json:"dst_location,omitempty"`
	Protocol                 int64     `json:"protocol"`
	Signature                string    `json:"signature"`
	SignatureRevision        int64     `json:"signature_revision"`
	RuleRaw                  string    `json:"rule_raw"`
	RuleSourceFile           string    `json:"rule_source_file"`
	RuleSourceFileLineNumber int64     `json:"rule_source_file_line_number"`
	Blocked                  int64     `json:"blocked,omitempty"`
	Impact                   int64     `json:"impact,omitempty"`
	ImpactFlag               int64     `json:"impact_flag,omitempty"`
	MplsLabel                int64     `json:"mpls_label,omitempty"`
	VlanId                   int64     `json:"vlan_id,omitempty"`
}

type PacketU2Record struct {
	RecordType        string    `json:"record_type"`
	Timestamp         time.Time `json:"@timestamp"`
	IndexedAt         time.Time `json:"indexed_at"`
	SensorId          int64     `json:"sensor_id"`
	SensorHostname    string    `json:"sensor_hostname"`
	SensorInterface   string    `json:"sensor_interface"`
	SensorType        string    `json:"sensor_type"`
	EventId           int64     `json:"event_id"`
	EventSecond       int64     `json:"event_second"`
	InputType         string    `json:"input_type"`
	Source            string    `json:"source"`
	SourceOffset      int64     `json:"source_offset"`
	PacketSecond      int64     `json:"packet_second"`
	PacketMicrosecond int64     `json:"packet_microsecond"`
	PacketDump        string    `json:"packet_dump"`
	// There are lot's of other fields, but their
	// absence/presence depend on the packet
	// and it's layers. So it is impossible
	// to list all of them. However, the
	// PacketDump field does contain all of the
	// fields for each layer in the packet, and
	// it's in both hex and human readable format.
}

type ExtraDataU2Record struct {
	RecordType   string    `json:"record_type"`
	Timestamp    time.Time `json:"@timestamp"`
	IndexedAt    time.Time `json:"indexed_at"`
	SensorId     int64     `json:"sensor_id"`
	EventId      int64     `json:"event_id"`
	EventSecond  int64     `json:"event_second"`
	InputType    string    `json:"input_type"`
	Source       string    `json:"source"`
	SourceOffset int64     `json:"source_offset"`
	EventType    int64     `json:"event_type,omitempty"`
	EventLength  int64     `json:"event_length,omitempty"`
	XType        int64     `json:"extradata_type,omitempty"`
	XDataType    int64     `json:"extradata_data_type,omitempty"`
	XDataLength  int64     `json:"extradata_data_length,omitempty"`
	XData        string    `json:"extradata_data,omitempty"`
}

var EventRecs []EventU2Record
var PacketRecs []PacketU2Record
var ExtraDataRecs []ExtraDataU2Record

type Entity interface {
	UnmarshalEsSource(json.RawMessage) error
}

func GetEntity(jr json.RawMessage, v Entity) error {
	return v.UnmarshalEsSource(jr)
}

func (e *EventU2Record) UnmarshalEsSource(jr json.RawMessage) error {
	if err := json.Unmarshal(jr, &e); err != nil {
		log.Printf("controller: records#UnmarshalEsSource: json Unmarshal failed: error: %s\n", err)
		return err
	}
	EventRecs = append(EventRecs, *e)
	return nil
}

func (p *PacketU2Record) UnmarshalEsSource(jr json.RawMessage) error {
	if err := json.Unmarshal(jr, &p); err != nil {
		log.Printf("controller: records#UnmarshalEsSource: json Unmarshal failed: error: %s\n", err)
		return err
	}
	PacketRecs = append(PacketRecs, *p)
	return nil
}

func (ed *ExtraDataU2Record) UnmarshalEsSource(jr json.RawMessage) error {
	if err := json.Unmarshal(jr, &ed); err != nil {
		log.Printf("controller: records#UnmarshalEsSource: json Unmarshal failed: error: %s\n", err)
		return err
	}
	ExtraDataRecs = append(ExtraDataRecs, *ed)
	return nil
}

func RecordShowGET(w http.ResponseWriter, r *http.Request) {
	// clear/reset the record arrays, otherwise every request continues to append:
	EventRecs = nil
	PacketRecs = nil
	ExtraDataRecs = nil

	v := view.New(r)
	v.Name = "records/show"
	v.Vars["error"] = false

	esIndex := r.FormValue("esindex")
	esIndex = strings.Replace(esIndex, "\"", "", -1) // remove quotes
	esId := r.FormValue("esid")
	esId = strings.Replace(esId, "\"", "", -1) // remove quotes
	if len(esId) <= 0 || len(esIndex) <= 0 {
		http.Redirect(w, r, "/records", http.StatusFound)
		return
	}
	v.Vars["index"] = esIndex
	v.Vars["id"] = esId

	eventId := r.FormValue("eventid")
	eventSecond := r.FormValue("eventsecond")

	client, err := search.EsClient()
	if err != nil {
		log.Printf("controller: records#RecordShowGET: EsClient failed: Elasticsearch error: %s\n", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// get the record the user clicked on by using the ES doc "_id":
	doc, err := search.EsGetId(client, esIndex, esId)
	if err != nil {
		log.Printf("controller: records#RecordShowGET: EsGetId failed: Elasticsearch error: %s\n", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// doc is a "type GetResult struct"
	// with fields: doc.Id, doc.Index, doc.Type, and doc.Source(*json.RawMessage)
	// https://github.com/olivere/elastic/blob/33747b6d450839d8186edd535f36ae1c90e75730/get.go#L255
	recordType := doc.Type

	var u2Rec Entity
	switch recordType {
	case "event":
		u2Rec = new(EventU2Record)
	case "packet":
		u2Rec = new(PacketU2Record)
	case "extradata":
		u2Rec = new(ExtraDataU2Record)
	}
	err = GetEntity(*doc.Source, u2Rec)
	if err != nil {
		log.Printf("controller: records#RecordShowGET: GetEntity: error: %s\n", err)
	}

	// link to other records using eventId, eventSecond, and searchInDocType
	var searchInDocType string
	switch recordType {
	case "event":
		searchInDocType = "packet"
	case "packet", "extradata":
		searchInDocType = "event"
	}
	q := fmt.Sprintf("event_id:%s AND event_second:%s", eventId, eventSecond)
	// *********************************************************************************
	// search.EsQuery should NOT use a specific index, rather it should always
	// use "unifiedbeat-*" coz linked records do NOT have to be in the same index
	// *********************************************************************************
	searchResult, err := search.EsQuery(client, searchInDocType, q, "", false)
	if err != nil {
		log.Printf("controller: records#RecordShowGET: GetEntity: error: %s\n", err)
		v.Vars["error"] = true
		v.Vars["errMsg"] = err.Error()
	}

	if searchResult != nil && searchResult.Hits != nil {
		v.Vars["totalHits"] = searchResult.TotalHits()
		v.Vars["maxSearchResults"] = search.EsMaxSearchResults()
		for _, hit := range searchResult.Hits.Hits {
			// "hit" is a "type SearchHit struct"
			// with fields: hit.Id, hit.Index, hit.Type, and hit.Source(*json.RawMessage)
			// https://github.com/olivere/elastic/blob/33747b6d450839d8186edd535f36ae1c90e75730/search.go#L360

			// WTF?
			// read this: http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go
			// In short, each record type represented in the "json.RawMessage" should
			// be responsible for unmarshalling itself from json into a struct. Basically,
			// the "json.RawMessage" contains the "_source" field for a document from
			// ElasticSearch's response.
			// This may going too far since there are only 3 unified2 record types,
			// and it is very unlikely that the unified2 format will ever change.
			var u2Rec Entity
			switch hit.Type {
			case "event":
				u2Rec = new(EventU2Record)
			case "packet":
				u2Rec = new(PacketU2Record)
			case "extradata":
				u2Rec = new(ExtraDataU2Record)
			}
			err := GetEntity(*hit.Source, u2Rec)
			if err != nil {
				log.Printf("controller: records#RecordShowGET: GetEntity: error: %s\n", err)
			}
		}
	}

	v.Vars["recordType"] = recordType
	v.Vars["EventRecs"] = EventRecs
	v.Vars["PacketRecs"] = PacketRecs
	v.Vars["ExtraDataRecs"] = ExtraDataRecs

	v.Render(w)
}

func RecordsPOST(w http.ResponseWriter, r *http.Request) {
	// this looks weird, but works ... after all, this
	// isn't a true POST rather it's a GET with a search query:
	RecordsGET(w, r)
}

func RecordsGET(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "records/index"

	q := r.FormValue("q")
	if len(q) <= 0 {
		q = "*"
	}
	v.Vars["q"] = q

	v.Vars["error"] = false
	client, err := search.EsClient()
	if err != nil {
		log.Printf("\nEsClient failed: Elasticsearch error: %s\n", err)
		return
	}

	var matches []AnyU2Record

	searchResult, err := search.EsQuery(client, "", q, "", false)
	if err != nil {
		v.Vars["error"] = true
		v.Vars["errMsg"] = err.Error()
	} else {
		v.Vars["totalHits"] = searchResult.TotalHits()
		v.Vars["maxSearchResults"] = search.EsMaxSearchResults()
		if searchResult != nil && searchResult.Hits != nil {
			for _, hit := range searchResult.Hits.Hits {
				var aRecord AnyU2Record
				// see AnyU2Record struct in "controller/application.go"
				// *********************************************************
				// the fields in "*hit.Source" (of type "*json.RawMessage")
				// are IGNORED if they are NOT DEFINED in the
				// AnyU2Record struct when unmarshalled, which is
				// useful for packet type records where there
				// are many unknown fields (depending on it's layers)
				// *********************************************************
				err := json.Unmarshal(*hit.Source, &aRecord)
				if err != nil {
					continue // bad json
				}
				aRecord.Id = hit.Id
				aRecord.Index = hit.Index
				matches = append(matches, aRecord)
			}
		}
	}
	v.Vars["matches"] = matches
	v.Vars["msg"] = "Records"
	v.Render(w)
}

// cls: keep in case we ever need to prettify some json:
// func prettifyRawJson(recordType string, data json.RawMessage) (int64, int64, string, string, error) {
// 	var dat map[string]interface{}
// 	var packetDump string
// 	if err := json.Unmarshal(data, &dat); err != nil {
// 		log.Printf("controller: records#prettifyRawJson: json Unmarshal failed: error: %s\n", err)
// 		return 0, 0, "", "", err
// 	}
// 	eventId := int64(dat["event_id"].(float64))
// 	eventSecond := int64(dat["event_second"].(float64))
// 	if recordType == "packet" {
// 		// extract "FULL PACKET DATA" (already formatted in hex/human):
// 		packetDump = dat["packet_dump"].(string)
// 		// remove fields with non-printable data:
// 		delete(dat, "packet_dump")
// 		delete(dat, "packet_payload")
// 		delete(dat, "packet_data_hex")
// 	}
// 	bePretty, err := json.MarshalIndent(dat, "", "  ")
// 	if err != nil {
// 		log.Printf("controller: records#prettifyRawJson: json MarshalIndent failed: error: %s\n", err)
// 		return 0, 0, "", "", err
// 	}
// 	return eventId, eventSecond, string(bePretty), packetDump, nil
// }
