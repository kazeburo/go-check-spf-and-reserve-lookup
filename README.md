# go-check-spf-and-reserve-lookup

Verify that the SPF records and reverse DNS records are correct.

## Usage

```
Usage:
  check-spf-and-reserve-lookup check-spf-and-reserve-lookup [OPTIONS] $ip $domain

Application Options:
  -v, --version  Show version

Help Options:
  -h, --help     Show this help message
```

Example

```
./check-spf-and-reserve-lookup x.x.x.x example.com
SPF and Reverse OK: OK: spf:Pass, reserve-lookup:xx.example.com
```

