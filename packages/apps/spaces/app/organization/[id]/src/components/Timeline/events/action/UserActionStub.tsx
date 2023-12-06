import { FC } from 'react';

import { Action, ActionType } from '@graphql/types';

import { ServiceUpdatedAction } from './service/ServiceUpdatedAction';
import { ContractStatusUpdatedAction } from './contract/ContractStatusUpdatedAction';

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
      return <ContractStatusUpdatedAction data={data} />;
    case ActionType.ServiceLineItemQuantityUpdated:
    case ActionType.ServiceLineItemPriceUpdated:
    case ActionType.ServiceLineItemBilledTypeUpdated:
      return <ServiceUpdatedAction data={data} />;
    case ActionType.ServiceLineItemBilledTypeOnceCreated:
    case ActionType.ServiceLineItemBilledTypeUsageCreated:
    case ActionType.ServiceLineItemBilledTypeRecurringCreated:
      return <ServiceUpdatedAction data={data} mode='created' />;
  }

  return null;
};
