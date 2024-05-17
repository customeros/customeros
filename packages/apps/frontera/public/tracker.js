(() => {
  const loader = document.querySelector('[data-customeros]');
  loader.parentNode.removeChild(loader);

  let trackerData = JSON.parse(
    localStorage.getItem('customeros-tracker') ??
      `{ "identity": { "sessionId": "${crypto.randomUUID()}" }, "activity": "" }`,
  );

  function persist() {
    localStorage.setItem('customeros-tracker', JSON.stringify(trackerData));
  }

  function identify(identifier) {
    trackerData.identity = { ...trackerData.identity, identifier };
    persist();

    const payload = JSON.stringify(trackerData);

    fetch(`https://user-admin-api.customeros.ai/track`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Tracker-Payload': payload,
      },
    });
  }

  function track(event) {
    trackerData.activity = trackerData.activity.concat(
      `${new Date().valueOf()},${event}`,
    );
    persist();
  }

  window.CustomerOS = {
    identify,
    track,
  };
})();
