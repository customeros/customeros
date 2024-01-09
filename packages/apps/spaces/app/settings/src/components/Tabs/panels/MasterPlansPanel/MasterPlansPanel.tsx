'use client';

import { useRouter, useSearchParams } from 'next/navigation';

import {
  MasterPlansQuery,
  useMasterPlansQuery,
} from '@settings/graphql/masterPlans.generated';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { Grid, GridItem } from '@ui/layout/Grid';
import { Collapse } from '@ui/transitions/Collapse';
import { Skeleton } from '@ui/presentation/Skeleton';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const MasterPlansPanel = () => {
  const client = getGraphQLClient();
  const { data, isLoading } = useMasterPlansQuery(client);

  const [activePlans, retiredPlans] = (data?.masterPlans ?? []).reduce(
    (acc, curr) => {
      if (curr?.retired) {
        acc[1]?.push(curr);
      } else {
        acc[0]?.push(curr);
      }

      return acc;
    },
    [[], []] as [
      MasterPlansQuery['masterPlans'],
      MasterPlansQuery['masterPlans'],
    ],
  );

  return (
    <Grid templateColumns='1fr 2fr' h='full'>
      <GridItem
        p='4'
        display='flex'
        flexDir='column'
        borderRight='1px solid'
        borderRightColor='gray.200'
      >
        <ActiveMasterPlans isLoading={isLoading} activePlans={activePlans} />
        <RetiredMasterPlans isLoading={isLoading} retiredPlans={retiredPlans} />
      </GridItem>
    </Grid>
  );
};

interface MasterPlanNavItemProps {
  id: string;
  children?: React.ReactNode;
}

const MasterPlanItem = ({ id, children }: MasterPlanNavItemProps) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isActive = searchParams?.get('planId') === id;

  const handleClick = () => {
    const newParams = new URLSearchParams(searchParams ?? '');
    newParams.set('planId', id);

    router.push(`?${newParams.toString()}`);
  };

  return (
    <Button
      px='3'
      w='full'
      fontSize='sm'
      fontWeight='normal'
      onClick={handleClick}
      justifyContent='flex-start'
      bg={isActive ? 'gray.100' : 'transparent'}
      _hover={{
        bg: 'gray.100',
      }}
      _active={{
        bg: 'gray.200',
      }}
    >
      {children}
    </Button>
  );
};

interface MasterPlansProps {
  isLoading?: boolean;
  masterPlans?: MasterPlansQuery['masterPlans'];
}

const MasterPlans = ({ masterPlans, isLoading }: MasterPlansProps) => {
  if (isLoading) return <LoadingMasterplans />;
  if (!masterPlans) return <NoMasterplans />;

  return (
    <VStack align='flex-start'>
      {masterPlans.map(({ id, name }) => (
        <MasterPlanItem key={id} id={id}>
          {name}
        </MasterPlanItem>
      ))}
    </VStack>
  );
};

interface ActiveMasterPlansProps {
  isLoading?: boolean;
  activePlans?: MasterPlansQuery['masterPlans'];
}

const ActiveMasterPlans = ({
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

interface RetiredMasterPlansProps {
  isLoading?: boolean;
  retiredPlans?: MasterPlansQuery['masterPlans'];
}

const RetiredMasterPlans = ({
  isLoading,
  retiredPlans,
}: RetiredMasterPlansProps) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isOpen = searchParams?.get('show') === 'retired';

  const toggle = () => {
    const newParams = new URLSearchParams(searchParams ?? '');
    newParams.set('show', isOpen ? 'active' : 'retired');

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

const LoadingMasterplans = () => {
  return (
    <VStack align='flex-start'>
      {Array.from({ length: 5 }).map((_, i) => (
        <Skeleton key={i} h='6' w='full' />
      ))}
    </VStack>
  );
};

const NoMasterplans = () => {
  return <Flex>No master plans created yet</Flex>;
};
