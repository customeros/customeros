import { FC } from 'react';

import { Action } from '@graphql/types';

import { RenewalForecastUpdatedAction } from './renewal-forecast/RenewalForecastUpdatedAction';
import { RenewalLikelihoodUpdatedAction } from './renewal-likelihood/RenewalLikelihoodUpdatedAction';

interface ActionStubProps {
  data: Action;
}

export const UserActionStub: FC<ActionStubProps> = ({ data }) => {
  if (data.actionType === 'RENEWAL_FORECAST_UPDATED') {
    return <RenewalForecastUpdatedAction data={data} />;
  }
  if (data.actionType === 'RENEWAL_LIKELIHOOD_UPDATED') {
    return <RenewalLikelihoodUpdatedAction data={data} />;
  }
  // This should be handled too as it currently appears in the timeline
  // if (data.actionType === 'CREATED') {
  //   return <p>CREATED</p>;
  // }

  return null;
};
