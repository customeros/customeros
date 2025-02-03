document.addEventListener("DOMContentLoaded", function () {
  const betaButton = document.getElementById("beta-button");

  if (!betaButton) return;

  betaButton.addEventListener("click", function () {
    const email = "support@customeros.ai";
    const subject = "LinkedIn helper request";
    const mailtoLink = `mailto:${email}?subject=${encodeURIComponent(subject)}`;

    window.open(mailtoLink, "_blank");
  });
});

document.addEventListener("DOMContentLoaded", function () {
  const linkedinButton = document.getElementById("linkedin-button");

  linkedinButton!.addEventListener("click", async () => {
    chrome.tabs.query({ active: true, currentWindow: true }, async (tabs) => {
      if (tabs.length > 0) {
        const url = tabs.find(
          (tab) => tab.url && tab.url.includes("linkedin.com/in")
        )?.url;
        if (!url) {
          console.error("No LinkedIn profile tab found.");
          return;
        }
        // console.log("Current URL:", url);

        const linkedinUrl = url;
        const sessionData = await new Promise<SessionData | null>((resolve) => {
          const onMessage = (message: any) => {
            // console.log("Message received:", message);
            if (message.action === "COS_SESSION_DATA") {
              chrome.runtime.onMessage.removeListener(onMessage);
              resolve({ email: message.email, apiKey: message.apiKey });
            }
          };

          chrome.runtime.onMessage.addListener(onMessage);
        });
        const tenantApiKey = sessionData?.apiKey;

        if (!tenantApiKey) return;

        fetch("https://api.customeros.ai/browserExtension/contact", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Accept: "*/*",
            "X-CUSTOMER-OS-API-KEY": `${tenantApiKey}`,
          },
          body: JSON.stringify({
            linkedinUrl,
          }),
        })
          .then((response) => response.json())
          .then((data) => {
            console.log("Success:", data);
          })
          .catch((error) => {
            console.error("Error:", error);
          });
      }
    });
  });
});
