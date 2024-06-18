import { useState, useLayoutEffect } from 'react';

import { ProgressBar } from '@ui/feedback/Progress/ProgressBar.tsx';

export const SimulatedProgress = ({ accelerate }: { accelerate: boolean }) => {
  const [progress, setProgress] = useState(0);

  useLayoutEffect(() => {
    if (!accelerate) {
      let start = 0;
      const increment = 100 / (10000 / 100); //  duration 10sek
      const interval = setInterval(() => {
        start += increment;
        if (start >= 100) {
          start = 100;
          clearInterval(interval);
        }
        setProgress(start);
      }, 100);

      return () => clearInterval(interval);
    }
    if (accelerate) {
      setProgress(100);
    }
  }, [accelerate]);

  return <ProgressBar value={progress} className='h-1 w-full' />;
};
