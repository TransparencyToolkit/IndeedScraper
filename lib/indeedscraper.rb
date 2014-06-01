require 'open-uri'
require 'json'
require 'nokogiri'
require 'date'

class IndeedScraper
  def initialize(searchterm)
    @searchterm = searchterm
    @output = Array.new
    # Add location specification
    # Specify if company or resume
  end

  # Get all results
  def search
    url = "http://www.indeed.com/resumes/" + @searchterm
    html = Nokogiri::HTML(open(url))
    results = html.css("ol#results")
    results.css("li").each do |l|
      getResume("http://indeed.com"+l.css("a")[0]["href"].gsub("?sp=0",""))
    end
    
    # Handle multiple pages
    # Add locations
    # Add jobs
  end

  # Process and save resume data
  def getResume(url)
    #begin
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
      positionhash[:company] = w.css("div.work_company").css("span")[0].text
      if w.css("div.work_company").css("span")[1]
        positionhash[:company_location] = w.css("div.work_company").css("span")[1].text
      end

      # Process date info
      daterange = w.css("p.work_dates").text.split(" to ")
      positionhash[:start_date] = DateTime.parse(dateCheck(daterange[0]))
      if daterange[1] == "Present"
        positionhash[:end_date] = "Present"
      else
        positionhash[:end_date] = DateTime.parse(dateCheck(daterange[1]))
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
      puts m
    end
      # Add military service section
      
    #rescue
     # puts url
    #end
  end

  # Handle year only dates
  def dateCheck(date)
    if date.length == 4
      return "January " + date
    else
      return date
    end
  end

  def getlisting
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


