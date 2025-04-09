# Website content analyzer

- once a website is enrolled in eventstream, (tracker.{domain}.provisioned), we:
- crawl website, scrape content, throw event
  scraper.page.no_content
  scraper.page.scraped

- classify webpage (what problem does it solve? for whom?)
  scraper.page.category.classified
  scraper.page.topics.identified
  scraper.page.content_stage.determined

## Events

scraper.error
scraper.page.no_content
scraper.page.scraped
scraper.page.category.classified
scraper.page.topics.identified
scraper.page.content_stage.determined
