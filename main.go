package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/clbanning/mxj"
)

func main() {
	url := "http://www.dneonline.com/calculator.asmx"
	client := &http.Client{}

	const sRequestContent = `<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<Add xmlns="http://tempuri.org/">
		  <intA>1</intA>
		  <intB>2</intB>
		</Add>
	  </soap:Body>
	</soap:Envelope>`

	requestContent := []byte(sRequestContent)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		log.Fatalf("new request err: %s", err)
	}

	req.Header.Add("SOAPAction", "http://tempuri.org/Add")
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("Accept", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("response err: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("error Respose %s", resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("response err: %s", err)
	}
	xmlMap, _ := mxj.NewMapXml(contents, true)
	fmt.Println(xmlMap)

	addResult, err := xmlMap.ValueForPath("Envelope.Body.AddResponse.AddResult")
	if err != nil {
		log.Fatalf("add result err: %s", err)
	}
	fmt.Println(addResult)
}

// Output :
// map[Envelope:map[-soap:http://schemas.xmlsoap.org/soap/envelope/ -xsi:http://www.w3.org/2001/XMLSchema-instance -xsd:http://www.w3.org/2001/XMLSchema
// Body:map[AddResponse:map[-xmlns:http://tempuri.org/ AddResult:3]]]]
// 3
