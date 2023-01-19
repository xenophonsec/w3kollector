package main

import (
	"fmt"
	"log"
	"net"
)

func lookupDomain(domain string) {
	fmt.Println()
	lookupCNAME(domain)
	fmt.Println()
	getIPs(domain)
	fmt.Println()
	lookupNS(domain)
	fmt.Println()
	lookupMX(domain)
	fmt.Println()
	lookupDNSTXT(domain)
}

func lookupDNSTXT(domain string) {
	records, err := net.LookupTXT(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		fmt.Println("DNS TXT Record:", record)
	}
}

func lookupCNAME(domain string) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CNAME Record:", cname)
}

func getIPs(domain string) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		fmt.Println("IP address:", ip)
	}
}

func lookupNS(domain string) {
	records, err := net.LookupNS(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		fmt.Println("Name Server:", record.Host)
	}
}

func lookupMX(domain string) {
	mxs, err := net.LookupMX(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, mx := range mxs {
		fmt.Println("MX Record:", mx.Host)
	}
}
