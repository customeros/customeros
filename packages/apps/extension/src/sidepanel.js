const BASE_URL = "https://bas.customeros.ai";
const DEFAULT_LOGO = "/src/assets/customeros.png";

async function getSessionData(retries = 3) {
    return new Promise((resolve) => {
        function attemptGetSessionData() {
            chrome.runtime.sendMessage({ action: "GET_SESSION_DATA" }, (response) => {
                console.log("Received session data in sidepanel:", response);
                if (response && response.email && response.apiKey) {
                    resolve(response);
                } else if (retries > 0) {
                    console.log(`Retrying to get session data. Attempts left: ${retries}`);
                    setTimeout(() => attemptGetSessionData(), 1000);
                    retries--;
                } else {
                    console.error("Failed to get valid session data after retries");
                    resolve(null);
                }
            });
        }
        attemptGetSessionData();
    });
}

async function fetchWorkspaceInfo() {
    try {
        const sessionData = await getSessionData();

        if (!sessionData || !sessionData.email || !sessionData.apiKey) {
            console.error("Session data is incomplete:", sessionData);
            throw new Error("Session data is incomplete or missing");
        }

        const response = await fetch(`${BASE_URL}/browser/config`, {
            headers: {
                "x-openline-api-key": sessionData.apiKey,
                "x-openline-username": sessionData.email,
            },
            method: "GET",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch workspace info");
        }

        const data = await response.json();
        return {
            workspaceName: data.data.workspaceName || "CustomerOS",
            workspaceLogo: data.data.workspaceLogo || DEFAULT_LOGO,
        };
    } catch (error) {
        console.error("Error fetching workspace info:", error);
        return {
            workspaceName: "CustomerOS",
            workspaceLogo: DEFAULT_LOGO,
        };
    }
}

async function updateWorkspaceInfo() {
    const workspaceLogoElement = document.getElementById("workspace-logo");
    const workspaceNameElement = document.getElementById("workspace-name");

    try {
        const { workspaceName, workspaceLogo } = await fetchWorkspaceInfo();
        workspaceLogoElement.src = workspaceLogo;
        workspaceNameElement.textContent = workspaceName;
    } catch (error) {
        console.error("Error updating workspace info:", error);
        workspaceLogoElement.src = DEFAULT_LOGO;
        workspaceNameElement.textContent = "CustomerOS";
    }
}

document.addEventListener("DOMContentLoaded", updateWorkspaceInfo);
