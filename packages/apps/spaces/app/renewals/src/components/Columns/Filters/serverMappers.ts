import { mapOwnerToAtom } from './Owner';
import { mapForecastToAtom } from './Forecast';
import { mapOrganizationToAtom } from './Organization';
import { mapTimeToRenewalToAtom } from './TimeToRenewal';
import { mapLastTouchpointToAtom } from './LastTouchpoint';
import { mapRenewalLikelihoodToAtom } from './RenewalLikelihood';

const serverMappers = {
  OWNER_ID: mapOwnerToAtom,
  NAME: mapOrganizationToAtom,
  FORECAST_ARR: mapForecastToAtom,
  RENEWAl_DATE: mapTimeToRenewalToAtom,
  LAST_TOUCHPOINT_AT: mapLastTouchpointToAtom,
  LAST_TOUCHPOINT_TYPE: mapLastTouchpointToAtom,
  RENEWAL_LIKELIHOOD: mapRenewalLikelihoodToAtom,
};

export const getServerToAtomMapper = (property: string) =>
  serverMappers[property as keyof typeof serverMappers];
