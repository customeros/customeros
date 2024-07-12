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
    console.log(
      '🏷️ ----- service.tempValue.closed: ',
      service.tempValue.closed,
    );
    const showEditView =
      (isDraft && !service.tempValue?.closed) ||
      (isFutureVersion && !service.tempValue?.closed) ||
      (!service?.tempValue?.metadata.id && !service?.tempValue?.closed) ||
      (!service?.tempValue?.closed &&
        service?.tempValue?.metadata?.id?.includes('new'));
    console.log('🏷️ ----- showEditView: ', showEditView);

    return (
      <>
        {showEditView ? (
          <ServiceItemEdit
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
