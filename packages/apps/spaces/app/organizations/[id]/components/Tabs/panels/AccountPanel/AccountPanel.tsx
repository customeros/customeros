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
        renewalLikelihood={
          data?.organization?.accountDetails
            ?.renewalLikelihood as RenewalLikelihoodType
        }
      />
      <RenewalForecast
        renewalForecast={
          data?.organization?.accountDetails
            ?.renewalForecast as RenewalForecastType
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
