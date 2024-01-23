import { useEffect, useCallback } from 'react';

import throttle from 'lodash/throttle';

type Callback<T> = (...args: T[]) => void;

export function useThrottle<T>(
  callback: Callback<T>,
  time = 500,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  deps: any[] = [],
) {
  const throttled = useCallback(
    throttle(callback, time, { trailing: true }),
    deps,
  );

  useEffect(() => {
    return () => {
      throttled.cancel();
    };
  }, []);

  return throttled;
}
