# tlsh_foobar
Fuzzy matching webserver based on TLSH library.

Tool is build based on Trend Micro Locality Sensitive Hash (https://github.com/trendmicro/tlsh) and the
Golang Lib from https://github.com/glaslos/tlsh.
This tool would not work without the public available data from MalwareBazaar - https://bazaar.abuse.ch/!

The idea behind this tool is to send "new-file" creation events from a monitored directory (with the TLSH) to the tlshServer and get some results back if the distance is <150.

For detailed information about TLSH please see https://github.com/trendmicro/tlsh, for more usecases you may also watch Enhancing Malware Code Similarity Detection through Vectorsearch and TLSH on https://www.youtube.com/watch?v=W_r6Unpr8ZA.

## Change Log

v0.2 - small refactor and updated README
v0.1 - initial version

## Usage

Pre-Requirements:
- Download the CSV File from [https://bazaar.abuse.ch/export/csv/full/](https://bazaar.abuse.ch/export/#csv)
- Build the tlshServer with make `make build` or with
  `CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -ldflags "-s -w" -o build/${SERVER_BINARY_NAME} server/main.go`
- Start the tlshServer: `./build/tlshServer -csv path_to_downloaded_csv_file`

Here's a breakdown of each part for interacting with:

"GET /distance" `curl -i http://localhost:8080/distance?q=TXXXX`
- Calculates distances for a query parameter (q) using tlshSvc.Distance(query).
- Returns the result if < 150 as JSON.

Example Output
```
$ curl -i "http://localhost:8080/distance?q=65A4BF181BB98C13F54BA6BAC4D942C9E2FCD57B8907F759D41129D60F0ABA7AC023C7"
HTTP/1.1 200 OK
Date: Tue, 01 Apr 2025 12:23:36 GMT
Content-Type: text/plain; charset=utf-8
Transfer-Encoding: chunked

[{"distance":146,"signature":" Gafgyt"},{"distance":149,"signature":" SnakeKeylogger"},{"distance":139,"signature":" n/a"},{"distance":148,"signature":" Mirai"},{"distance":149,"signature":" StrelaStealer"},{"distance":146,"signature":" Mirai"},{"distance":144,"signature":" Mirai"},{"distance":146,"signature":" RemcosRAT"},   
```

"GET /search" `curl -i http://localhost:8080/search?q=TXXXX`
- Searches for tlsh items based on a query parameter (q).
- Uses tlshSvc.Search(query) to fetch matching results and returns them as JSON.


"GET /tlsh"
`curl -i http://localhost:8080/tlsh`
- Retrieves all tlsh items managed by the service.
- Uses tlshSvc.GetAll() to fetch data and returns it as JSON.
  (not recommended with a big-file)

"POST /tlsh" `curl -i -X POST --data '{"name":"TEST","hash":"T13XXXX","signature":"TEST"}'  http://localhost:8080/tlsh`
- Adds a new tlsh item to the service.
- Decodes the request body into a models.Item structure.
  (warning data is not persistent)



