import { useEffect, useState } from "react";
import "./styles/tailwind.css";
import { Button } from "@ui/form/Button/Button";

type Contact = {
  contactId: string;
  email: string;
  linkedinUrl: string;
};

export const App = () => {
  const [workspaceName, setWorkspaceName] = useState<string | null>(null);
  const [tenantApiKey, setTenantApiKey] = useState<string | null>(null);
  const [linkedInUrl, setLinkedInUrl] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [contact, setContact] = useState<Contact | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isAddingContact, setIsAddingContact] = useState(false);

  useEffect(() => {
    const handleMessage = (message: any) => {
      if (message.action === "COS_SESSION_DATA") {
        setWorkspaceName(message.workspaceName || null);
        setTenantApiKey(message.apiKey || null);
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
    setIsLoading(true);

    const handleActiveTabChange = (activeInfo: chrome.tabs.TabActiveInfo) => {
      chrome.tabs.get(activeInfo.tabId, (tab) => {
        updateLinkedInUrlFromTab(tab);
      });
    };

    const parseLinkedInUrl = (url: string): string | null => {
      try {
        if (url.includes("linkedin.com/sales/lead/")) {
          const commaIndex = url.indexOf(",");
          if (commaIndex !== -1) {
            return url.slice(0, commaIndex);
          }
          return url;
        } else if (url.includes("linkedin.com/in/")) {
          return url;
        }
        return null;
      } catch {
        return null;
      }
    };

    const updateLinkedInUrlFromTab = async (tab: chrome.tabs.Tab) => {
      const isLinkedInProfile =
        tab.url?.includes("linkedin.com/in") ||
        tab.url?.includes("linkedin.com/sales/lead");

      setContact(null);

      const parsedUrl = tab.url ? parseLinkedInUrl(tab.url) : null;
      setLinkedInUrl(parsedUrl);

      if (!isLinkedInProfile || !parsedUrl) {
        setIsLoading(false);
        return;
      }

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

      if (sessionData?.apiKey) {
        try {
          const response = await fetch(
            `https://api.customeros.ai/browserExtension/contact?linkedin=${parsedUrl}`,
            {
              method: "GET",
              headers: {
                "Content-Type": "application/json",
                Accept: "*/*",
                "X-CUSTOMER-OS-API-KEY": sessionData?.apiKey,
              },
            }
          );
          if (response.ok) {
            const data = await response.json();
            setContact(data.contact);
          }
          if (
            !tab.url?.includes("linkedin.com/in") &&
            !tab.url?.includes("linkedin.com/sales/lead")
          ) {
            setContact(null);
          }
        } catch (error) {
          console.error("Error fetching contact:", error);
          setContact(null);
        } finally {
          setIsLoading(false);
        }
      }
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

    setIsAddingContact(true);

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
      setIsAddingContact(false);
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
        setErrorMessage("We couldn't add this contact");
        setTimeout(() => setErrorMessage(null), 5000);
      }
    } catch (error) {
      console.error("Error adding contact:", error);
      setErrorMessage("We couldn't add this contact");
      setTimeout(() => setErrorMessage(null), 5000);
    } finally {
      setIsAddingContact(false);
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
        {!tenantApiKey && !linkedInUrl && !isLoading && (
          <span className="text-center max-w-[250px] text-sm">
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
        {!tenantApiKey && linkedInUrl && !isLoading && (
          <span className="text-center max-w-[250px] text-sm ">
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
        {tenantApiKey && (
          <span className="bg-gray-100 px-1 text-sm">
            Signed into {workspaceName || "CustomerOS"}
          </span>
        )}
        {!linkedInUrl && tenantApiKey && !isLoading && (
          <span className="text-center max-w-[250px] text-sm ">
            Go to any LinkedIn profile to instantly add contacts to CustomerOS
          </span>
        )}

        {tenantApiKey && linkedInUrl && !contact?.contactId && !isLoading && (
          <span className="text-center max-w-[250px] text-sm">
            Instantly add LinkedIn contacts to CustomerOS{" "}
          </span>
        )}
      </div>
      {isLoading && !contact?.contactId && (
        <img
          className="mt-4"
          src={chrome.runtime.getURL("src/assets/spinner.svg")}
          alt="loading"
        />
      )}

      {linkedInUrl && tenantApiKey && !contact?.contactId && !isLoading && (
        <Button
          className="mt-4"
          size="xs"
          onClick={handleAddContact}
          colorScheme="primary"
          isDisabled={isAddingContact}
          leftIcon={
            !isAddingContact ? (
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
          {!isAddingContact
            ? "Add contact to CustomerOS"
            : "Adding to CustomerOS…"}
        </Button>
      )}
      {contact?.contactId && (
        <div className="mt-1">
          <span className="text-sm">This contact is already in CustomerOS</span>
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
