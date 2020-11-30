# check-mysql-uptime

Mackerel check plugin for checking MySQL uptime

## Usage

```
Usage:
  check-mysql-uptime [OPTIONS]

Application Options:
  -H, --host=     Hostname (default: localhost)
  -p, --port=     Port (default: 3306)
  -u, --user=     Username (default: root)
  -P, --password= Password
      --timeout=  Timeout to connect mysql (default: 5s)
  -c, --critical= critical if uptime seconds is less than this number
  -w, --warning=  warning if uptime seconds is less than this number
  -v, --version   Show version

Help Options:
  -h, --help      Show this help message
```

Example

```
$ ./check-mysql-uptime--user=xxx --password=xxx  -w 110 -c 110
MySQL Uptime OK: up 28 days, 22:08:40
```

  ## Install

Please download release page or `mkr plugin install kazeburo/go-check-mysql-uptime`.
