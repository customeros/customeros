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
    console.log(
      '🏷️ ----- process.env.JUNE_ENABLED: ',
      process.env.JUNE_ENABLED,
    );
    if (`${process.env.JUNE_ENABLED}` === 'true') {
      console.log(
        '🏷️ ----- : Test log that will be remove as soon as issue is resolved on prod',
      );
      loadAnalytics();
    }
  }, []);

  return analytics;
}
