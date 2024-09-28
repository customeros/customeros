function sendSessionData() {
  console.log("Attempting to send session data from content script");
  const request: IDBOpenDBRequest = indexedDB.open("customerDB", 2);

  request.onerror = function (event: Event) {
    console.error(
      "Error opening IndexedDB:",
      (event.target as IDBRequest)?.error
    );
  };

  request.onsuccess = function (event: Event) {
    const db: IDBDatabase = (event.target as IDBOpenDBRequest).result;

    const transaction: IDBTransaction = db.transaction(
      ["customer_os"],
      "readonly"
    );
    const objectStore: IDBObjectStore = transaction.objectStore("customer_os");

    const getRequest: IDBRequest = objectStore.get("SessionStore");

    getRequest.onerror = function (event: Event) {
      console.error(
        "Error reading from IndexedDB:",
        (event.target as IDBRequest)?.error
      );
    };

    getRequest.onsuccess = function (event: Event) {
      const sessionData = (event.target as IDBRequest).result;
      console.log("Session data retrieved from IndexedDB:", sessionData);

      if (sessionData) {
        const email: string | null = sessionData.value.profile.email || null;
        const apiKey: string | null = sessionData.tenantApiKey || null;

        console.log("Sending session data to background:", { email, apiKey });
        chrome.runtime.sendMessage({
          action: "COS_SESSION_DATA",
          email,
          apiKey,
        });
      } else {
        console.log("No session data found in IndexedDB");
      }
    };
  };

  request.onupgradeneeded = function (event: IDBVersionChangeEvent) {
    const db: IDBDatabase = (event.target as IDBOpenDBRequest).result;
    if (!db.objectStoreNames.contains("customer_os")) {
      db.createObjectStore("customer_os");
    }
  };
}

sendSessionData();

// Add a listener for messages from the background script
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === "GET_SESSION_DATA") {
    console.log("Content script received GET_SESSION_DATA request");
    sendSessionData();
  }
});
