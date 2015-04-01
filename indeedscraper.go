package indeedscraper

import (
  "net/http"
  "crypto/tls"
  "io/ioutil"
	"strings"
	"strconv"
	"encoding/json"
  "github.com/moovweb/gokogiri"
)

//TODO:
// Add concurrency (in parsing and pages?)
// Add same for job listing
// Clean up result names/types

var overall []map[string]string

// Download all the resumes
func GetResumes(searchterm string, location string) string {
  // Generate search URL
  searchterm = cleanString("golang")
  location = cleanString("")
  url := "http://indeed.com/resumes?"

  // Add search term to URL
  if searchterm != "" {
    url += "q="+searchterm
  }

  // Add location to URL
  if location != "" {
    if strings.Contains(url, "?q="){
      url += "&"
    }
    url += "l="+location
  }


  // Get page with results
  body := getPage(url)
  numPages := getPageCount(body)

  // Loop through all pages
  for i := 0; i < numPages; i++ {
    getResults(url + "&start="+strconv.Itoa(i*50))
  }

  out, _ := json.MarshalIndent(overall, "", "    ")
  return string(out)
}

// Gets the results for a single page
func getResults(resultsurl string) {
  body := getPage(resultsurl)
  
	// Get a list of all profile links on page
	doc, _ := gokogiri.ParseHtml(body)
	results, _ := doc.NodeById("results").Search("//li[@itemtype='http://schema.org/Person']")
	names, _ := results[0].Search("//a[@class='app_link']")

	// Send link of each profile on page to parser
	for _, profile := range(names){
    parseProfile("http://indeed.com"+profile.Attr("href"))
	}
}

// Gets the total number of result pages
func getPageCount(firstpage []uint8) int {
	parsed, _ := gokogiri.ParseHtml(firstpage)
	numresults, _ := parsed.Search("//div[@id='result_count']")
	num, _ := strconv.Atoi(strings.Split(numresults[0].InnerHtml(), " ")[1])
  
	numpages := num/50
	if num % 50 != 0 {
	  numpages += 1
	}

	return numpages
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
