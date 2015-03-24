package main

import (
  "fmt"
  "net/http"
  "crypto/tls"
  "io/ioutil"
	"strings"
	"strconv"
  "time"
  //"reflect"
	"encoding/json"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/xml"
)

//TODO: 
// Call externally (input opts), mult files, package
// Add concurrency (in parsing and pages?)

// Parse and save profiles
// Output in JSON (MarshalIndent?)

// Add same for job listing

var overall []map[string]string

func main(){
  // Generate search URL
  searchterm := cleanString("golang")
  location := cleanString("")
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
  
  //parseProfile("http://www.indeed.com/r/Nolan-Knight/c56382c1cfeddd01")
  out, _ := json.Marshal(overall)
  fmt.Println(string(out))
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

// Return cleaned up value if input isn't empty
func checkVal(input []xml.Node) string {
  if(len(input) != 0){
    return input[0].InnerHtml()
  } else {
    return ""
  }
}

// Parses and saves data in profile
func parseProfile(url string) {
  body := getPage(url)
  doc, _ := gokogiri.ParseHtml(body)
  
  // Set the values that are the same for all items in profile
  personvals := make(map[string]string)
	
  name, _ := doc.Search("//h1[@itemprop='name']")
  personvals["name"] = checkVal(name)
	
  personvals["url"] = url
  personvals["time_scraped"] = time.Now().Format(time.RFC850)
  //personvals["fulltext"] = string(body)
	
  location, _ := doc.Search("//p[@id='headline_location']")
  personvals["location"] = checkVal(location)
	
  current_title, _ := doc.Search("//h2[@id='headline']")
  personvals["current_title"] = checkVal(current_title)
	
  skills, _ := doc.Search("//div[@id='skills-section']//p")
  personvals["skills"] = checkVal(skills)

  summary, _ := doc.Search("//p[@id='res_summary']")
  personvals["summary"] = checkVal(summary)
	
  additional_info, _ := doc.Search("//div[@id='additionalinfo-section']//p")
  personvals["additional_info"] = checkVal(additional_info)

  // Get all awards person received
  awardslice := make([]map[string]string, 100)
  awards, _ := doc.Search("//div[contains(concat(' ',normalize-space(@class),' '),' award-section ')]")
  for _, award := range(awards){
    awardvals := make(map[string]string)

    award_title, _ := award.Search(".//p[@class='award_title']")
    awardvals["award_title"] = checkVal(award_title)

    award_date, _ := award.Search(".//p[@class='award_date']")
    awardvals["award_date"] = checkVal(award_date)

    award_description, _ := award.Search(".//p[@class='award_description']")
    awardvals["award_description"] = checkVal(award_description)

    awardslice = append(awardslice, awardvals)
    // Push to awardvals
    // Push to overall award hash (part of personvals)
    //fmt.Println(awardvals["award_description"])
  }
  //fmt.Println(awardslice)
  //personvals["awards"] = string(awardslice)

  // TODO
  // Add awards
  // Add certifications
  // Fix date parsing to handle months
  // Fix last vs non-last classes
  // Finish military parser
  // Just append military and education to jobs?


  // Work history
	jobs, _ := doc.Search("//div[@class='work-experience-section ']")
	for _, job := range(jobs) {
    jobvals := make(map[string]string)
    
    job_title, _ := job.Search(".//p[@class='work_title title']")
    jobvals["job_title"] = checkVal(job_title)

    company, _ := job.Search(".//div[@class='work_company']//span")
    jobvals["company"] = checkVal(company)

    company_location, _ := job.Search(".//div[@class='work_company']//div[@class='inline-block']//span")
    jobvals["company_location"] = checkVal(company_location)

    job_description, _ := job.Search(".//p[@class='work_description']")
    jobvals["job_description"] = checkVal(job_description)

    job_dates, _ := job.Search(".//p[@class='work_dates']")
    jobvals["start_date"], jobvals["end_date"] = parseDates(checkVal(job_dates))
    
    addPersonVals(personvals, jobvals)
	}

  // Education info
  degrees, _ := doc.Search("//div[@itemtype='http://schema.org/EducationalOrganization']")
  for _, degree := range(degrees){
    degreevals := make(map[string]string)
    
    school, _ := degree.Search(".//span[@itemprop='name']")
    degreevals["school"] = checkVal(school)

    degree_title, _ := degree.Search(".//p[@class='edu_title']")
    degreevals["degree_title"] = checkVal(degree_title)

    school_location, _ := degree.Search(".//span[@itemprop='addressLocality']")
    degreevals["school_location"] = checkVal(school_location)

    degree_dates, _ := degree.Search(".//p[@class='edu_dates']")
    degreevals["start_date"], degreevals["end_date"] = parseDates(checkVal(degree_dates))

    addPersonVals(personvals, degreevals)
  }
  
  // Military service info
  //milservice, _ := doc.Search("//div[@id='military-items']//div[@class='data_display']")
  //for _, mil := range(milservice){
    // Finish adding
    //fmt.Println(mil)
  //}
}

// Adds the values that are the same for all items
func addPersonVals(personvals map[string]string, itemvals map[string]string){
  for key, val := range(personvals) {
    itemvals[key] = val
  }
  
  overall = append(overall, itemvals)
}

// Handles date parsing
func parseDates(dates string) (string, string){
  split := strings.Split(dates, " ")
  start_date := ""
  end_date := ""
  
  if len(split) >= 3 {
    start_date = split[0]
    end_date = split[len(split)-1]
  } else {
    start_date = dates
    end_date = dates
  }
  return start_date, end_date
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
