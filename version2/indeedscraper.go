package main

import (
  "fmt"
  "net/http"
  "crypto/tls"
  "io/ioutil"
	"strings"
	"strconv"
	//"encoding/json"
  "github.com/moovweb/gokogiri"
)

// Parse and save profiles
// Output in JSON (clean up formatting)
// Add timestamp

// Call externally (input opts), mult files, package
// Add concurrency (in parsing and pages?)

// Add same for job listing

func main(){
  // Generate search URL
  searchterm := cleanString("golang")
  location := cleanString("Boston, MA")
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
    //getResults(url + "&start="+strconv.Itoa(i*50))
  }
  
  parseProfile("http://www.indeed.com/r/Adonis-Peralta/805b36f3e9d6f2fa")
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

// Parses and saves data in profile
func parseProfile(url string) {
  body := getPage(url)
  doc, _ := gokogiri.ParseHtml(body)
  
  // Set the values that are the same for all items in profile
  personvals := make(map[string]string)
	
  name, _ := doc.Search("//h1[@itemprop='name']")
  personvals["name"] = name[0].InnerHtml()
	
  personvals["url"] = url
  personvals["fulltext"] = string(body)
	
  location, _ := doc.Search("//p[@id='headline_location']")
  personvals["location"] = location[0].InnerHtml()
	
  current_title, _ := doc.Search("//h2[@id='headline']")
  personvals["current_title"] = current_title[0].InnerHtml()
	
  skills, _ := doc.Search("//div[@id='skills-section']//p")
  personvals["skills"] = skills[0].InnerHtml()

  summary, _ := doc.Search("//p[@id='res_summary']")
  personvals["summary"] = summary[0].InnerHtml()
	
  additional_info, _ := doc.Search("//div[@id='additionalinfo-section']//p")
  personvals["additional_info"] = additional_info[0].InnerHtml()

	//out, _ := json.Marshal(personvals)
	//fmt.Println(string(out))
  // Work history
	//jobs, _ := doc.Search("//div[@class='work-experience-section ']")
	//data, _ := jobs[0].Search("//div[@class='data_display']")
	//for _, job := range(jobs) {       
  //job_title, _ := job.Search("//p[@class='work_title title']")
  company, _ := doc.Search("//div[@class='work_company']")
  fmt.Println(company)
	//for i, _ := range(job_title){
	//reflect.TypeOf(data)
	//fmt.Println(data)
  //}
	//}
	// Company
	// Job title
	// Years (start and end- maybe parsing)
	// Company Location
	// Job Description

  // Education info

  // Military service info
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
