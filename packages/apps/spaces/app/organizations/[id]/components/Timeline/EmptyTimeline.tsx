import React from 'react';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import EmptyTimelineIlustration from '@spaces/atoms/icons/EmptyTimelineIlustration';
import { useOrganization } from '@organization/hooks/useOrganization';
import { useParams } from 'next/navigation';

export const EmptyTimeline: React.FC = () => {
  const id = useParams()?.id as string;

  const { data } = useOrganization({ id });

  return (
    <Flex
      direction='column'
      alignItems='center'
      flex={1}
      backgroundImage='/backgrounds/organization/dotted-bg-pattern.svg'
      backgroundRepeat='no-repeat'
      backgroundSize='contain'
      backgroundPosition='center'
      maxH='50%'
      as='article'
    >
      <Flex
        direction='column'
        alignItems='center'
        justifyContent='center'
        height='100%'
        maxWidth='390px'
      >
        <EmptyTimelineIlustration />
        <Text
          color='gray.900'
          fontSize='lg'
          as='h1'
          fontWeight={600}
          mt={3}
          mb={2}
        >
          {data?.organization?.name || 'Unknown'} has no events yet
        </Text>
        <Text color='gray.600' size='xs' textAlign='center'>
          This organization’s events will show up here once a data source has
          been linked
        </Text>
      </Flex>
    </Flex>
  );
};
