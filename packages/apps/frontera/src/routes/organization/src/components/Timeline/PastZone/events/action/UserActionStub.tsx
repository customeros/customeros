import { FC } from 'react';

import { Action, ActionType } from '@graphql/types';

import { ServiceUpdatedAction } from './service/ServiceUpdatedAction';
import InvoiceStatusChangeAction from './invoice/InvoiceStatusChangeAction';
import { ContractStatusUpdatedAction } from './contract/ContractStatusUpdatedAction';
import { OnboardingStatusChangedAction } from './onboarding/OnboardingStatusChangedAction';

interface ActionStubProps {
  data: Action;
}

export const UserActionStub: FC<ActionStubProps> = ({ data }) => {
  // This should be handled too as it currently appears in the timeline
  // if (data.actionType === 'CREATED') {
  //   return <p>CREATED</p>;
  // }

  switch (data.actionType) {
    case ActionType.ContractStatusUpdated:
    case ActionType.ContractRenewed:
      return <ContractStatusUpdatedAction data={data} />;
    case ActionType.ServiceLineItemQuantityUpdated:
    case ActionType.ServiceLineItemPriceUpdated:
    case ActionType.ServiceLineItemBilledTypeUpdated:
      return <ServiceUpdatedAction data={data} />;
    case ActionType.ServiceLineItemBilledTypeOnceCreated:
    case ActionType.ServiceLineItemBilledTypeUsageCreated:
    case ActionType.ServiceLineItemBilledTypeRecurringCreated:
      return <ServiceUpdatedAction data={data} mode='created' />;
    case ActionType.ServiceLineItemRemoved:
      return <ServiceUpdatedAction data={data} mode='removed' />;
    case ActionType.OnboardingStatusChanged:
      return <OnboardingStatusChangedAction data={data} />;
    case ActionType.InvoiceIssued:
    case ActionType.InvoicePaid:
    case ActionType.InvoiceSent:
    case ActionType.InvoiceVoided:
    case ActionType.InvoiceOverdue:
      return <InvoiceStatusChangeAction data={data} mode={data.actionType} />;
  }

  return null;
};
