package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gopistolet/gospf"
	"github.com/gopistolet/gospf/dns"
)

func main() {
	os.Exit(_main())
}

func getSpf(domain string) (*gospf.SPF, error) {
	i := 0
	var spf *gospf.SPF
	var err error
	for {
		spf, err = gospf.New(domain, &dns.GoSPFDNS{})
		if err == nil {
			break
		}
		i++
		if i > 2 {
			break
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	return spf, err
}

func getReverse(ip string) (string, error) {
	i := 0
	var result string
	var err error
	for {
		out, derr := exec.Command("dig", "+short", "-x", ip).Output()
		if derr != nil {
			// do not retry
			err = derr
			break
		}
		result = strings.TrimRight(string(out), "\n")
		result = strings.TrimRight(result, ".")
		if len(result) == 0 {
			// make error and retry
			err = fmt.Errorf("no result\n")
		} else {
			// success
			break
		}
		i++
		if i > 2 {
			break
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	return result, err

}

func _main() (st int) {
	st = 1
	if len(os.Args) < 3 || (len(os.Args) >= 2 && (os.Args[1] == "-h" || os.Args[1] == "--help")) {
		fmt.Printf("check-spf-and-reserve-lookup $ip $domain\n")
		if len(os.Args) >= 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
			st = 0
		}
		return
	}
	ip := os.Args[1]
	domain := os.Args[2]
	spf, err := getSpf(domain)
	if err != nil {
		fmt.Printf("NG: DNS lookup failed: %s\n", err)
		return
	}

	check, err := spf.CheckIP(ip)
	if err != nil {
		fmt.Printf("NG: Check SPF failed: %s\n", err)
		return
	}

	if check != "Pass" {
		fmt.Printf("NG: spf check failed: result=%s\n", check)
		return
	}

	result, err := getReverse(ip)
	if err != nil {
		fmt.Printf("NG: reverse lookup dig failed: %s\n", err)
		return
	}

	if strings.HasSuffix(result, domain) == false {
		fmt.Printf("NG: reverse lookup failed: no contains domain - %s\n", result)
		return
	}

	fmt.Printf("OK: spf:%s, reserve-lookup:%s\n", check, result)
	st = 0
	return
}
