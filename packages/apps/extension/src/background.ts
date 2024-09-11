const BASE_URL = "https://bas.customeros.ai";

type SessionData = {
  email: string | null;
  apiKey: string | null;
};

type data = {
  succes: boolean;
  data: {
    id: number;
    tenant: string | null;
    userId: string | null;
    cookies: string | null;
    createAt: string | null;
    updateAt: string | null;
    userAgent: string | null;
  };
};

async function getCookiesFromLinkedInTab() {
  try {
    const customerOSTab = await new Promise<chrome.tabs.Tab>(
      (resolve, reject) => {
        chrome.tabs.query({}, (tabs) => {
          const foundTab = tabs.find(
            (tab) =>
              tab.url &&
              (tab.url.includes("localhost") || tab.url.includes("customeros"))
          );
          if (foundTab) {
            resolve(foundTab);
          } else {
            reject(new Error("No app.customeros.ai tab found"));
          }
        });
      }
    );

    await new Promise<void>((resolve, reject) => {
      chrome.scripting.executeScript(
        {
          target: { tabId: customerOSTab.id as any },
          files: ["contentScript.js"],
        },
        () => {
          if (chrome.runtime.lastError) {
            reject(
              new Error(
                `Error injecting content script: ${chrome.runtime.lastError.message}`
              )
            );
          } else {
            resolve();
          }
        }
      );
    });

    const sessionData = await new Promise<SessionData | null>((resolve) => {
      const onMessage = (message: any) => {
        if (message.action === "COS_SESSION_DATA") {
          chrome.runtime.onMessage.removeListener(onMessage);
          resolve({ email: message.email, apiKey: message.apiKey });
        }
      };

      chrome.runtime.onMessage.addListener(onMessage);
    });

    const userNameCustomerOs = sessionData?.email;
    const apiKey = sessionData?.apiKey;

    if (!userNameCustomerOs || !apiKey) {
      console.error(
        "Error: Could not retrieve userNameCustomerOs or apiKey from session"
      );
      return;
    }

    const cookies = await new Promise<chrome.cookies.Cookie[]>(
      (resolve, reject) => {
        chrome.cookies.getAll({ domain: ".linkedin.com" }, (result) => {
          if (chrome.runtime.lastError) {
            reject(
              new Error(
                `Error retrieving cookies: ${chrome.runtime.lastError.message}`
              )
            );
          } else {
            resolve(result);
          }
        });
      }
    );

    const liAtCookie = cookies.find((cookie) => cookie.name === "li_at");

    if (!liAtCookie) return;

    const userAgent = navigator.userAgent || "unknown";

    chrome.storage.local.get(["storedCookies"], async (result) => {
      const storedCookies = result.storedCookies;
      let previousCookies = null;

      if (storedCookies) {
        previousCookies = storedCookies;
      }

      if (previousCookies && previousCookies.value === liAtCookie.value) {
        return;
      }

      const response = await fetch(`${BASE_URL}/browser/config`, {
        headers: {
          "x-openline-api-key": `${sessionData.apiKey}`,
          "x-openline-username": `${sessionData.email}`,
        },
        method: "GET",
      });
      const data: data = await response.json();

      if (data.data.cookies && data.data.cookies.length > 0) {
        const prevLiAtCookie = data.data.cookies;

        if (prevLiAtCookie !== liAtCookie.value) {
          console.log("Different cookie detected");

          await fetch(`${BASE_URL}/browser/config`, {
            method: "PATCH",
            headers: {
              "Content-Type": "application/json",
              "x-openline-api-key": `${sessionData.apiKey}`,
              "x-openline-username": `${sessionData.email}`,
            },
            body: JSON.stringify({
              cookies: `[{\"name\":\"li_at\",\"value\":\"${liAtCookie?.value}\",\"domain\":\"www.linkedin.com\",\"path\":\"/\",\"secure\":true,\"httpOnly\":true,\"sameSite\":\"Lax\"}]`,
              userAgent: userAgent,
            }),
          });
        } else {
          console.log("Cookie is the same");
        }
      } else {
        console.log("No previous cookies, creating new record");

        await fetch(`${BASE_URL}/browser/config`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "x-openline-api-key": `${sessionData.apiKey}`,
            "x-openline-username": `${sessionData.email}`,
          },
          body: JSON.stringify({
            cookies: `[{\"name\":\"li_at\",\"value\":\"${liAtCookie?.value}\",\"domain\":\"www.linkedin.com\",\"path\":\"/\",\"secure\":true,\"httpOnly\":true,\"sameSite\":\"Lax\"}]`,
            userAgent: userAgent,
          }),
        });
      }

      chrome.storage.local.set({ storedCookies: liAtCookie });
    });
  } catch (error) {
    console.error("Error in getCookiesFromLinkedInTab:", error);
  }
}

chrome.alarms.create("checkLinkedInCookies", { periodInMinutes: 0.1 });

chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === "checkLinkedInCookies") {
    getCookiesFromLinkedInTab();
  }
});
