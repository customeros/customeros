import { useSearchParams } from 'next/navigation';

import { MasterPlansQuery } from '@settings/graphql/masterPlans.generated';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { Collapse } from '@ui/transitions/Collapse';

import { MasterPlans } from './MasterPlans';

interface ActiveMasterPlansProps {
  isLoading?: boolean;
  activePlans?: MasterPlansQuery['masterPlans'];
}

export const ActiveMasterPlans = ({
  isLoading,
  activePlans,
}: ActiveMasterPlansProps) => {
  const searchParams = useSearchParams();
  const isOpen = searchParams?.get('show') !== 'retired';

  return (
    <Flex flexDir='column' flex={isOpen ? 1 : 0}>
      <Flex align='center' justify='space-between' mb='2'>
        <Text fontWeight='semibold'>Your plans</Text>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Add Master Plan'
          icon={<Plus color='gray.400' />}
        />
      </Flex>

      <Collapse in={isOpen} animateOpacity>
        <MasterPlans isLoading={isLoading} masterPlans={activePlans} />
      </Collapse>
    </Flex>
  );
};
