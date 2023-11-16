'use client';
import { produce } from 'immer';
import { Command, CommandMenu, useCommands, CommandWrapper } from 'kmenu';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';
import {
  LastTouchpointType,
  RenewalLikelihoodProbability,
} from '@graphql/types';
import {
  useOwnerFilter,
  defaultState as ownerDefaultState,
} from '@organizations/components/Columns/Filters/Owner/OwnerFilter.atom';
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
  useTimeToRenewalFilter,
  defaultState as timeToRenewalDefaultState,
} from '@organizations/components/Columns/Filters/TimeToRenewal/TimeToRenewalFilter.atom';
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
  const client = getGraphQLClient();
  const [_owner, setOwnerFilter] = useOwnerFilter();
  const [_forecast, setForecastFilter] = useForecastFilter();
  const [_touchpoint, setTouchpointFilter] = useLastTouchpointFilter();
  const [_relationship, setRelationshipFilter] = useRelationshipFilter();
  const [_organization, setOrganizationFilter] = useOrganizationFilter();
  const [_timeToRenewal, setTimeToRenewalFilter] = useTimeToRenewalFilter();
  const [_renewal, setRenewalLikelihoodFilter] = useRenewalLikelihoodFilter();

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
                draft.value = [RenewalLikelihoodProbability.High];
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
                  RenewalLikelihoodProbability.Low,
                  RenewalLikelihoodProbability.Zero,
                ];
              }),
            );
          },
        },
        {
          text: 'who messaged recently',
          perform: () => {
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
  const { data } = useTenantNameQuery(client);

  if (data?.tenant !== 'openlineai') return null;

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
