package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func lookupDomain(domain string) {
	fmt.Println()
	lookupCNAME(domain)
	lookupHTTP(domain)
	fmt.Println()
	lookupTLS(domain)
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
		fmt.Println("DNS TXT Record: ", record)
	}
}

func lookupCNAME(domain string) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CNAME Record: ", cname)
}

func getIPs(domain string) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		fmt.Println("IP address: ", ip)
	}
}

func lookupNS(domain string) {
	records, err := net.LookupNS(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		fmt.Println("Name Server: ", record.Host)
	}
}

func lookupMX(domain string) {
	mxs, err := net.LookupMX(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, mx := range mxs {
		fmt.Println("MX Record: ", mx.Host)
	}
}

func lookupHTTP(domain string) {
	url := "https://" + domain
	res, err := http.Get(url)
	if err == nil {
		serverHeader := res.Header.Get("server")
		poweredBy := res.Header.Get("X-Powered-By")
		fmt.Println("Protocol:\t" + res.Proto)
		if len(serverHeader) > 0 {
			fmt.Println("Server Engine: ", serverHeader)
		}
		if len(poweredBy) > 0 {
			fmt.Println("Powered By: ", poweredBy)
		}
		fmt.Println()
		for key := range res.Header {
			if (strings.HasPrefix(key, "X-") || strings.HasPrefix(key, "x-")) && key != "X-Powered-By" {
				fmt.Println(key+": ", res.Header.Get(key))
			}
		}
	}
}

func lookupTLS(domain string) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", domain+":443", conf)
	if err != nil {
		log.Println("Error in Dial", err)
		return
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		fmt.Println("Certificate(s)")
	}
	for _, cert := range certs {
		fmt.Printf("Issuer Name: %s\n", cert.Issuer)
		fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
		fmt.Printf("Version: %s \n", strconv.Itoa(cert.Version))
		fmt.Printf("Common Name: %s \n", cert.Issuer.CommonName)

		for _, email := range cert.EmailAddresses {
			fmt.Println("\tAssociated Email Address: " + email)
		}
		for _, name := range cert.DNSNames {
			fmt.Println("   " + name)
		}
		fmt.Println()
	}
}
