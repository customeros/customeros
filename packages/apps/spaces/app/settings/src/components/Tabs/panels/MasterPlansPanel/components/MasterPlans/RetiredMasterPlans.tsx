import { useRouter, useSearchParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Collapse } from '@ui/transitions/Collapse';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { MasterPlansQuery } from '@shared/graphql/masterPlans.generated';

import { MasterPlans } from './MasterPlans';

interface RetiredMasterPlansProps {
  isLoading?: boolean;
  activePlanFallbackId?: string;
  retiredPlanFallbackId?: string;
  retiredPlans?: MasterPlansQuery['masterPlans'];
}

export const RetiredMasterPlans = ({
  isLoading,
  retiredPlans,
  activePlanFallbackId,
  retiredPlanFallbackId,
}: RetiredMasterPlansProps) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isOpen = searchParams?.get('show') === 'retired';

  const toggle = () => {
    const newParams = new URLSearchParams(searchParams ?? '');
    newParams.set('show', isOpen ? 'active' : 'retired');

    newParams.set(
      'planId',
      isOpen ? activePlanFallbackId ?? '' : retiredPlanFallbackId ?? '',
    );

    router.push(`?${newParams.toString()}`);
  };

  return (
    <Flex flexDir='column' flex={isOpen ? 1 : 0}>
      <Button
        mt='4'
        w='full'
        size='sm'
        variant='ghost'
        onClick={toggle}
        justifyContent='space-between'
        rightIcon={
          isOpen ? (
            <ChevronDown color='gray.400' />
          ) : (
            <ChevronRight color='gray.400' />
          )
        }
        colorScheme='gray'
        sx={{
          '> span': {
            '> span': {
              ml: 1,
              color: 'gray.500',
            },
          },
        }}
      >
        <span>
          Retired plans<span>• {retiredPlans?.length}</span>
        </span>
      </Button>
      <Collapse in={isOpen} animateOpacity>
        <MasterPlans isLoading={isLoading} masterPlans={retiredPlans} />
      </Collapse>
    </Flex>
  );
};
