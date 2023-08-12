'use client';

import { RenewalLikelihood } from './RenewalLikelihood';
import { useParams } from 'next/navigation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useQueryClient } from '@tanstack/react-query';
import { useOrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';
import { RenewalForecast } from './RenewalForecast';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import { BillingDetailsCard } from './BillingDetailsCard/BillingDetailsCard';
export const AccountPanel = () => {
  const id = useParams()?.id as string;

  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data } = useOrganizationAccountDetailsQuery(client, { id });
  const invalidateQuery = () =>
    queryClient.invalidateQueries(
      useOrganizationAccountDetailsQuery.getKey({ id }),
    );

  return (
    <OrganizationPanel title='Account'>
      <RenewalLikelihood
      // likelyhoodData={data?.organization?.accountDetails?.renewalLikelihood}
      />
      <RenewalForecast
      // forecastData={data?.organization?.accountDetails?.renewalForecast}
      />
      <BillingDetailsCard
        billingDetailsData={data?.organization?.accountDetails?.billingDetails}
      />
    </OrganizationPanel>
  );
};
