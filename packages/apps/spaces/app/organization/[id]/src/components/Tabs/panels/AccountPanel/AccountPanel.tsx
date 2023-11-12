'use client';

// import { useMemo } from 'react';
import { useParams } from 'next/navigation';

// import { useDisclosure } from '@ui/utils';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ContractCard } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ContractCard';
import { ARRForecast } from '@organization/src/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast';
import { useOrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';

import { Notes } from './Notes';
import { RenewalForecastType } from './RenewalForecast';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';

export const AccountPanel = () => {
  const id = useParams()?.id as string;
  // Moved to upperscope due to error in safari https://linear.app/customer-os/issue/COS-619/scrollbar-overlaps-the-renewal-modals-in-safari

  // Todo modify and connect modals
  // const renewalLikelihoodUpdateModal = useDisclosure({
  //   id: 'renewal-likelihood-update-modal',
  // });
  // const renewalLikelihoodInfoModal = useDisclosure({
  //   id: 'renewal-likelihood-info-modal',
  // });
  //
  // const renewalForecastUpdateModal = useDisclosure({
  //   id: 'renewal-renewal-update-modal',
  // });
  // const renewalForecastInfoModal = useDisclosure({
  //   id: 'renewal-renewal-info-modal',
  // });

  const client = getGraphQLClient();
  const { data, isInitialLoading } = useOrganizationAccountDetailsQuery(
    client,
    { id },
  );
  // const isModalOpen = useMemo(() => {
  //   return (
  //     renewalForecastUpdateModal.isOpen ||
  //     renewalLikelihoodUpdateModal.isOpen ||
  //     renewalForecastInfoModal.isOpen ||
  //     renewalLikelihoodInfoModal.isOpen
  //   );
  // }, [
  //   renewalForecastUpdateModal.isOpen,
  //   renewalLikelihoodUpdateModal.isOpen,
  //   renewalForecastInfoModal.isOpen,
  //   renewalLikelihoodInfoModal.isOpen,
  // ]);

  if (isInitialLoading) {
    return <AccountPanelSkeleton />;
  }

  return (
    <OrganizationPanel
      title='Account'
      withFade
      // shouldBlockPanelScroll={isModalOpen}
    >
      <ARRForecast
        name={data?.organization?.name || ''}
        isInitialLoading={isInitialLoading}
        renewalProbability={
          data?.organization?.accountDetails?.renewalLikelihood?.probability
        }
        aRRForecast={
          data?.organization?.accountDetails
            ?.renewalForecast as RenewalForecastType
        }
      />

      <ContractCard name={data?.organization?.name || ''} data={null} />

      <Notes id={id} data={data?.organization} />
    </OrganizationPanel>
  );
};
