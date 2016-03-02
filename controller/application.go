package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/cleesmith/pakquery/search"
	"github.com/cleesmith/pakquery/shared/view"
)

type AnyU2Record struct {
	Id                string    `json:"_id"`
	Index             string    `json:"_index"`
	Label             string    `json:"label"`
	RecordType        string    `json:"record_type"`
	Timestamp         time.Time `json:"@timestamp"`
	SensorId          int64     `json:"sensor_id"`
	EventId           int64     `json:"event_id"`
	EventSecond       int64     `json:"event_second"`
	EventMicrosecond  int64     `json:"event_microsecond"`
	GeneratorId       int64     `json:"generator_id"`
	SignatureId       int64     `json:"signature_id"`
	SrcIP             string    `json:"src_ip"`
	SPort             int64     `json:"sport"`
	DstIP             string    `json:"dst_ip"`
	DPort             int64     `json:"dport"`
	Protocol          int64     `json:"protocol"`
	Signature         string    `json:"signature"`
	PacketSecond      int64     `json:"packet_second"`
	PacketMicrosecond int64     `json:"packet_microsecond"`
	PacketDump        string    `json:"packet_dump"`
	EventType         int64     `json:"event_type"`
	EventLength       int64     `json:"event_length"`
	XType             int64     `json:"extradata_type"`
	XDataType         int64     `json:"extradata_data_type"`
	XDataLength       int64     `json:"extradata_data_length"`
	XData             string    `json:"extradata_data"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "home/index"
	// v.Vars["baseURI"] = v.BaseURI

	client, err := search.EsClient()
	if err != nil {
		log.Printf("controller: application#Index: EsClient failed: Elasticsearch error: %s\n", err)
		ErrorES(err.Error(), w, r)
		return
	}

	v.Vars["sigBucketCount"] = 0
	v.Vars["sigTerm"] = "Signatures"
	aField := "signature.raw"
	aggResult, err := search.EsTermsAgg(client, 10, aField)
	if aggResult != nil {
		v.Vars["sigBucketCount"] = len(aggResult.Buckets)
		v.Vars["sigBuckets"] = aggResult.Buckets
	} else {
		v.Vars["sigMsg"] = "No Signatures found!"
	}

	v.Vars["sipBucketCount"] = 0
	v.Vars["sipTerm"] = "Source IPs"
	aField = "src_ip.raw"
	aggResult, err = search.EsTermsAgg(client, 10, aField)
	if aggResult != nil {
		v.Vars["sipBucketCount"] = len(aggResult.Buckets)
		v.Vars["sipBuckets"] = aggResult.Buckets
	} else {
		v.Vars["sipMsg"] = "No Source IPs found!"
	}

	v.Vars["dipBucketCount"] = 0
	v.Vars["dipTerm"] = "Destination IPs"
	aField = "dst_ip.raw"
	aggResult, err = search.EsTermsAgg(client, 10, aField)
	if aggResult != nil {
		v.Vars["dipBucketCount"] = len(aggResult.Buckets)
		v.Vars["dipBuckets"] = aggResult.Buckets
	} else {
		v.Vars["dipMsg"] = "No Destination IPs found!"
	}

	v.Vars["sccBucketCount"] = 0
	v.Vars["sccTerm"] = "Source IP Country"
	aField = "src_country_code.raw"
	aggResult, err = search.EsTermsAgg(client, 10, aField)
	if aggResult != nil {
		v.Vars["sccBucketCount"] = len(aggResult.Buckets)
		v.Vars["sccBuckets"] = aggResult.Buckets
	} else {
		v.Vars["sccMsg"] = "No Source IP Countries found!"
	}

	v.Vars["dccBucketCount"] = 0
	v.Vars["dccTerm"] = "Destination IP Country"
	aField = "dst_country_code.raw"
	aggResult, err = search.EsTermsAgg(client, 10, aField)
	if aggResult != nil {
		v.Vars["dccBucketCount"] = len(aggResult.Buckets)
		v.Vars["dccBuckets"] = aggResult.Buckets
	} else {
		v.Vars["dccMsg"] = "No Destination IP Countries found!"
	}

	v.Render(w)
}
