'use client';

import { RenewalLikelihood } from './RenewalLikelihood';
import { useParams } from 'next/navigation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';
import { RenewalForecast } from './RenewalForecast';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import { BillingDetailsCard } from './BillingDetailsCard/BillingDetailsCard';
import { AccountPanelSkeleton } from '@organization/components/Tabs/panels/AccountPanel/AccountPanelSkeleton';
import { BillingDetails as BillingDetailsType } from '@graphql/types';
export const AccountPanel = () => {
  const id = useParams()?.id as string;

  const client = getGraphQLClient();
  const { data, isInitialLoading } = useOrganizationAccountDetailsQuery(
    client,
    { id },
  );

  if (isInitialLoading) {
    return <AccountPanelSkeleton />;
  }

  return (
    <OrganizationPanel title='Account'>
      <RenewalLikelihood
      // likelyhoodData={data?.organization?.accountDetails?.renewalLikelihood}
      />
      <RenewalForecast
      // forecastData={data?.organization?.accountDetails?.renewalForecast}
      />
      <BillingDetailsCard
        id={id}
        billingDetailsData={
          data?.organization?.accountDetails
            ?.billingDetails as BillingDetailsType & { amount: string }
        }
      />
    </OrganizationPanel>
  );
};
