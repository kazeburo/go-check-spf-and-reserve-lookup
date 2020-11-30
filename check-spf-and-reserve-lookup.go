package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gopistolet/gospf"
	"github.com/gopistolet/gospf/dns"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// Version by Makefile
var version string

type options struct {
	Version bool `short:"v" long:"version" description:"Show version"`
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
		result = strings.ReplaceAll(string(out), "\n", "")
		result = strings.ReplaceAll(result, "\r", "")
		result = strings.TrimRight(result, ".")
		if len(result) == 0 {
			// make error and retry
			err = fmt.Errorf("no result")
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

func main() {
	ckr := checkSpfAndReverse()
	ckr.Name = "SPF and Reverse"
	ckr.Exit()
}

func checkSpfAndReverse() *checkers.Checker {
	opts := options{}
	psr := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	psr.Usage = "check-spf-and-reserve-lookup [OPTIONS] $ip $domain"
	args, err := psr.Parse()
	if opts.Version {
		fmt.Fprintf(os.Stderr, "Version: %s\nCompiler: %s %s\n",
			version,
			runtime.Compiler,
			runtime.Version())
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if len(args) < 2 {
		psr.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	ip := args[0]
	domain := args[1]
	spf, err := getSpf(domain)
	if err != nil {
		return checkers.Critical(fmt.Sprintf("DNS lookup failed: %v", err))
	}

	check, err := spf.CheckIP(ip)
	if err != nil {
		return checkers.Critical(fmt.Sprintf("Check SPF failed: %v", err))
	}

	if check != "Pass" {
		return checkers.Critical(fmt.Sprintf("spf check failed: result=%s", check))
	}

	result, err := getReverse(ip)
	if err != nil {
		return checkers.Critical(fmt.Sprintf("reverse lookup dig failed: %v", err))
	}

	if strings.HasSuffix(result, domain) == false {
		return checkers.Critical(fmt.Sprintf("reverse lookup failed: no contains domain - %s", result))
	}

	return checkers.Ok(fmt.Sprintf("OK: spf:%s, reserve-lookup:%s", check, result))
}
