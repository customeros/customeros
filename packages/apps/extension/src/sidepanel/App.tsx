import { useEffect, useState } from "react";
import "./styles/tailwind.css";
import { Button } from "@ui/form/Button/Button";

export const App = () => {
  const [workspaceName, setWorkspaceName] = useState<string | null>(null);
  const [_appIsOpen, setAppIsOpen] = useState(false);
  const [linkedInUrl, setLinkedInUrl] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
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
      if (changeInfo.url && tab.active) {
        updateLinkedInUrlFromTab(tab);
      }
    };

    const handleActiveTabChange = (activeInfo: chrome.tabs.TabActiveInfo) => {
      chrome.tabs.get(activeInfo.tabId, (tab) => {
        updateLinkedInUrlFromTab(tab);
      });
    };

    const updateLinkedInUrlFromTab = (tab: chrome.tabs.Tab) => {
      const isLinkedInProfile = tab.url?.includes("linkedin.com/in");
      setLinkedInUrl(isLinkedInProfile ? tab.url || null : null);
    };

    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      if (tabs[0]) {
        updateLinkedInUrlFromTab(tabs[0]);
      }
    });

    chrome.tabs.onUpdated.addListener(handleTabUpdate);
    chrome.tabs.onActivated.addListener(handleActiveTabChange);

    return () => {
      chrome.tabs.onUpdated.removeListener(handleTabUpdate);
      chrome.tabs.onActivated.removeListener(handleActiveTabChange);
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
        setSuccessMessage("Contact added");

        setTimeout(() => setSuccessMessage(null), 3000);
      } else {
        console.error("Failed to add contact");
      }
    } catch (error) {
      console.error("Error adding contact:", error);
      setErrorMessage("We couldn't add this contact");
      setTimeout(() => setErrorMessage(null), 3000);
    } finally {
      setIsLoading(false);
    }
  };

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
        {!workspaceName && !linkedInUrl && (
          <span className="text-center max-w-[250px]">
            Sign into the{" "}
            <a
              id="customerOSLink"
              href="#"
              className="underline"
              onClick={() => {
                chrome.runtime.sendMessage({
                  action: "openTab",
                  url: "https://app.customeros.ai",
                });
              }}
            >
              CustomerOS app
            </a>{" "}
            and go to any LinkedIn profile to start adding contacts to your
            workspace
          </span>
        )}
        {!workspaceName && linkedInUrl && (
          <span className="text-center max-w-[250px] ">
            Sign into the{" "}
            <a
              id="customerOSLink"
              href="#"
              className="underline"
              onClick={() => {
                chrome.runtime.sendMessage({
                  action: "openTab",
                  url: "https://app.customeros.ai",
                });
              }}
            >
              CustomerOS app
            </a>{" "}
            app to start adding contacts to your workspace
          </span>
        )}
        {workspaceName && (
          <span className="bg-gray-100">Signed into {workspaceName}</span>
        )}
        {!linkedInUrl && workspaceName && (
          <span className="text-center max-w-[250px] ">
            Go to any LinkedIn profile to instantly add contacts to CustomerOS
          </span>
        )}

        {workspaceName && linkedInUrl && (
          <span className="text-center max-w-[250px] ">
            Instantly add LinkedIn contacts to CustomerOS{" "}
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
            {!isLoading ? "Add contact to CustomerOS" : "Adding to CustomerOS…"}
          </Button>
        </div>
      )}
      {successMessage && (
        <div className="flex items-center gap-2 mt-2">
          <img
            src={chrome.runtime.getURL("src/assets/check-circle.svg")}
            alt="Success Icon"
            className="text-success-500"
          />
          <span id="success-message" className="text-success-500">
            {successMessage}
          </span>
        </div>
      )}
      {errorMessage && (
        <div className="flex items-center gap-2 mt-2">
          <img
            src={chrome.runtime.getURL("src/assets/x-circle.svg")}
            alt="Error Icon"
            className="text-error-500"
          />
          <span id="error-message" className="text-error-500">
            {errorMessage}
          </span>
        </div>
      )}
    </div>
  );
};
