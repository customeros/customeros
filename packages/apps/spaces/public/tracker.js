(() => {
  console.info('Hello from tracker.js')
  const loader = document.querySelector('[data-customeros]')
  loader.parentNode.removeChild(loader)

  let trackerData = JSON.parse(localStorage.getItem('customeros-tracker') ?? `{ "identity": { "sessionId": "${crypto.randomUUID()}" }, "activity": [] }`);

  function persist() {
    localStorage.setItem('customeros-tracker', JSON.stringify(trackerData))
  }

  function identify(identifier) {
    trackerData.identity = { ...trackerData.identity, identifier}
    persist();

    const payload = JSON.stringify(trackerData)

    fetch(`http://localhost:3001/api/tracker?event=identify`, {
      headers: {
        'X-Tracker-Payload': payload,
      }
    });
  }

  function track(event) {
    trackerData.activity.push(event)
    persist();
  }

  window.CustomerOS = {
    identify,
    track
  }
})()
