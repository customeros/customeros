'use client';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';

export const HelpContent = () => {
  return (
    <Box mt='1'>
      <Text fontSize='md' fontWeight='normal'>
        Retention Rate is the percentage of customers who continue to subscribe
        to your service over a specific period.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        To calculate this rate, we look at the number of customers with a
        renewal in the current period and determine what percentage in fact
        renewed their subscription.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        For example, customers with an Annual renewal in January will only be
        included in this metric during the January period as they are not
        eligible for renewal from February to December.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        The higher this percentage, the more effectively you are maintaining
        your customer base.
      </Text>
    </Box>
  );
};
