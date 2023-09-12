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
    if (`${process.env.NEXT_PUBLIC_JUNE_ENABLED}` === 'true') {
      loadAnalytics();
    }
  }, []);

  return analytics;
}
