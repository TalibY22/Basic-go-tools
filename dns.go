package main


import (
	

	 "github.com/miekg/dns"
	
) 


func main()  {
	
  var msg dns.Msg
  
  fqdn := dns.Fqdn("github.com")//Change the domain name into fully qualified domain name 
  msg.SetQuestion(fqdn,dns.TypeA)//Indicates the ip address of a domain 
  dns.Exchange(&msg, "8.8.8.8:53")//The 8.8.8.8;53 is googles dns address 




}

