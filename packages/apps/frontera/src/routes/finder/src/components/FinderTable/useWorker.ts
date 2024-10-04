import { useState, useEffect } from 'react';

const useWebWorker = (workerFunction: () => void, inputData: unknown) => {
  const [result, setResult] = useState(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!inputData) return;

    setLoading(true);
    setError(null);
    setResult(null);

    const blob = new Blob(
      [
        `
          self.onmessage = function(event) {
            (${workerFunction})(event);
          };
        `,
      ],
      { type: 'application/javascript' },
    );

    const workerScriptUrl = URL.createObjectURL(blob);
    const worker = new Worker(workerScriptUrl);

    worker.onmessage = (event) => {
      console.info('Worker result:', event.data);
      setResult(event.data);
      setLoading(false);
    };

    worker.onerror = (event) => {
      console.error('Worker error:', event.message);
      setError(event.message);
      setLoading(false);
    };

    console.info('Posting message to worker:', inputData);
    worker.postMessage(inputData);

    return () => {
      worker.terminate();
      URL.revokeObjectURL(workerScriptUrl);
    };
  }, [inputData, workerFunction]);

  return { result, error, loading };
};

export default useWebWorker;
