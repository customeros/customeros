import { useState } from 'react';
import { useParams } from 'next/navigation';

import { VStack } from '@ui/layout/Stack';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationOnboardingPlansQuery } from '@organization/src/graphql/organizationOnboardingPlans.generated';

import { Plan } from './Plan';

export const OnboardingPlans = () => {
  const [openPlanId, setOpenPlanId] = useState<string | null>(null);
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();
  const { data } = useOrganizationOnboardingPlansQuery(client, {
    organizationId,
  });

  const activePlans = data?.organizationPlansForOrganization?.filter(
    (plan) => !plan.retired,
  );

  const handleTogglePlan = (planId: string) => {
    setOpenPlanId((prevPlanId) => (prevPlanId === planId ? null : planId));
  };

  return (
    <VStack w='full' overflowY='auto' maxH='calc(100vh - 148px)' mt='4'>
      {activePlans?.map((plan) => (
        <Plan
          plan={plan}
          key={plan.id}
          onToggle={handleTogglePlan}
          isOpen={plan.id === openPlanId}
        />
      ))}
    </VStack>
  );
};
