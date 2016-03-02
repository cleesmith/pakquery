# pakquery

> Go web app for searching unified2 records indexed by [Unifiedbeat](https://github.com/cleesmith/unifiedbeat) into ElasticSearch.

***

##### Usage

* ```cd $GOPATH```
* ```git clone https://github.com/cleesmith/pakquery```
* ```cd pakquery```
* ```go build```
  * or cross-compile for linux ```env GOOS=linux GOARCH=amd64 go build```
* if remote deploy:
  * copy binary: ```scp pakquery user@host:/home/user/folder```
  * copy css/js: ```scp -r static user@host:/home/user/folder```
  * copy views: ```scp -r views user@host:/home/user/folder```
  * ```scp config.json user@host:/home/user/folder```
* edit ```config.json```
* **./pakquery**
* browse to ```http://host:8080/```

##### Features

* fast access to the unified2 records as indexed into Elasticsearch using Unifiedbeat
* Overview page with **top** counts linked to the matching records
* Records page
  * search via URL or the simple form
  * search queries are similar to Kibana's discover feature
  * click a record to get the full details

***
***
