'use client';

import { RenewalLikelihood, RenewalLikelihoodType } from './RenewalLikelihood';
import { useParams } from 'next/navigation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';
import { RenewalForecast, RenewalForecastType } from './RenewalForecast';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import {
  BillingDetailsCard,
  BillingDetailsType,
} from './BillingDetailsCard/BillingDetailsCard';
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
    <OrganizationPanel title='Account'>
      <RenewalLikelihood
        name={data?.organization?.name || ''}
        renewalLikelihood={
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
        renewalCycle={
          data?.organization?.accountDetails?.billingDetails?.renewalCycle
        }
        renewalCycleStart={
          data?.organization?.accountDetails?.billingDetails?.renewalCycleStart
        }
      />
      <BillingDetailsCard
        id={id}
        billingDetailsData={
          data?.organization?.accountDetails
            ?.billingDetails as BillingDetailsType
        }
      />
    </OrganizationPanel>
  );
};
