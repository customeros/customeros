import { useState, useEffect } from 'react';

import { loadingStatements } from './messages.ts';

export const LoadingMessage: React.FC = () => {
  const getRandomMessage = () => {
    const randomIndex = Math.floor(Math.random() * loadingStatements.length);

    return loadingStatements[randomIndex];
  };

  const [loadingMessage, setLoadingMessage] = useState<string>(
    getRandomMessage(),
  );

  useEffect(() => {
    const intervalId = setInterval(() => {
      setLoadingMessage(getRandomMessage());
    }, 3000);

    return () => clearInterval(intervalId);
  }, []);

  return <div className='text-sm '>{loadingMessage}</div>;
};
