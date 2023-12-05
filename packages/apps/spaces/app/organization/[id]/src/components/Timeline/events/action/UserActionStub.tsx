import { FC } from 'react';

import { Action } from '@graphql/types';

import { ServiceUpdatedAction } from './service/ServiceUpdatedAction';
import { ContractStatusUpdatedAction } from './contract/ContractStatusUpdatedAction';

interface ActionStubProps {
  data: Action;
}

export const UserActionStub: FC<ActionStubProps> = ({ data }) => {
  if (data.actionType === 'CONTRACT_STATUS_UPDATED') {
    return <ContractStatusUpdatedAction data={data} />;
  }
  if (data.actionType === 'SERVICE_LINE_ITEM_QUANTITY_UPDATED') {
    return <ServiceUpdatedAction data={data} />;
  }
  // This should be handled too as it currently appears in the timeline
  // if (data.actionType === 'CREATED') {
  //   return <p>CREATED</p>;
  // }

  return null;
};
