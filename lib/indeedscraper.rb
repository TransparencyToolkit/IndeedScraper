require 'open-uri'
require 'json'
require 'nokogiri'
require 'date'

class IndeedScraper
  def initialize(searchterm, location)
    @searchterm = searchterm
    @location = location
    @output = Array.new
  end

  # Get all results
  def searchResumes
    @searchterm.gsub!(" ", "-")
    if @location != nil
      @location.gsub!(", ", "-")
      @location.gsub!(" ", "-")
      url = "http://www.indeed.com/resumes/" + @searchterm + "/in-" + @location
    else
      url = "http://www.indeed.com/resumes/" + @searchterm
    end
    html = Nokogiri::HTML(open(url))
    
    # Handle multiple pages
    numresults = html.css("div#result_count").text.split(" ")
    fresult = numresults[0].to_i/50.0
    if fresult != numresults[0].to_i/50
      count = fresult +1
    else
      count = numresults[0].to_i/50
    end
    
    # Loop through pages and get results
    i = 1
    while i <= count
      results = html.css("ol#results")
      results.css("li").each do |l|
        getResume("http://indeed.com"+l.css("a")[0]["href"].gsub("?sp=0",""))
      end
      i += 1
      nextstart = (i-1)*50
      html = Nokogiri::HTML(open(url+"?start="+nextstart.to_s))
    end
  end

  # Process and save resume data
  def getResume(url)
    page = Nokogiri::HTML(open(url))
    name = page.css('h1[itemprop="name"]').text
    location = page.css('p.locality').text
    currtitle = page.css('h2[itemprop="jobTitle"]').text
    summary = page.css('p#res_summary').text
    additionalinfo = page.css('div#additionalinfo-section').text
    skills = page.css("div#skills-section").text

    # Get work info
    page.css("div.work-experience-section").each do |w|
      positionhash = Hash.new
      positionhash[:name] = name
      positionhash[:url] = url
      positionhash[:title] = w.css("p.work_title").text
      if w.css("div.work_company").css("span")[0]
        positionhash[:company] = w.css("div.work_company").css("span")[0].text
      end
      if w.css("div.work_company").css("span")[1]
        positionhash[:company_location] = w.css("div.work_company").css("span")[1].text
      end

      # Process date info
      dates = dateParse(w.css("p.work_dates"))
      if dates
        positionhash[:start_date] = dates[0]
        positionhash[:end_date] = dates[1]
      end

      positionhash[:description] =  w.css("p.work_description").text

      # Info for all positions
      positionhash[:skills] = skills
      positionhash[:current_location] = location
      positionhash[:current_title] = currtitle
      positionhash[:summary] = summary
      positionhash[:additional_info] = additionalinfo
      @output.push(positionhash)
    end

    # Get education info
    page.css("div.education-section").each do |e|
      eduhash = Hash.new
      eduhash[:name] = name
      eduhash[:url] = url
      eduhash[:degree] = e.css("p.edu_title").text
      eduhash[:school] = e.css('span[itemprop="name"]').text
      eduhash[:dates] = e.css("p.edu_dates").text

      # Info for all degrees
      eduhash[:skills] = skills
      eduhash[:current_location] = location
      eduhash[:current_title] = currtitle
      eduhash[:summary] = summary
      eduhash[:additional_info] = additionalinfo
      @output.push(eduhash)
    end

    # Get military service info
    page.css("div.military-section").each do |m|
      milhash = Hash.new
      milhash[:name] = name
      milhash[:url] = url
      milhash[:service_country] = m.css("p.military_country").text.gsub("Service Country: ", "") 
      milhash[:branch] = m.css("p.military_branch").text.gsub("Branch: ", "")
      milhash[:rank] = m.css("p.military_rank").text.gsub("Rank: ", "")
    
      # Parse dates                                                                                                     
      dates = dateParse(m.css("p.military_date"))
      milhash[:start_date] = dates[0]
      milhash[:end_date] = dates[1]

      milhash[:military_description] = m.css("p.military_description").text
      milhash[:military_commendations] = m.css("p.military_commendations").text.split("\n")

      # Info for all items
      milhash[:skills] = skills
      milhash[:current_location] = location
      milhash[:current_title] = currtitle
      milhash[:summary] = summary
      milhash[:additional_info] = additionalinfo
      @output.push(milhash)
    end
  end

  # Process dates
  def dateParse(date)
    datearray = Array.new
    daterange = date.text.split(" to ")
    if daterange[0] != nil
      datearray[0] = DateTime.parse(dateCheck(daterange[0]))
    else
      datearray[0] = nil
    end
    
    if daterange[1] == "Present"
      datearray[1] = "Present"
    else
      if daterange[1] != nil
        datearray[1] = DateTime.parse(dateCheck(daterange[1]))
      else
        datearray = nil
      end
    end

    return datearray
  end

  # Handle year only dates
  def dateCheck(date)
    if date.length == 4
      return "January " + date
    else
      return date
    end
  end

  def searchJobs
    # Save company/job name
    # Save location
    # Save url
    # Download url
    # Get page text
    # Add each to the output
  end

  # Generates JSON output
  def getOutput
    JSON.pretty_generate(@output)
  end
end

