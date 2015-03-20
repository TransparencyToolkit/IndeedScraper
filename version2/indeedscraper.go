package main

import (
       "fmt"
       "net/http"
       "crypto/tls"
       "io/ioutil"
       "strings"
       "github.com/moovweb/gokogiri"
)

// Parse and save profiles
// Output in JSON

// Go to next page and repeat
// Also search by location option
// Call externally (input opts), mult files, package
// Add concurrency (in parsing and pages?)
// Add same for job listings

var saved = make(map[string]bool)

func main(){
     // Generate search URL
     searchterm := cleanString("golang")
     url := "http://indeed.com/resumes?q="+searchterm

     // Get page with results
     body := getPage(url)
     
     // Get a list of all profile links on page
     doc, _ := gokogiri.ParseHtml(body)
     results, _ := doc.NodeById("results").Search("//li[@itemtype='http://schema.org/Person']")
     names, _ := results[0].Search("//a[@class='app_link']")

     // Send link of each profile on page to parser
     for _, profile := range(names){
	 parseProfile("http://indeed.com"+profile.Attr("href"))
     }
}

// Parses and saves data in profile
func parseProfile(url string) {
    body := getPage(url)
    doc, _ := gokogiri.ParseHtml(body)
    
    fmt.Println(doc)
}

// Gets the body of a webpage
func getPage(url string) []uint8 {
     // SSL config
     tlsConfig := &tls.Config{
       InsecureSkipVerify: true,
     }
     transport := &http.Transport{
       TLSClientConfig: tlsConfig,
     }
     client := http.Client{Transport: transport}

     // Get page for search term
     resp, _ := client.Get(url)
     defer resp.Body.Close()
     body, _ := ioutil.ReadAll(resp.Body)

     return body
}

// Format search string as needed for URL params
func cleanString(input_term string) string {
 outstr := strings.Replace(input_term, " ", "+", -1)
 outstr = strings.Replace(outstr, ",", "%2C", -1)
 
 return outstr
}
