import { useEffect, useState } from "react";
import "./styles/tailwind.css";
import { Button } from "@ui/form/Button/Button";

export const App = () => {
  const [workspaceName, setWorkspaceName] = useState<string | null>(null);
  const [appIsOpen, setAppIsOpen] = useState(false);
  const [linkedInUrl, setLinkedInUrl] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const handleMessage = (message: any) => {
      if (message.action === "COS_SESSION_DATA") {
        setWorkspaceName(message.workspaceName || null);
      }
    };

    chrome.runtime.onMessage.addListener(handleMessage);

    return () => {
      chrome.runtime.onMessage.removeListener(handleMessage);
    };
  }, []);

  useEffect(() => {
    const handleTabUpdate = (
      _tabId: number,
      changeInfo: chrome.tabs.TabChangeInfo,
      tab: chrome.tabs.Tab
    ) => {
      if (changeInfo.url && tab.active && tab.url?.length) {
        setLinkedInUrl(tab.url.includes("linkedin.com/in") ? tab.url : null);
        chrome.runtime.sendMessage({ action: "REQUEST_COS_SESSION_DATA" });
      }
    };

    chrome.tabs.onUpdated.addListener(handleTabUpdate);

    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      const url = tabs.find((tab) => tab.url?.includes("linkedin.com/in"))?.url;
      setLinkedInUrl(url || null);
    });

    chrome.runtime.sendMessage({ action: "REQUEST_COS_SESSION_DATA" });

    return () => {
      chrome.tabs.onUpdated.removeListener(handleTabUpdate);
    };
  }, []);

  useEffect(() => {
    const checkCustomerOSTab = () => {
      chrome.tabs.query({}, (tabs) => {
        const customerOSTab = tabs.find(
          (tab) =>
            tab.url?.includes("app.customeros.ai") ||
            tab.url?.includes("localhost:5173")
        );
        if (customerOSTab) {
          setAppIsOpen(true);
        } else {
          setAppIsOpen(false);
        }
      });
    };

    checkCustomerOSTab();

    const handleTabUpdate = (
      _tabId: number,
      changeInfo: chrome.tabs.TabChangeInfo,
      tab: chrome.tabs.Tab
    ) => {
      if (changeInfo.url || tab.status === "complete") {
        checkCustomerOSTab();
      }
    };

    chrome.tabs.onUpdated.addListener(handleTabUpdate);
    chrome.tabs.onRemoved.addListener(checkCustomerOSTab);

    return () => {
      chrome.tabs.onUpdated.removeListener(handleTabUpdate);
      chrome.tabs.onRemoved.removeListener(checkCustomerOSTab);
    };
  }, []);

  const handleAddContact = async () => {
    if (!linkedInUrl) return;

    setIsLoading(true);

    const sessionData = await new Promise<{
      email: string;
      apiKey: string;
      workspaceName: string;
    } | null>((resolve) => {
      const messageHandler = (message: any) => {
        if (message.action === "COS_SESSION_DATA") {
          chrome.runtime.onMessage.removeListener(messageHandler);
          resolve({
            email: message.email,
            apiKey: message.apiKey,
            workspaceName: message.workspaceName,
          });
        }
      };

      chrome.runtime.onMessage.addListener(messageHandler);
      chrome.runtime.sendMessage({ action: "REQUEST_COS_SESSION_DATA" });
    });

    if (!sessionData?.apiKey) {
      setIsLoading(false);
      return;
    }
    setWorkspaceName(sessionData.workspaceName || null);

    try {
      const response = await fetch(
        "https://api.customeros.ai/browserExtension/contact",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Accept: "*/*",
            "X-CUSTOMER-OS-API-KEY": sessionData.apiKey,
          },
          body: JSON.stringify({ linkedinUrl: linkedInUrl }),
        }
      );

      if (response.ok) {
        setSuccessMessage("Contact added successfully");
        setTimeout(() => setSuccessMessage(null), 3000);
      } else {
        console.error("Failed to add contact");
      }
    } catch (error) {
      console.error("Error adding contact:", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const handleClick = (event: MouseEvent) => {
      event.preventDefault();
      chrome.runtime.sendMessage({
        action: "openTab",
        url: "https://app.customeros.ai",
      });
    };

    const link = document.getElementById("customerOSLink");
    if (link) {
      link.addEventListener("click", handleClick);
    }

    return () => {
      const link = document.getElementById("customerOSLink");
      if (link) {
        link.removeEventListener("click", handleClick);
      }
    };
  }, []);
  return (
    <div className="flex flex-col items-center justify-center h-full p-4">
      <div className="flex items-center flex-col gap-1">
        <img
          id="workspace-logo"
          src={chrome.runtime.getURL("src/assets/customeros.png")}
          alt="CustomerOS Logo"
          className="size-8"
        />
        <span className="font-semibold text-[16px]">CustomerOS</span>
        {workspaceName && appIsOpen && (
          <span className="bg-gray-100">Signed into {workspaceName}</span>
        )}
        {!linkedInUrl && (
          <span className="text-center max-w-[250px] ">
            Go to any LinkedIn profile to instantly add contacts to CustomerOS
          </span>
        )}
        {!appIsOpen && (
          <span className="text-center max-w-[250px] ">
            Open the{" "}
            <a id="customerOSLink" href="#" className="underline">
              CustomerOS app
            </a>{" "}
            to start adding contacts to your workspace instantly
          </span>
        )}
      </div>
      {linkedInUrl && workspaceName && (
        <div className="mt-4" onClick={handleAddContact}>
          <Button
            className="linkedin-button"
            size="xs"
            colorScheme="primary"
            isDisabled={isLoading}
            leftIcon={
              !isLoading ? (
                <img
                  id="plus-icon"
                  src={chrome.runtime.getURL("src/assets/plus.svg")}
                  alt="Plus Icon"
                />
              ) : (
                <img
                  src={chrome.runtime.getURL("src/assets/spinner.svg")}
                  alt="loading"
                />
              )
            }
          >
            {isLoading ? "Add contact to CustomerOS" : "Adding to CustomerOS…"}
          </Button>
        </div>
      )}
      {successMessage && (
        <div className="flex items-center gap-2 mt-2">
          <img
            src={chrome.runtime.getURL("src/assets/check-circle.svg")}
            alt="Success Icon"
            className="text-success-700"
          />
          <span id="success-message">{successMessage}</span>
        </div>
      )}
    </div>
  );
};
