'use client';

import { useMemo, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

import {
  MasterPlansQuery,
  useMasterPlansQuery,
} from '@settings/graphql/masterPlans.generated';

import { Grid, GridItem } from '@ui/layout/Grid';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  Milestones,
  ActiveMasterPlans,
  MasterPlanDetails,
  RetiredMasterPlans,
} from './components';

export const MasterPlansPanel = () => {
  const router = useRouter();
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

  const selectedPlan = plans.find((plan) => plan?.id === planId);
  const selectedMilestones = selectedPlan?.milestones ?? [];

  useEffect(() => {
    if (!planId) {
      const newParams = new URLSearchParams(searchParams ?? '');
      const firstId = showRetired
        ? retiredPlans?.[0]?.id
        : activePlans?.[0]?.id;

      if (!firstId) return;

      newParams.set('planId', firstId);
      router.push(`/settings?${newParams.toString()}`);
    }
  }, [showRetired, planId]);

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
        <RetiredMasterPlans
          isLoading={isLoading}
          retiredPlans={retiredPlans}
          activePlanFallbackId={activePlans[0]?.id}
          retiredPlanFallbackId={retiredPlans[0]?.id}
        />
      </GridItem>

      <GridItem p='4'>
        {planId && selectedPlan && (
          <>
            <MasterPlanDetails
              id={planId}
              isRetired={selectedPlan?.retired}
              name={selectedPlan?.name ?? 'Unnamed master plan'}
            />
            <Milestones milestones={selectedMilestones} />
          </>
        )}
      </GridItem>
    </Grid>
  );
};
