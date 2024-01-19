'use client';

import { useMemo } from 'react';
import { useSearchParams } from 'next/navigation';

import {
  MasterPlansQuery,
  useMasterPlansQuery,
} from '@settings/graphql/masterPlans.generated';

import { Flex } from '@ui/layout/Flex';
import { Grid, GridItem } from '@ui/layout/Grid';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  Milestones,
  ActiveMasterPlans,
  MasterPlanDetails,
  RetiredMasterPlans,
} from './components';

export const MasterPlansPanel = () => {
  const client = getGraphQLClient();
  const { data, isLoading } = useMasterPlansQuery(client);

  const searchParams = useSearchParams();
  const planId = searchParams?.get('planId');
  const showRetired = searchParams?.get('show') === 'retired';

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

  const plans = useMemo(() => {
    if (showRetired) return retiredPlans;

    return activePlans;
  }, [data?.masterPlans, showRetired]);

  if (!plans) return <Flex>No master plan created yet</Flex>;
  if (!planId) return <Flex>Select a master plan</Flex>;

  const selectedPlan = plans.find((plan) => plan?.id === planId);
  const selectedMilestones = selectedPlan?.milestones ?? [];

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

      <GridItem p='4'>
        <MasterPlanDetails name={selectedPlan?.name ?? 'Unnamed master plan'} />
        <Milestones milestones={selectedMilestones} />
      </GridItem>
    </Grid>
  );
};
