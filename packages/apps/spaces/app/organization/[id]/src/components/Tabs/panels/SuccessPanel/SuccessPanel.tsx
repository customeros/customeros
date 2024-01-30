'use client';
import { useParams } from 'next/navigation';

import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';

import { PanelContainer } from './PanelContainer';
import { OnboardingMenu } from './OnboardingMenu';
import { OnboardingPlans } from './OnboardingPlans';
import { OnboardingStatus } from './OnboardingStatus';

export const SuccessPanel = () => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;
  const { data } = useOrganizationQuery(client, { id });
  const isFeatureOn = useFeatureIsOn('onboarding-plans');

  return (
    <PanelContainer
      title='Success'
      actionItem={isFeatureOn ? <OnboardingMenu /> : undefined}
    >
      <OnboardingStatus data={data?.organization?.accountDetails?.onboarding} />

      <OnboardingPlans />
    </PanelContainer>
  );
};
