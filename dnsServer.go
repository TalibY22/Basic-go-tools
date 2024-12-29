package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

func main() {
	// Define a DNS handler


	fmt.Println("am alive ")
	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		var response dns.Msg
		response.SetReply(req)

		// Loop through the questions in the DNS query
		for _, q := range req.Question { // `req.Question` instead of `req.question`
			// Create a DNS A record
			a := &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP("127.0.0.1").To4(),
			}

			// Add the answer to the response
			response.Answer = append(response.Answer, a)
		}

		// Write the response message
		w.WriteMsg(&response)
	})

	// Start the DNS server
	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
