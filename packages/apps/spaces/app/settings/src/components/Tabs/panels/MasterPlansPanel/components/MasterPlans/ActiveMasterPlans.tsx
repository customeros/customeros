import { useSearchParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Collapse } from '@ui/transitions/Collapse';
import { MasterPlansQuery } from '@shared/graphql/masterPlans.generated';

import { MasterPlans } from './MasterPlans';
import { MasterPlansMenu } from './MasterPlansMenu';
import { useMasterPlansMethods } from '../../hooks/useMasterPlansMethods';

interface ActiveMasterPlansProps {
  isLoading?: boolean;
  activePlans?: MasterPlansQuery['masterPlans'];
}

export const ActiveMasterPlans = ({
  isLoading,
  activePlans,
}: ActiveMasterPlansProps) => {
  const { isPending, handleCreateDefault, handleCreateFromScratch } =
    useMasterPlansMethods();

  const searchParams = useSearchParams();
  const isOpen = searchParams?.get('show') !== 'retired';

  return (
    <Flex flexDir='column' flex={isOpen ? 1 : 0}>
      <Flex align='center' justify='space-between' mb='2'>
        <Text fontWeight='semibold'>Your plans</Text>
        <MasterPlansMenu
          isLoading={isPending}
          onCreateDefault={handleCreateDefault}
          onCreateFromScratch={handleCreateFromScratch}
        />
      </Flex>

      <Collapse in={isOpen} animateOpacity>
        <MasterPlans isLoading={isLoading} masterPlans={activePlans} />
      </Collapse>
    </Flex>
  );
};
