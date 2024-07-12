(() => {
  const loader = document.querySelector('[data-customeros]');
  loader.parentNode.removeChild(loader);

  let trackerData = JSON.parse(
    localStorage.getItem('customeros-tracker') ??
      `{ "identity": { "sessionId": "${crypto.randomUUID()}" }, "activity": "" }`,
  );

  function persist() {
    if (trackerData.activity.startsWith(',')) {
      trackerData.activity = trackerData.activity.substring(1);
    }
    localStorage.setItem('customeros-tracker', JSON.stringify(trackerData));
  }

  function identify(identifier) {
    trackerData.identity = { ...trackerData.identity, identifier };
    persist();

    const payload = JSON.stringify(trackerData);

    fetch(`https://user-admin-api.customeros.ai/track`, {
      headers: {
        'Content-Type': 'application/json',
        'X-Tracker-Payload': payload,
      },
    });
  }

  function track(event) {
    trackerData.activity = trackerData.activity
      .split(',')
      .concat([new Date().valueOf(), event])
      .join(',');

    persist();
  }

  window.CustomerOS = {
    identify,
    track,
  };
})();
