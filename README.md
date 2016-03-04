# pakquery

> Go web app for searching unified2 records indexed by [Unifiedbeat](https://github.com/cleesmith/unifiedbeat) into ElasticSearch.
>
> Built using Go 1.6 and tested on Ubuntu 14.04.3

***

#### Installation and usage

* Download and install:
```
wget https://github.com/cleesmith/pakquery/releases/download/v1.0.0/pakquery_v1.0.0.tar.gz
tar -zxvf pakquery_v1.0.0.tar.gz
cd pakquery_v1.0.0
edit config.json
./pakquery
```
* the downloaded tarball includes:
  * **pakquery** - the executable binary for Linux 64bit (_linux amd64_)
  * **views** - html templates used by pakquery
  * **static** - folder of web assets (_mostly css and javascript files_)
    * from https://www.datatables.net/download/index including:
      * jQuery v2.2.0
      * Bootstrap v3.3.6
      * DataTables 1.10.11
      * DataTables Bootstrap 3 integration

***

#### or to build from source:

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

***

#### Features

* fast access to the unified2 records as indexed into Elasticsearch using Unifiedbeat
* an overview page with **top** counts linked to the matching records
* Records page
  * search via URL or an HTML form
  * search queries are similar to Kibana's discover feature
  * click a record to get the full details

***
***
