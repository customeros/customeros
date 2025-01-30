import { FC } from 'react';

import { observer } from 'mobx-react-lite';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store.ts';

import { DateTimeUtils } from '@utils/date.ts';
import { ContractStatus } from '@graphql/types';

import { ProductItemEdit } from './ProductItemEdit.tsx';
import { ProductItemPreview } from './ProductItemPreview.tsx';

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

export const ProductItem: FC<ServiceItemProps> = observer(
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
          <ProductItemEdit
            type={type}
            service={service}
            currency={currency}
            allServices={allServices}
            isModification={isModification}
            contractStatus={contractStatus}
          />
        ) : (
          <ProductItemPreview
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
