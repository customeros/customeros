import { FC } from 'react';

import { observer } from 'mobx-react-lite';

import { DateTimeUtils } from '@utils/date.ts';
import { ContractStatus } from '@graphql/types';

import { ServiceItemEdit } from './ServiceItemEdit.tsx';
import { ServiceItemPreview } from './ServiceItemPreview.tsx';
import ServiceLineItemStore from '../../../stores/Service.store.ts';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  usedDates?: string[];
  billingEnabled: boolean;
  isModification?: boolean;
  service: ServiceLineItemStore;
  allowIndividualRestore?: boolean;
  type: 'subscription' | 'one-time';
  allServices?: ServiceLineItemStore[];
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
      service?.serviceLineItem?.serviceStarted &&
      DateTimeUtils.isFuture(service?.serviceLineItem?.serviceStarted);

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const showEditView =
      (isDraft && !service.serviceLineItem?.isDeleted) ||
      (isFutureVersion && !service.serviceLineItem?.isDeleted) ||
      (service?.serviceLineItem?.isNew &&
        !service.serviceLineItem.isDeleted &&
        !service.serviceLineItem.closedVersion);

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
