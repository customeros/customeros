import { useEffect, useState } from 'react';
import { AnalyticsBrowser } from '@june-so/analytics-next';

export function useJune(): AnalyticsBrowser | undefined {
  const [analytics, setAnalytics] = useState(undefined);

  useEffect(() => {
    const loadAnalytics = async () => {
      const response = AnalyticsBrowser.load({
        writeKey: 'M2QnaR2vqHiuu3W2',
      }) as any;
      setAnalytics(response);
    };

    // @ts-expect-error fixme type process in the future
    if (process.env.JUNE_ENABLED === true) {
      loadAnalytics();
    }
  }, []);

  return analytics;
}
