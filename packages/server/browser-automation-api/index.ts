import "dotenv/config";
import { App } from "@/app";

new App();

// (async () => {
//   const scheduler = Scheduler.getInstance();
//   const browser = await Browser.getInstance();
//   const scraper = new Scraper(browser);

//   scheduler.schedule("* * * * *", () => {
//     scraper.scrape("https://www.customeros.ai");
//   });

//   scheduler.startJobs();
// })();
