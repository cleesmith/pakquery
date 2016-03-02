package search

import (
	"errors"
	"log"
	"strings"

	"gopkg.in/olivere/elastic.v3" // see: https://github.com/olivere/elastic
)

type Elasticsearch struct {
	EsHostPort       string `json:"EsHostPort"`
	EsIndex          string `json:"EsIndex"`
	MaxSearchResults int    `json:"MaxSearchResults"`
}

var EsHostPort = "http://127.0.0.1:9200"
var EsIndex = "unifiedbeat-*"
var MaxSearchResults = 100

const (
	TsLayout = "2006-01-02T15:04:05.000Z"
	TopTen   = 10
)

func Configure(s Elasticsearch) {
	if len(s.EsHostPort) > 0 {
		EsHostPort = s.EsHostPort
	} else {
		log.Printf("search#Configure: 'EsHostPort' setting is missing in 'config.json' default: %v", EsHostPort)
	}
	if len(s.EsIndex) > 0 {
		EsIndex = s.EsIndex
	} else {
		log.Printf("search#Configure: 'EsIndex' setting is missing in 'config.json' default: %v", EsIndex)
	}
	if s.MaxSearchResults > 0 {
		MaxSearchResults = s.MaxSearchResults
	} else {
		log.Printf("search#Configure: 'MaxSearchResults' setting is missing in 'config.json' default: %v", MaxSearchResults)
	}
}

func EsMaxSearchResults() int {
	return MaxSearchResults
}

func EsClient() (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(EsHostPort))
}

func EsPing(client *elastic.Client) (*elastic.PingResult, int, error) {
	// note: this is really pinging a node so "EsHostPort" is required
	info, code, err := client.Ping(EsHostPort).Do()
	if err != nil {
		return nil, 0, err
	}
	return info, code, err
}

func EsGetId(client *elastic.Client,
	esIndex string, esId string) (*elastic.GetResult, error) {
	aDoc, err := client.Get().
		Index(esIndex).
		Id(esId).
		Do()
	if err != nil {
		return nil, err
	}
	if aDoc.Found {
		return aDoc, nil
	}
	return nil, nil
}

func EsQuery(
	client *elastic.Client, docType string, q string,
	sortBy string, sortAsc bool) (*elastic.SearchResult, error) {

	// search via string query: lucene AND/OR/NOT/etc.
	// stringQuery := elastic.NewQueryStringQuery(`packet_payload:localhost AND event_id:2`)
	// stringQuery := elastic.NewQueryStringQuery(`192.168.1.102 AND ftp AND ACK=true`)
	// stringQuery := elastic.NewQueryStringQuery(`"5f 25 25 7c"`)
	// stringQuery := elastic.NewQueryStringQuery(`event_id:200`)
	// stringQuery := elastic.NewQueryStringQuery(`*`) // may NOT be empty or blank !!!

	stringQuery := elastic.NewQueryStringQuery(q)
	if len(strings.TrimSpace(q)) <= 0 {
		// "q" param may not be empty or blank
		stringQuery = elastic.NewQueryStringQuery("*")
	}
	// test error: stringQuery = elastic.NewQueryStringQuery("NOT AND OR")
	if len(strings.TrimSpace(sortBy)) <= 0 {
		sortBy = "@timestamp"
	}
	// see: https://github.com/olivere/elastic/blob/release-branch.v3/search.go#L315
	// for "searchResult" struct:
	// type SearchResult struct {
	//  TookInMillis int64         `json:"took"`         // search time in milliseconds
	//  ScrollId     string        `json:"_scroll_id"`   // only used with Scroll and Scan operations
	//  Hits         *SearchHits   `json:"hits"`         // the actual search hits
	//  Suggest      SearchSuggest `json:"suggest"`      // results from suggesters
	//  Aggregations Aggregations  `json:"aggregations"` // results from aggregations
	//  TimedOut     bool          `json:"timed_out"`    // true if the search timed out
	//  Error *ErrorDetails `json:"error,omitempty"` // only used in MultiGet
	// }
	searchResult, err := client.Search().
		Index(EsIndex).
		Type(docType).
		Query(stringQuery).
		Size(MaxSearchResults).
		Sort(sortBy, sortAsc).
		Do()
	if err != nil {
		elasticErr, ok := err.(*elastic.Error)
		if !ok {
			return nil, err
		}
		errReason := elasticErr.Details.RootCause[0].Reason
		if len(errReason) <= 0 {
			errReason = err.Error()
		}
		return nil, errors.New(errReason)
	}
	return searchResult, err
}

func EsTermsAgg(
	client *elastic.Client, maxAggResults int,
	term string) (*elastic.AggregationBucketKeyItems, error) {

	matchAll := elastic.NewMatchAllQuery()
	termAgg := elastic.NewTermsAggregation().Field(term).Size(maxAggResults).OrderByCountDesc()
	aQuery := client.Search().Index(EsIndex).Query(matchAll).Size(0)
	aQuery = aQuery.Aggregation(term, termAgg)
	searchResult, err := aQuery.Do()
	if err != nil {
		return nil, err
	}
	agg := searchResult.Aggregations
	if agg == nil {
		return nil, err
	}
	termsAggResult, ok := agg.Terms(term)
	if ok {
		return termsAggResult, err
	}

	return nil, err
}
