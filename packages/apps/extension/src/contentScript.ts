function sendSessionData() {
  // console.log("Attempting to send session data from content script");
  const request: IDBOpenDBRequest = indexedDB.open("customerDB_shared");

  request.onerror = function (event: Event) {
    console.error(
      "Error opening IndexedDB:",
      (event.target as IDBRequest)?.error
    );
  };

  request.onsuccess = function (event: Event) {
    // console.log("IndexedDB opened successfully");
    const db: IDBDatabase = (event.target as IDBOpenDBRequest).result;

    const transaction: IDBTransaction = db.transaction(["Session"], "readonly");

    const objectStore: IDBObjectStore = transaction.objectStore("Session");
    const getRequest: IDBRequest = objectStore.get("value");

    getRequest.onerror = function (event: Event) {
      console.error(
        "Error reading from IndexedDB:",
        (event.target as IDBRequest)?.error
      );
    };

    getRequest.onsuccess = function (event: Event) {
      const sessionData = (event.target as IDBRequest).result;
      // console.log("Session data retrieved from IndexedDB:", sessionData);

      if (sessionData && sessionData.profile) {
        const email: string | null = sessionData.profile.email || null;
        const apiKeyRequest: IDBRequest = objectStore.get("tenantApiKey");
        const tenantName: string | null = sessionData.tenant;
        apiKeyRequest.onerror = function (event: Event) {
          console.error(
            "Error reading tenantApiKey from IndexedDB:",
            (event.target as IDBRequest)?.error
          );
        };

        apiKeyRequest.onsuccess = function (event: Event) {
          const apiKey = (event.target as IDBRequest).result;

          if (email && apiKey) {
            chrome.runtime.sendMessage({
              action: "COS_SESSION_DATA",
              email,
              apiKey,
              tenantName,
            });
          }
        };
      } else {
        console.log(
          "No session data found in IndexedDB or session data is incomplete"
        );
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
