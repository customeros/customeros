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

let sessionData: SessionData;

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

    console.log("CustomerOS tab found:", customerOSTab);

    await new Promise<void>((resolve, reject) => {
      console.log("Attempting to inject content script...");
      chrome.scripting.executeScript(
        {
          target: { tabId: customerOSTab.id as any },
          files: ["contentScript.js"],
        },
        () => {
          if (chrome.runtime.lastError) {
            console.error(
              "Error injecting content script:",
              chrome.runtime.lastError.message
            );
            reject(
              new Error(
                `Error injecting content script: ${chrome.runtime.lastError.message}`
              )
            );
          } else {
            console.log("Content script injected successfully");
            resolve();
          }
        }
      );
    });

    // Retrieve session data
    const sessionData = await new Promise<SessionData | null>((resolve) => {
      const onMessage = (message: any) => {
        console.log("Message received:", message);
        if (message.action === "COS_SESSION_DATA") {
          chrome.runtime.onMessage.removeListener(onMessage);
          resolve({ email: message.email, apiKey: message.apiKey });
        }
      };

      chrome.runtime.onMessage.addListener(onMessage);
    });

    console.log("Session Data:", sessionData);
    if (!sessionData) {
      console.error("Error: sessionData is null");
      return;
    }

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

    console.log("LinkedIn Cookies:", cookies);

    if (!cookies) return;

    const liAtCookie = cookies?.find((cookie) => cookie.name === "li_at");
    console.log("li_at Cookie:", liAtCookie);

    if (!liAtCookie) return;

    const userAgent = navigator.userAgent || "unknown";

    chrome.storage.local.get(["storedCookies"], async (result) => {
      const storedCookies = result.storedCookies;
      let previousCookies = null;

      if (storedCookies) {
        previousCookies = storedCookies;
      }

      if (previousCookies && previousCookies.value === liAtCookie.value) {
        console.log("Cookie is the same, no update needed");
        return;
      }

      console.log("Sending new cookie data...");

      const response = await fetch(`${BASE_URL}/browser/config`, {
        headers: {
          "x-openline-api-key": `${sessionData.apiKey}`,
          "x-openline-username": `${sessionData.email}`,
        },
        method: "GET",
      });
      const data: data = await response.json();

      if (data?.data?.cookies && data?.data?.cookies.length > 0) {
        const prevLiAtCookie = data?.data?.cookies;

        if (prevLiAtCookie !== liAtCookie.value) {
          console.log("Different cookie detected, updating...");

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

async function handleExtensionButtonClick(tab: chrome.tabs.Tab) {
  chrome.windows.getCurrent((window) => {
    chrome.sidePanel.open({
      tabId: Number(tab.id),
      windowId: window?.id,
    });
  });
  await chrome.sidePanel.setOptions({
    tabId: tab.id,
    path: "sidepanel/index.html",
    enabled: true,
  });
}

chrome.action.onClicked.addListener((tab: chrome.tabs.Tab) =>
  handleExtensionButtonClick(tab)
);

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status === "complete" && tab.url) {
    const url = new URL(tab.url);
    if (url.host === "www.linkedin.com" || url.host === "linkedin.com") {
      chrome.sidePanel.setOptions({
        tabId: tabId,
        path: "sidepanel/index.html",
        enabled: true,
      });
    } else {
      chrome.sidePanel.setOptions({
        tabId: tabId,
        enabled: false,
      });
    }
  }
});

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === "openTab") {
    chrome.tabs.create({ url: message.url, active: true });
    sendResponse({ success: true });
  }
});

chrome.tabs.onActivated.addListener(async (activeInfo) => {
  const tab = await chrome.tabs.get(activeInfo.tabId);
  if (tab.url) {
    const url = new URL(tab.url);
    if (url.host === "www.linkedin.com" || url.host === "linkedin.com") {
      chrome.sidePanel.setOptions({
        tabId: activeInfo.tabId,
        path: "sidepanel/index.html",
        enabled: true,
      });
    } else {
      chrome.sidePanel.setOptions({
        tabId: activeInfo.tabId,
        enabled: false,
      });
    }
  }
});
