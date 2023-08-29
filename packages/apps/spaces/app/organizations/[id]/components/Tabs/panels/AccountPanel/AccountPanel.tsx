'use client';

import { RenewalLikelihood, RenewalLikelihoodType } from './RenewalLikelihood';
import { useParams } from 'next/navigation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';
import { RenewalForecast, RenewalForecastType } from './RenewalForecast';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import { BillingDetailsCard } from './BillingDetailsCard/BillingDetailsCard';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import { TimeToRenewal } from '@organization/components/Tabs/panels/AccountPanel/TimeToRenewal/TimeToRenewal';

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
    <OrganizationPanel title='Account' withFade>
      <RenewalLikelihood
        name={data?.organization?.name || ''}
        data={
          data?.organization?.accountDetails
            ?.renewalLikelihood as RenewalLikelihoodType
        }
      />
      <RenewalForecast
        name={data?.organization?.name || ''}
        renewalProbability={
          data?.organization?.accountDetails?.renewalLikelihood?.probability
        }
        renewalForecast={
          data?.organization?.accountDetails
            ?.renewalForecast as RenewalForecastType
        }
      />
      <TimeToRenewal
        id={data?.organization?.id || ''}
        data={data?.organization?.accountDetails?.billingDetails}
      />
      <BillingDetailsCard
        id={id}
        data={data?.organization?.accountDetails?.billingDetails}
      />
    </OrganizationPanel>
  );
};
