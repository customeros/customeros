'use client';
import { useParams, useRouter } from 'next/navigation';

import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';

import { PanelContainer } from './PanelContainer';
import { OnboardingStatus } from './OnboardingStatus';

export const SuccessPanel = () => {
  const router = useRouter();
  const client = getGraphQLClient();
  const id = useParams()?.id as string;
  const { data } = useOrganizationQuery(client, { id });
  const isFeatureOn = useFeatureIsOn('onboarding-status');

  if (!isFeatureOn) {
    router.replace(`/organization/${id}?tab=about`);
  }

  return (
    <PanelContainer title='Success'>
      <OnboardingStatus data={data?.organization?.accountDetails?.onboarding} />
    </PanelContainer>
  );
};
