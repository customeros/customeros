import { FC } from 'react';

import { observer } from 'mobx-react-lite';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store.ts';

import { DateTimeUtils } from '@utils/date.ts';
import { ContractStatus } from '@graphql/types';

import { ServiceItemEdit } from './ServiceItemEdit.tsx';
import { ServiceItemPreview } from './ServiceItemPreview.tsx';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  dataTest?: string;
  usedDates?: string[];
  isModification?: boolean;
  service: ContractLineItemStore;
  allowIndividualRestore?: boolean;
  type: 'subscription' | 'one-time';
  allServices?: ContractLineItemStore[];
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
  }) => {
    const isFutureVersion =
      service?.tempValue?.serviceStarted &&
      DateTimeUtils.isFuture(service?.tempValue?.serviceStarted);

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const showEditView =
      (isDraft && !service.tempValue?.closed) ||
      (isFutureVersion && !service.tempValue?.closed) ||
      (!service?.tempValue?.metadata.id && !service?.tempValue?.closed) ||
      (!service?.tempValue?.closed &&
        service?.tempValue?.metadata?.id?.includes('new'));

    return (
      <>
        {showEditView ? (
          <ServiceItemEdit
            type={type}
            service={service}
            currency={currency}
            allServices={allServices}
            isModification={isModification}
            contractStatus={contractStatus}
          />
        ) : (
          <ServiceItemPreview
            type={type}
            service={service}
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
