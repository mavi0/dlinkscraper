# Dlinkscraper 

Gets statistics from CPE and pushes to InfluxDB. 

Developed for the D-Link DWR-1010 CPE.

```bash
# build
go build -o out/dlinkscraper ./cmd/dlinkscraper/...

# run
./out/dlinkscraper
```

Usage:
```text
dlinkscraper scrapes data from a D-Link router

Usage:
  dlinkscraper [flags]

Flags:
  -a, --address string         telnet server address of the dlink router (default "router_url")
  -h, --help                   help for dlinkscraper
  -n, --hostname string        hostname of the system 
  -b, --influx_bucket string   bucket on influxdb server to utilise 
  -o, --influx_org string      organisation on influxdb server to utilise 
  -t, --influx_token string    api token for influxdb server 
  -s, --influx_url string      url of influxdb server 
  -i, --interval duration      time between each poll - minimum is 20 secs
  -p, --password string        password for the dlink router 
  -u, --username string        username for the dlink router 
  -v, --verbose                output debug logs

```
Args can be passed via command line, json (config.json in the working directory) or environment variables (see compose.yml)

Minimum duration between polls is 20 secs due to limitations on how frequently the router can be accessed. If no interval is specified, no looping will ocour. 