'use client';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';

export const HelpContent = () => {
  return (
    <Box mt='1'>
      <Text fontSize='md' fontWeight='normal'>
        This is the number of new customers you’ve added, meaning that their
        contract’s start date falls in the given period.
      </Text>
    </Box>
  );
};
