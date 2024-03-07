import React from 'react';

import { Heading } from '@ui/typography/Heading';

export const NotificationsHeader: React.FC = () => {
  return (
    <Heading fontSize='md' px={4} py={1} mb={3}>
      Up next
    </Heading>
  );
};
