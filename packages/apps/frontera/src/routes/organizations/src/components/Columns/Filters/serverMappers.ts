import { mapOwnerToAtom } from './Owner';
import { mapWebsiteToAtom } from './Website';
import { mapForecastToAtom } from './Forecast';
import { mapOnboardingToAtom } from './Onboarding';
import { mapRelationshipToAtom } from './Relationship';
import { mapOrganizationToAtom } from './Organization';
import { mapTimeToRenewalToAtom } from './TimeToRenewal';
import { mapLastTouchpointToAtom } from './LastTouchpoint';
import { mapRenewalLikelihoodToAtom } from './RenewalLikelihood';

const serverMappers = {
  OWNER_ID: mapOwnerToAtom,
  WEBSITE: mapWebsiteToAtom,
  NAME: mapOrganizationToAtom,
  FORECAST_ARR: mapForecastToAtom,
  IS_CUSTOMER: mapRelationshipToAtom,
  RENEWAl_DATE: mapTimeToRenewalToAtom,
  ONBOARDING_STATUS: mapOnboardingToAtom,
  LAST_TOUCHPOINT_AT: mapLastTouchpointToAtom,
  LAST_TOUCHPOINT_TYPE: mapLastTouchpointToAtom,
  RENEWAL_LIKELIHOOD: mapRenewalLikelihoodToAtom,
};

export const getServerToAtomMapper = (property: string) =>
  serverMappers[property as keyof typeof serverMappers];
