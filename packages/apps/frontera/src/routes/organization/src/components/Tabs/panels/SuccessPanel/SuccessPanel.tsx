'use client';
import { useParams } from 'react-router-dom';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationQuery } from '@organization/graphql/organization.generated';

import { PanelContainer } from './PanelContainer';
import { OnboardingStatus } from './OnboardingStatus';

export const SuccessPanel = () => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;
  const { data, isPending } = useOrganizationQuery(client, { id });

  return (
    <PanelContainer title='Success'>
      <OnboardingStatus
        isLoading={isPending}
        data={data?.organization?.accountDetails?.onboarding}
      />
    </PanelContainer>
  );
};
