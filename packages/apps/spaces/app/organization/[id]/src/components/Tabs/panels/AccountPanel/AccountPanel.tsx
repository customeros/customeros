'use client';

import { RenewalLikelihood, RenewalLikelihoodType } from './RenewalLikelihood';
import { useParams } from 'next/navigation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';
import { RenewalForecast, RenewalForecastType } from './RenewalForecast';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import { BillingDetailsCard } from './BillingDetailsCard/BillingDetailsCard';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import { TimeToRenewal } from './TimeToRenewal';
import { Notes } from './Notes';
import { useDisclosure } from '@ui/utils';
import { useMemo } from 'react';

export const AccountPanel = () => {
  const id = useParams()?.id as string;
  // Moved to upperscope due to error in safari https://linear.app/customer-os/issue/COS-619/scrollbar-overlaps-the-renewal-modals-in-safari
  const renewalLikelihoodUpdateModal = useDisclosure({
    id: 'renewal-likelihood-update-modal',
  });
  const renewalLikelihoodInfoModal = useDisclosure({
    id: 'renewal-likelihood-info-modal',
  });

  const renewalForecastUpdateModal = useDisclosure({
    id: 'renewal-renewal-update-modal',
  });
  const renewalForecastInfoModal = useDisclosure({
    id: 'renewal-renewal-info-modal',
  });

  const client = getGraphQLClient();
  const { data, isInitialLoading } = useOrganizationAccountDetailsQuery(
    client,
    { id },
  );
  const isModalOpen = useMemo(() => {
    return (
      renewalForecastUpdateModal.isOpen ||
      renewalLikelihoodUpdateModal.isOpen ||
      renewalForecastInfoModal.isOpen ||
      renewalLikelihoodInfoModal.isOpen
    );
  }, [
    renewalForecastUpdateModal.isOpen,
    renewalLikelihoodUpdateModal.isOpen,
    renewalForecastInfoModal.isOpen,
    renewalLikelihoodInfoModal.isOpen,
  ]);

  if (isInitialLoading) {
    return <AccountPanelSkeleton />;
  }

  return (
    <OrganizationPanel
      title='Account'
      withFade
      shouldBlockPanelScroll={isModalOpen}
    >
      <RenewalLikelihood
        infoModal={renewalLikelihoodInfoModal}
        updateModal={renewalLikelihoodUpdateModal}
        name={data?.organization?.name || ''}
        data={
          data?.organization?.accountDetails
            ?.renewalLikelihood as RenewalLikelihoodType
        }
      />
      <RenewalForecast
        infoModal={renewalForecastInfoModal}
        updateModal={renewalForecastUpdateModal}
        isInitialLoading={isInitialLoading}
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
      <Notes id={id} data={data?.organization} />
    </OrganizationPanel>
  );
};
