async function initMocks() {
  if (typeof window === 'undefined') {
    console.log('🏷️ ----- : SERVER');
    const { server } = await import('./server');
    server.listen();
  } else {
    console.log('🏷️ ----- : BROWSER');
    const { worker } = await import('./browser');
    worker.start();
  }
}

initMocks();

export {};
