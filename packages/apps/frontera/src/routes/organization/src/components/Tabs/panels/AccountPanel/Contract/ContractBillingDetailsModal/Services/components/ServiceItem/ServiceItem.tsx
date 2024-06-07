import { FC } from 'react';

import { Store } from '@store/store.ts';
import { observer } from 'mobx-react-lite';

import { DateTimeUtils } from '@utils/date.ts';
import { ContractStatus, ServiceLineItem } from '@graphql/types';

import { ServiceItemEdit } from './ServiceItemEdit.tsx';
import { ServiceItemPreview } from './ServiceItemPreview.tsx';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  usedDates?: string[];
  billingEnabled: boolean;
  isModification?: boolean;
  service: Store<ServiceLineItem>;
  allowIndividualRestore?: boolean;
  type: 'subscription' | 'one-time';
  allServices?: Store<ServiceLineItem>[];
  contractStatus?: ContractStatus | null;
}

export const ServiceItem: FC<ServiceItemProps> = observer(
  ({
    service,
    allServices,
    isEnded,
    currency,
    isModification,
    type,
    contractStatus,
    allowIndividualRestore,
    billingEnabled,
  }) => {
    const isFutureVersion =
      service?.value?.serviceStarted &&
      DateTimeUtils.isFuture(service?.value?.serviceStarted);

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const showEditView =
      (isDraft && !service.value?.closed) ||
      (isFutureVersion && !service.value?.closed) ||
      (!service?.value?.metadata.id && !service?.value?.closed);

    return (
      <>
        {showEditView ? (
          <ServiceItemEdit
            billingEnabled={billingEnabled}
            service={service}
            type={type}
            allServices={allServices}
            isModification={isModification}
            currency={currency}
            contractStatus={contractStatus}
          />
        ) : (
          <ServiceItemPreview
            service={service}
            type={type}
            isEnded={isEnded}
            currency={currency}
            contractStatus={contractStatus}
            allowIndividualRestore={allowIndividualRestore}
          />
        )}
      </>
    );
  },
);
