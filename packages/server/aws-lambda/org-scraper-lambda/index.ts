import { Configuration, OpenAIApi } from "openai";
import { APIGatewayEvent, APIGatewayProxyResult, Context } from "aws-lambda";
import * as https from "https";

const isValidDomain = require("is-valid-domain");
const safeHtml = require("safe-html");
const cheerio = require("cheerio");

const { isUri } = require("valid-url");

export const handler = async (event: APIGatewayEvent, context: Context): Promise<APIGatewayProxyResult> => {
  try {
    // Ensure OPENAI_SECRET_KEY is defined
    if (!process.env.OPENAI_SECRET_KEY) {
      throw new Error("Missing environment variable: OPENAI_SECRET_KEY");
    }

    if (!process.env.X_OPENLINE_TENANT_KEY) {
      throw new Error("Missing environment variable: X_OPENLINE_TENANT_KEY");
    }

    if (!event.body) {
      return {
        statusCode: 400,
        body: JSON.stringify({ error: "Invalid request body" })
      };
    }

    const keys: string[] = process.env.X_OPENLINE_TENANT_KEY.split(" ");
    const containsKey = keys.includes(event.headers["x-openline-tenant-key"]);

    if (!containsKey) {
      return {
        statusCode: 404,
        body: JSON.stringify({ error: "Invalid API Key" }),
      };
    }

    const domain: string = JSON.parse(event?.body).scrapDomain;
    // reject uri's
    if (isUri(domain)) {
      return {
        statusCode: 422,
        body: "Invalid Domain"
      };
    }

    if (isValidDomain(domain)) {
      const response = await fetch(domain);
      var config = {
        allowedTags: ["div", "span", "b", "i", "a"],
        allowedAttributes: {
          "class": {
            allTags: true
          },
          "href": {
            allowedTags: ["a"],
            filter: function(value: any) {
              // Only let through http urls
              if (/^https?:/i.exec(value)) {
                return value;
              }
            }
          }
        }
      };
      if (response.status === 200 || response.body !== undefined) {
        const safeHtmlData = safeHtml(response.body, config);
        const text = extractRelevantText(safeHtmlData);
        const socialLinks = extractSocialLinks(safeHtmlData);
        const analysis = await analyze(domain, text, socialLinks);
        let body = JSON.stringify(analysis);
        return {
          statusCode: 200,
          body: body
        };
      } else {
        return {
          statusCode: 422,
          body: "Unable to retrieve information about domain"
        };
      }
    } else {
      return {
        statusCode: 422,
        body: "Invalid domain"
      };
    }
  } catch (error) {
    console.error("Error:", error);

    return {
      statusCode: 500,
      body: JSON.stringify({ error: "Internal Server Error", message: error.message })
    };
  }
};

async function analyze(website: string, text: string, socials: { [key: string]: string }): Promise<{
  analysis: string
}> {
  // OpenAI configuration creation
  const configuration = new Configuration({
    apiKey: process.env.OPENAI_SECRET_KEY
  });
  const openaiClient = new OpenAIApi(configuration);

  // OpenAI instance creation
  let c_prompt = COMPANY_PROMPT.replace("{{text}}", text);
  const analysis: string = await openaiEnhance(openaiClient, c_prompt);
  // Define the replacements as an object

  const replacements: Record<string, string> = {
    "{{ANALYSIS}}": analysis,
    "{{WEBSITE}}": website,
    "{{SOCIALS}}": JSON.stringify(socials)
  };
  // Perform variable replacements in the file data using a callback function

  const s_prompt: string = SCRAPED_PROMPT.replace(
    /{{ANALYSIS}}|{{WEBSITE}}|{{SOCIALS}}/g,
    (match) => replacements[match]
  );
  const cleanText = await openaiEnhance(openaiClient, s_prompt);

  // Try parsing the cleaned analysis as JSON
  let parsedAnalysis;
  try {
    parsedAnalysis = JSON.parse(cleanText);
  } catch (err) {
    console.error("Error parsing analysis as JSON:", err);
    throw err;
  }

  return parsedAnalysis;
}

async function openaiEnhance(openaiClient: OpenAIApi, prompt: string) {

  const completion = await openaiClient.createCompletion({
    model: "text-davinci-003",
    prompt,
    max_tokens: 400
  });

  // Ensure choices exist and contain at least one item
  if (!completion.data.choices || completion.data.choices.length === 0) {
    throw new Error("No data returned from OpenAI API");
  }

  // Ensure text exists
  const cleanText = completion.data.choices[0].text;
  if (!cleanText) {
    throw new Error("No text returned from OpenAI API");
  }
  return cleanText.trim();
}

function extractSocialLinks(html: string): { [key: string]: string } {
  const $ = cheerio.load(html);

  // Define a mapping from social media site names to URL patterns
  const socialMediaSites = {
    linkedin: "linkedin.com",
    twitter: "twitter.com",
    instagram: "instagram.com",
    facebook: "facebook.com",
    youtube: "youtube.com",
    github: "github.com"
  };

  const socialLinks: { [key: string]: string } = {};

  // Search for all links in the footer
  const links = $("a");

  links.each((_: any, element: any) => {
    const link = $(element).attr("href");
    if (link) {
      for (const site in socialMediaSites) {
        // Type assertion here
        if (link.includes(socialMediaSites[site as keyof typeof socialMediaSites])) {
          socialLinks[site] = link;
          break;
        }
      }
    }
  });

  return socialLinks;
}

function extractRelevantText(html: string): string {
  const $ = cheerio.load(html);

  // Remove script and style tags
  $("script, style").remove();

  // Extract text from all leaf nodes (elements with no child elements)
  const leafNodes = $("*:not(:has(*))");
  const texts: string[] = [];
  leafNodes.each((_: any, element: any) => {
    const text = $(element).text().trim();
    if (text.length > 0 && !texts.includes(text)) {
      texts.push(text);
    }
  });

  return texts.join(" ");
}

const COMPANY_PROMPT = "  Analyze the following text from a company website.\n" +
  "  \n" +
  "  {{text}}\n" +
  "  \n" +
  "  Analyze the text and respond (in English) as defined below:\n" +
  "\n" +
  "    {\n" +
  "    companyName:  the name of the company\n" +
  "    market: options are B2B, B2C, or Marketplace\n" +
  "    industry: Industry per the Global Industry Classification Standard (GISB),\n" +
  "    industryGroup: Industry Group per the Global Industry Classification Standard (GISB),\n" +
  "    subIndustry: Sub-industry per the Global Industry Classification Standard (GISB),\n" +
  "    targetAudience: analysis of the company's target audience,\n" +
  "    valueProposition: analysis of the company's core value proposition,\n" +
  "    }";

const SCRAPED_PROMPT = "  The following is data scraped from a website:  Please combine and format the data into a clean json response\n" +
  "\n" +
  "  {{ANALYSIS}}\n" +
  "\n" +
  "  website: {{WEBSITE}}\n" +
  "\n" +
  "  {{SOCIALS}}\n" +
  "\n" +
  "  --------\n" +
  "\n" +
  "  Put the data above in the following JSON structure\n" +
  "\n" +
  "  {\n" +
  "    \"companyName\": \"..\",\n" +
  "    \"website\": \"..\",\n" +
  "    \"market\": \"..\",\n" +
  "    \"industry\": \"..\",\n" +
  "    \"industryGroup\": \"..\",\n" +
  "    \"subIndustry\": \"..\",\n" +
  "    \"targetAudience\": \"..\",\n" +
  "    \"valueProposition\": \"..\",\n" +
  "    \"linkedin\": \"..\",\n" +
  "    \"twitter\": \"..\",\n" +
  "    \"instagram\": \"..\",\n" +
  "    \"facebook\": \"..\",\n" +
  "    \"youtube\": \"..\",\n" +
  "    \"github\": \"..\",\n" +
  "  }\n" +
  "\n" +
  "  If you do not have data for a given key, do not return it as part of the JSON object.\n" +
  "\n" +
  "  Ensure before you output that your response is valid JSON.  If it is not valid JSON, do your best to fix the formatting to align to valid JSON.\n" +
  "\n";