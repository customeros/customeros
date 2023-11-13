'use client';

import { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Star06 } from '@ui/media/icons/Star06';
import { Heading } from '@ui/typography/Heading';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

export const EmptyContracts: FC<{ name: string }> = ({ name }) => {
  return (
    <OrganizationPanel
      title='Account'
      actionItem={
        <Button
          size='xs'
          variant='outline'
          type='button'
          isDisabled
          borderRadius='16px'
        >
          Prospect
        </Button>
      }
    >
      <Flex
        mt={4}
        w='full'
        boxShadow={'none'}
        flexDir='column'
        justifyItems='center'
        alignItems='center'
      >
        <FeaturedIcon colorScheme='primary' mb={2} size='lg'>
          <Star06 boxSize={4} />
        </FeaturedIcon>
        <Heading mb={1} size='sm' fontWeight='semibold'>
          Create new contract
        </Heading>
        <Text fontSize='sm'>
          Create new contract for
          <Text as='span' fontWeight='medium' ml={1}>
            {name}
          </Text>
        </Text>
        <Button
          fontSize='sm'
          size='sm'
          colorScheme='primary'
          mt={6}
          variant='outline'
          width='fit-content'
        >
          New contract
        </Button>
      </Flex>
    </OrganizationPanel>
  );
};
