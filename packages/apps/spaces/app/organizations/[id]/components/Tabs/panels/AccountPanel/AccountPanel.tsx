'use client';
import { useState } from 'react';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { BillingDetailsCard } from '@organization/components/Tabs/panels/AccountPanel/BillingDetailsCard/BillingDetailsCard';
import {
  RenewalLikelihood,
  Value as RenewalLikelihoodValue,
} from './RenewalLikelihood';
import { useParams } from 'next/navigation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useQueryClient } from '@tanstack/react-query';
import { useOrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';
import {
    RenewalForecast,
    Value as RenewalForecastValue,
} from './RenewalForecast';
export const AccountPanel = () => {
  const id = useParams()?.id as string;

  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data } = useOrganizationAccountDetailsQuery(client, { id });
    console.log('🏷️ ----- data: '
        , data);
    const invalidateQuery = () =>
    queryClient.invalidateQueries(
      useOrganizationAccountDetailsQuery.getKey({ id }),
    );

  // TODO move set state to renewal likelyhood component if those values are not needed at this level
  const [renewalLikelihood, setRenewalLikelihood] =
    useState<RenewalLikelihoodValue>({ reason: '', likelihood: 'NOT_SET' });
  const [renewalForecast, setRenewalForecast] = useState<RenewalForecastValue>({
    reason: '',
    forecast: '',
  });

  return (
    <OrganizationPanel title='Account'>
      <RenewalLikelihood
        value={renewalLikelihood}
        onChange={setRenewalLikelihood}
      />
      <RenewalForecast value={renewalForecast} onChange={setRenewalForecast} />
      <BillingDetailsCard />
    </OrganizationPanel>
  );
};
