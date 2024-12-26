package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/miekg/dns"
)


// looktypea function to look up A records (IP addresses)
func looktypea(fqdn, address string) ([]string, error) {
	var m dns.Msg
	var ips []string

	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA) // Ask for A records (IP addresses)
	in, err := dns.Exchange(&m, address)

	if err != nil { // Proper error handling
		return ips, err
	}

	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String()) // Append the IP to the slice
		}
	}

	return ips, nil
}

// Cname function to look up CNAME records
func Cname(fqdn, address string) ([]string, error) {
	var m dns.Msg
	var fqdns []string

	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, address)

	if err != nil {
		return fqdns, err
	}

	if len(in.Answer) < 1 {
		return fqdns, errors.New("no answer")
	}

	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, c.Target) 
		}
	}

	return fqdns, nil
}

// lokkup function to look up CNAME and A records
func lokkup(fqdn, address string) []result {
	var results []result
	var cfqdn = fqdn

	for {
		// Get CNAME records
		cnames, err := Cname(cfqdn, address)
		if err == nil && len(cnames) > 0 {
			cfqdn = cnames[0] 
			continue           
		}

		ips, err := looktypea(cfqdn, address)
		if err != nil {
			break // Exit the loop if there's an error with A record lookup
		}

		// Append the IP addresses to results
		for _, ip := range ips {
			results = append(results, result{IPAddress: ip, Hostname: fqdn})
		}
		break // Exit the loop after processing A records
	}

	return results
}

// Worker function for concurrency
func worker(tracker chan empty, fqdns chan string, gather chan []result, address string) {
	for fqdn := range fqdns {
		results := lokkup(fqdn, address)

		if len(results) > 0 {
			gather <- results
		}
	}

	var e empty
	tracker <- e
}

type empty struct{}

type result struct {
	IPAddress string
	Hostname  string
}

func main() {
	// Usage to be passed when running te program 
	var (
		flDomain     = flag.String("domain", "", "The domain to perform guessing against.")
		flWordlist   = flag.String("wordlist", "", "The wordlist to use for guessing.")
		flWorkerCount = flag.Int("c", 100, "The amount of workers to use.")
		flServerAddr = flag.String("server", "8.8.8.8:53", "The DNS server to use.")
	)
	flag.Parse()

	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain and -wordlist are required")
		os.Exit(1)
	}

	var results []result
	fqdns := make(chan string, *flWorkerCount)
	gather := make(chan []result)
	tracker := make(chan empty)

	fh, err := os.Open(*flWordlist)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for i := 0; i < *flWorkerCount; i++ {
		go worker(tracker, fqdns, gather, *flServerAddr)
	}

	go func() {
		for r := range gather {
			results = append(results, r...)
		}
		var e empty
		tracker <- e
	}()

	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), *flDomain)
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	close(fqdns)

	for i := 0; i < *flWorkerCount; i++ {
		<-tracker
	}
	close(gather)
	<-tracker

	// Print results in a formatted table
	w := tabwriter.NewWriter(os.Stdout, 0, 8, ' ', 0, 0)

	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.IPAddress)
	}
	w.Flush()
}
