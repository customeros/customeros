import { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Inbox01 } from '@ui/media/icons/Inbox01';
export const EmptyIssueMessage: FC<{
  title?: string;
  description: string;
}> = ({ title, description }) => (
  <Flex direction='column' alignItems='center' mt='4'>
    <FeaturedIcon size='md' minW='10' colorScheme='gray' mb={2}>
      <Inbox01 color='gray.700' boxSize='6' />
    </FeaturedIcon>
    {title && (
      <Text color='gray.700' fontWeight={600} mb={1}>
        {title}
      </Text>
    )}

    <Text color='gray.500' mt={1} mb={6} textAlign='center'>
      {description}
    </Text>
  </Flex>
);
