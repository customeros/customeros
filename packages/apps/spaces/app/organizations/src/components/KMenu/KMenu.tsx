'use client';
import { produce } from 'immer';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { Command, CommandMenu, useCommands, CommandWrapper } from 'kmenu';

import {
  LastTouchpointType,
  OpportunityRenewalLikelihood,
} from '@graphql/types';
import {
  useOwnerFilter,
  defaultState as ownerDefaultState,
} from '@organizations/components/Columns/Filters/Owner/OwnerFilter.atom';
import {
  useTimeToRenewalFilter,
  defaultState as timeToRenewalDefaultState,
} from '@organizations/components/Columns/Filters/TimeToRenewal';
import {
  useForecastFilter,
  defaultState as forecastDefaultState,
} from '@organizations/components/Columns/Filters/Forecast/ForecastFilter.atom';
import {
  useRelationshipFilter,
  defaultState as relationshipDefaultState,
} from '@organizations/components/Columns/Filters/Relationship/RelationshipFilter.atom';
import {
  useOrganizationFilter,
  defaultState as organizationDefaultState,
} from '@organizations/components/Columns/Filters/Organization/OrganizationFilter.atom';
import {
  useLastTouchpointFilter,
  defaultState as touchpointDefaultState,
} from '@organizations/components/Columns/Filters/LastTouchpoint/LastTouchpointFilter.atom';
import {
  useRenewalLikelihoodFilter,
  defaultState as renewalLikelihoodDefaultState,
} from '@organizations/components/Columns/Filters/RenewalLikelihood/RenewalLikelihoodFilter.atom';

import 'kmenu/dist/index.css';

export const KMenu = () => {
  const [_owner, setOwnerFilter] = useOwnerFilter();
  const [_forecast, setForecastFilter] = useForecastFilter();
  const [_touchpoint, setTouchpointFilter] = useLastTouchpointFilter();
  const [_relationship, setRelationshipFilter] = useRelationshipFilter();
  const [_organization, setOrganizationFilter] = useOrganizationFilter();
  const [_timeToRenewal, setTimeToRenewalFilter] = useTimeToRenewalFilter();
  const [_renewal, setRenewalLikelihoodFilter] = useRenewalLikelihoodFilter();

  const isFeatureOn = useFeatureIsOn('kmenu');

  const resetFilters = () => {
    setOwnerFilter(ownerDefaultState);
    setForecastFilter(forecastDefaultState);
    setTouchpointFilter(touchpointDefaultState);
    setOrganizationFilter(organizationDefaultState);
    setRelationshipFilter(relationshipDefaultState);
    setTimeToRenewalFilter(timeToRenewalDefaultState);
    setRenewalLikelihoodFilter(renewalLikelihoodDefaultState);
  };

  const main: Command[] = [
    {
      category: 'Filter by',
      commands: [
        {
          text: 'who is likely to renew',
          perform: () => {
            resetFilters();
            setRenewalLikelihoodFilter((prev) =>
              produce(prev, (draft) => {
                draft.isActive = true;
                draft.value = [OpportunityRenewalLikelihood.HighRenewal];
              }),
            );
          },
        },
        {
          text: 'who is unlikely to renew',
          perform: () => {
            resetFilters();
            setRenewalLikelihoodFilter((prev) =>
              produce(prev, (draft) => {
                draft.isActive = true;
                draft.value = [
                  OpportunityRenewalLikelihood.LowRenewal,
                  OpportunityRenewalLikelihood.ZeroRenewal,
                ];
              }),
            );
          },
        },
        {
          text: 'who messaged recently',
          perform: () => {
            resetFilters();
            setTouchpointFilter((prev) =>
              produce(prev, (draft) => {
                draft.isActive = true;
                draft.value = [LastTouchpointType.InteractionEventChat];
              }),
            );
          },
        },
        {
          text: 'who did we met recently',
          perform: () => {
            resetFilters();
            setTouchpointFilter((prev) =>
              produce(prev, (draft) => {
                draft.isActive = true;
                draft.value = [LastTouchpointType.Meeting];
              }),
            );
          },
        },
        {
          text: 'who emailed recently',
          perform: () => {
            resetFilters();
            setTouchpointFilter((prev) =>
              produce(prev, (draft) => {
                draft.isActive = true;
                draft.value = [LastTouchpointType.InteractionEventEmailSent];
              }),
            );
          },
        },
        {
          text: 'who was created recently',
          perform: () => {
            resetFilters();
            setTouchpointFilter((prev) =>
              produce(prev, (draft) => {
                draft.isActive = true;
                draft.value = [LastTouchpointType.ActionCreated];
              }),
            );
          },
        },
      ],
    },
    {
      category: 'Actions',
      commands: [
        {
          text: 'Clear all filters',
          perform: resetFilters,
        },
      ],
    },
  ];

  const [mainCommands] = useCommands(main);

  if (!isFeatureOn) return null;

  return (
    <CommandWrapper>
      <CommandMenu
        commands={mainCommands}
        crumbs={['Organizations']}
        index={1}
        placeholder='What would you like to do?'
      />
    </CommandWrapper>
  );
};
