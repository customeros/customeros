import { FC } from 'react';

import { Action } from '@graphql/types';

import { ContractStatusUpdatedAction } from './contract/ContractStatusUpdatedAction';

interface ActionStubProps {
  data: Action;
}

export const UserActionStub: FC<ActionStubProps> = ({ data }) => {
  if (data.actionType === 'CONTRACT_STATUS_UPDATED') {
    return <ContractStatusUpdatedAction data={data} />;
  }
  // This should be handled too as it currently appears in the timeline
  // if (data.actionType === 'CREATED') {
  //   return <p>CREATED</p>;
  // }

  return null;
};
