'use client';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';

export const HelpContent = () => {
  return (
    <Box mt='1'>
      <Text fontSize='md' fontWeight='normal'>
        Monthly Recurring Revenue (MRR) is the total amount of money you can
        expect to receive from your customers per month.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        This chart is the average across your live customer base, excluding
        churned customers and customers without recurring revenue.
      </Text>
    </Box>
  );
};
