'use client';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';

export const HelpContent = () => {
  return (
    <Box mt='1'>
      <Text fontSize='md' fontWeight='normal'>
        Revenue at risk shows the forecasted revenue from customers whose
        renewal likelihood is rated medium, low or zero in the current period.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        In contrast, the high confidence segment shows the forecasted revenue
        with a high likelihood to renew.
      </Text>
    </Box>
  );
};
