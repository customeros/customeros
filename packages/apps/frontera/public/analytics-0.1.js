(function (w) {
  function generateUUID() {
    return (
      'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = (Math.random() * 16) | 0,
          v = c === 'x' ? r : (r & 0x3) | 0x8;

        return v.toString(16);
      }) + new Date().valueOf()
    );
  }

  function getIp() {
    return fetch('https://api.ipify.org?format=json')
      .then((response) => response.json())
      .then((data) => {
        window.cosUserIp = data.ip;
      });
  }

  function sendData(eventType, eventData) {
    const userAgent = navigator.userAgent;

    fetch('https://user-admin-api.customeros.ai/tracking', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      keepalive: true,
      body: JSON.stringify({
        ip: window.cosUserIp,
        userId: window.cosUserId,
        eventType: eventType,
        eventData: JSON.stringify(eventData),
        timestamp: new Date().valueOf(),
        href: window.location.href,
        origin: window.location.origin,
        search: window.location.search,
        hostname: window.location.hostname,
        pathname: window.location.pathname,
        referrer: document.referrer,
        userAgent: userAgent,
        language: navigator.language,
        cookiesEnabled: navigator.cookieEnabled,
        screenResolution: window.screen.width + 'x' + window.screen.height,
      }),
    });
  }

  var cosUserId = document.cookie.replace(
    /(?:(?:^|.*;\s*)cosUserId\s*=\s*([^;]*).*$)|^.*$/,
    '$1',
  );

  if (!cosUserId) {
    cosUserId = generateUUID();
    document.cookie = 'cosUserId=' + cosUserId + '; path=/';
  }
  window.cosUserId = cosUserId;

  getIp()
    .then(() => {
      sendData('page_view', {
        title: document.title,
      });
      window.cosPageLoadTime = new Date().valueOf();

      document.addEventListener('click', function (event) {
        const target = event.target;

        if (target.tagName === 'BODY') {
          return;
        }

        sendData('click', {
          tag: target.tagName,
          id: target.id,
          classes: target.className,
          text: target.innerText,
          url: target.href,
        });
      });

      document.addEventListener(
        'blur',
        function (event) {
          const target = event.target;

          if (target.tagName !== 'INPUT' && target.tagName !== 'TEXTAREA') {
            return;
          }

          if (target.type !== 'email') {
            return;
          }

          if (!target.value) {
            return;
          }

          sendData('identify', {
            tag: target.tagName,
            id: target.id,
            classes: target.className,
            email: target.value,
            dataset: target.dataset,
          });
        },
        true,
      );

      document.addEventListener('visibilitychange', function (event) {
        if (document.visibilityState === 'hidden') {
          sendData('page_exit', {
            title: document.title,
            sessionDuration: new Date().valueOf() - window.cosPageLoadTime,
          });

          event.preventDefault();
        }
      });
    })
    .catch((error) => {
      console.error('Error fetching IP:', error);
    });
})(window);
