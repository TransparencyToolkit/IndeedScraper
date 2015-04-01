package indeedscraper

import(
  "time"
  "encoding/json"
  "github.com/moovweb/gokogiri"
  "strings"
  "github.com/moovweb/gokogiri/xml"
)

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
  personvals["fulltext"] = string(body)

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

  links, _ := doc.Search("//div[contains(concat(' ',normalize-space(@class),' '),' link-section ')]")
  linkslice := make([]map[string]string, len(links))
  for n, link := range(links){
    linkvals := make(map[string]string)

    link_title, _ := link.Search(".//a")
    linkvals["link_title"] = checkVal(link_title)

    link_url, _ := link.Search(".//a")
    linkvals["link_url"] = link_url[0].Attr("href")

    linkslice[n] = linkvals
  }
  linklist, _ := json.Marshal(linkslice)
  personvals["links"] = string(linklist)

  // Get all awards person received
  awards, _ := doc.Search("//div[contains(concat(' ',normalize-space(@class),' '),' award-section ')]")
  awardslice := make([]map[string]string, len(awards))
  for n, award := range(awards){
    awardvals := make(map[string]string)

    award_title, _ := award.Search(".//p[@class='award_title']")
    awardvals["award_title"] = checkVal(award_title)

    award_date, _ := award.Search(".//p[@class='award_date']")
    awardvals["award_date"] = checkVal(award_date)

    award_description, _ := award.Search(".//p[@class='award_description']")
    awardvals["award_description"] = checkVal(award_description)

    awardslice[n] = awardvals
  }
  awardlist, _ := json.Marshal(awardslice)
  personvals["awards"] = string(awardlist)

  // Get all groups
  groups, _ := doc.Search("//div[contains(concat(' ',normalize-space(@class),' '),' group-section ')]")
  groupslice := make([]map[string]string, len(groups))
  for n, group := range(groups){
    groupvals := make(map[string]string)

    group_title, _ := group.Search(".//p[@class='group_title']")
    groupvals["group_title"] = checkVal(group_title)

    group_date, _ := group.Search(".//p[@class='group_date']")
    groupvals["group_start_date"], groupvals["group_end_date"] = parseDates(checkVal(group_date))

    group_description, _ := group.Search(".//p[@class='group_description']")
    groupvals["group_description"] = checkVal(group_description)

    groupslice[n] = groupvals
  }
  grouplist, _ := json.Marshal(groupslice)
  personvals["groups"] = string(grouplist)

  // Get all certifications
  certifications, _ := doc.Search("//div[contains(concat(' ',normalize-space(@class),' '),' certification-section ')]")
  certificationslice := make([]map[string]string, len(certifications))
  for n, certification := range(certifications){
    certvals := make(map[string]string)

    cert_title, _ := certification.Search(".//p[@class='certification_title']")
    certvals["cert_title"] = checkVal(cert_title)

    cert_date, _ := certification.Search(".//p[@class='certification_date']")
    certvals["cert_start_date"], certvals["cert_end_date"] = parseDates(checkVal(cert_date))

    cert_description, _ := certification.Search(".//p[@class='certification_description']")
    certvals["cert_description"] = checkVal(cert_description)

    certificationslice[n] = certvals
  }
  certlist, _ := json.Marshal(certificationslice)
  personvals["certifications"] = string(certlist)

  // Military service info
  milservice, _ := doc.Search("//div[contains(concat(' ',normalize-space(@class),' '),' military-section ')]")
  milslice := make([]map[string]string, len(milservice))
  for n, mil := range(milservice){
    milvals := make(map[string]string)

    military_country, _ := mil.Search(".//p[@class='military_country']")
    milvals["military_country"] = strings.Split(checkVal(military_country), "</span>")[1]

    military_branch, _ := mil.Search(".//p[@class='military_branch']")
    milvals["military_branch"] = strings.Split(checkVal(military_branch), "</span>")[1]

    military_rank, _ := mil.Search(".//p[@class='military_rank']")
    milvals["military_rank"] = strings.Split(checkVal(military_rank), "</span>")[1]

    military_description, _ := mil.Search(".//p[@class='military_description']")
    milvals["military_description"] = checkVal(military_description)

    military_commendations, _ := mil.Search(".//p[@class='military_commendations']")
    milvals["military_commendations"] = checkVal(military_commendations)

    military_date, _ := mil.Search(".//p[@class='military_date']")
    milvals["military_start_date"], milvals["military_end_date"] = parseDates(checkVal(military_date))

    milslice[n] = milvals
  }
  millist, _ := json.Marshal(milslice)
  personvals["military_service"] = string(millist)

  // Education info
  degrees, _ := doc.Search("//div[@itemtype='http://schema.org/EducationalOrganization']")
  degreeslice := make([]map[string]string, len(degrees))
  for n, degree := range(degrees){
    degreevals := make(map[string]string)

    school, _ := degree.Search(".//span[@itemprop='name']")
    degreevals["school"] = checkVal(school)

    degree_title, _ := degree.Search(".//p[@class='edu_title']")
    degreevals["degree_title"] = checkVal(degree_title)

    school_location, _ := degree.Search(".//span[@itemprop='addressLocality']")
    degreevals["school_location"] = checkVal(school_location)

    degree_dates, _ := degree.Search(".//p[@class='edu_dates']")
    degreevals["degree_start_date"], degreevals["degree_end_date"] = parseDates(checkVal(degree_dates))

    degreeslice[n] = degreevals
  }
  degreelist, _ := json.Marshal(degreeslice)
  personvals["degrees"] = string(degreelist)

  // Work history
  jobs, _ := doc.Search("//div[contains(concat(' ',normalize-space(@class),' '),' work-experience-section ')]")
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
}

// Adds the values that are the same for all items
func addPersonVals(personvals map[string]string, itemvals map[string]string){
  for key, val := range(personvals) {
    itemvals[key] = val
  }

  overall = append(overall, itemvals)
}

// Return cleaned up value if input isn't empty
func checkVal(input []xml.Node) string {
  if(len(input) != 0){
    return input[0].InnerHtml()
  } else {
    return ""
  }
}

// Handles date parsing
func parseDates(dates string) (string, string){
  var start_date string
  var end_date string

  // Deal with different date formats
  if strings.Contains(dates, "bis") {
    split := strings.Split(dates, " bis ")
    start_date = split[0]
    end_date = split[1]
  } else if strings.Contains(dates, " to "){
    split := strings.Split(dates, " to ")
    start_date = split[0]
    end_date = split[1]
  } else {
    start_date = dates
    end_date = dates
  }

  return start_date, end_date
}
