import React from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Lotus } from '@ui/media/icons/Lotus';
import { Heading } from '@ui/typography/Heading';

import HalfCirclePattern from '../../assets/HalfCirclePattern';

interface EmptyNotificationsProps {}

export const EmptyNotifications: React.FC<EmptyNotificationsProps> = () => {
  return (
    <Flex
      as='article'
      position='relative'
      flexDirection='column'
      alignItems='center'
      maxW='448px'
      px={4}
      py={1}
      mt={5}
      overflow='hidden'
    >
      <Box position='absolute' height='400px' width='448px'>
        <HalfCirclePattern />
      </Box>
      <FeaturedIcon colorScheme='primary'>
        <Lotus />
      </FeaturedIcon>
      <Heading fontSize='md' mt={4} mb={1}>
        No notifications for now
      </Heading>
      <Text textAlign='center' fontSize='sm' color='gray.500'>
        Enjoy the quiet moment. Explore other corners of the app or take a deep
        breath and savor the calm.
      </Text>
    </Flex>
  );
};
