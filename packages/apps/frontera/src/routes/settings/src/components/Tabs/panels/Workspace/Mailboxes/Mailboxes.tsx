import React, { useState } from 'react';

import { EmptyMailboxes } from './components/EmptyMailboxes';

export const Mailboxes = () => {
  const [isUpdated, setIsUpdated] = useState(false);

  const handleUpdate = () => {
    setIsUpdated(true);
  };

  return (
    <>
      {!isUpdated ? (
        <EmptyMailboxes onUpdate={handleUpdate} />
      ) : (
        <div>aici</div>
      )}
    </>
  );
};
