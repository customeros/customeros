import OpenGraph from "../components/openGraph";
import { SITE } from "../config";
import { getIsRtlFromLangCode, getLanguageFromURL } from "../languages";
import { getFont } from "../utils/ogFont";
import { SITE_URL } from "../utils/siteUrl";
import { Resvg } from "@resvg/resvg-js";
import { type APIRoute } from "astro";
import satori from "satori";

const removeEndingSlash = (str: string) => str.replace(/\/$/, "");

export const get: APIRoute = async (request) => {
  const params = request.url.searchParams;
  const title = params.get("title") ?? SITE.title;
  const description = params.get("description") ?? SITE.description;
  const pagePath = params.get("pagePath") ?? "";

  // Used for most languages
  const inter = await getFont({
    family: "Inter",
    weights: [400, 700] as const,
  });

  // Used for arabic text
  const bonaNova = await getFont({
    family: "Bona Nova",
    weights: [400, 700] as const,
  });

  // Used for chinese
  const notoSans = await getFont({
    family: "Noto Sans SC",
    weights: [400, 700] as const,
  });

  const hostname = request.site?.hostname.replace(/^https?:\/\//, "");
  const pageLang = getLanguageFromURL(pagePath);

  const svg = await satori(
    OpenGraph({
      title,
      description,
      imageBase: SITE_URL,
      pageUrl: `${hostname}${removeEndingSlash(pagePath)}`,
      rtl: getIsRtlFromLangCode(pageLang),
    }),
    {
      width: 1200,
      height: 630,
      fonts: [
        { name: "Inter", data: inter[400], weight: 400 },
        { name: "Inter", data: inter[700], weight: 700 },
        { name: "Noto Sans SC", data: notoSans[400], weight: 400 },
        { name: "Noto Sans SC", data: notoSans[700], weight: 700 },
        { name: "Bona Nova", data: bonaNova[400], weight: 400 },
        { name: "Bona Nova", data: bonaNova[700], weight: 700 },
      ],
      debug: import.meta.env.DEBUG_OG === "true" ?? false,
    },
  );

  const resvg = new Resvg(svg, {});
  const pngData = resvg.render();
  const pngBuffer = pngData.asPng();

  return new Response(pngBuffer, {
    headers: {
      "Content-Type": "image/png",
      "cache-control": "public, max-age=31536000, immutable",
    },
  });
};
