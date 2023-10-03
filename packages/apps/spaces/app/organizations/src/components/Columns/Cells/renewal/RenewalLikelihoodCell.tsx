import { useRef, useState, useEffect } from 'react';
import { produce } from 'immer';
import capitalize from 'lodash/capitalize';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import { useQueryClient, InfiniteData } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';
import { RenewalLikelihoodProbability } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  useInfiniteGetOrganizationsQuery,
  GetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';
import { useOrganizationsMeta } from '@organizations/shared/state';

import { useUpdateRenewalLikelihoodMutation } from '@organization/src/graphql/updateRenewalLikelyhood.generated';

import {
  getLikelihoodColor,
  isLikelihoodIncreased,
  renewalLikelihoodOptions,
} from './utils';

interface RenewalLikelihoodCellProps {
  updatedAt: string | null;
  organizationId: string;
  currentProbability?: RenewalLikelihoodProbability | null;
  previousProbability?: RenewalLikelihoodProbability | null;
}

export const RenewalLikelihoodCell = ({
  updatedAt,
  organizationId,
  currentProbability,
  previousProbability,
}: RenewalLikelihoodCellProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [organizationsMeta] = useOrganizationsMeta();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isEditing, setIsEditing] = useState(false);

  const { getOrganization } = organizationsMeta;
  const queryKey = useInfiniteGetOrganizationsQuery.getKey(getOrganization);

  const updateRenewalLikelihood = useUpdateRenewalLikelihoodMutation(client, {
    onMutate: (payload) => {
      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.cancelQueries(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            const row = draft?.pages[
              pageIndex
            ]?.dashboardView_Organizations?.content.find(
              (c) => c.id === organizationId,
            );
            const likelihood = row?.accountDetails?.renewalLikelihood;

            if (likelihood) {
              likelihood.probability = payload.input.probability;
              likelihood.previousProbability = currentProbability;
            }
          });
        },
      );
      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        context?.previousEntries,
      );
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
      }, 500);
    },
  });

  const isIncreased = isLikelihoodIncreased(
    currentProbability,
    previousProbability,
  );
  const value = currentProbability
    ? { label: capitalize(currentProbability), value: currentProbability }
    : undefined;

  const handleChange = (
    newValue: SelectOption<RenewalLikelihoodProbability>,
  ) => {
    setIsEditing(false);
    updateRenewalLikelihood.mutate({
      input: {
        id: organizationId,
        probability: newValue.value,
      },
    });
  };

  useEffect(() => {
    return () => {
      timeoutRef?.current && clearTimeout(timeoutRef.current);
    };
  }, []);

  return (
    <Flex flexDir='column' key={Math.random()}>
      {isEditing ? (
        <Select
          size='sm'
          variant='unstyled'
          placeholder='Not set'
          value={value}
          autoFocus
          onKeyDown={(e) => {
            if (e.key === 'Escape') {
              setIsEditing(false);
            }
          }}
          defaultMenuIsOpen
          openMenuOnClick={false}
          onBlur={() => setIsEditing(false)}
          onChange={handleChange}
          leftElement={<Flex w='3' h='3' />}
          options={renewalLikelihoodOptions}
          chakraStyles={{
            singleValue: (props) => ({
              ...props,
              color: getLikelihoodColor(currentProbability),
              paddingBottom: 0,
            }),
            control: (props) => ({
              ...props,
              minH: '0',
            }),
            placeholder: (props) => ({
              ...props,
              color: 'gray.400',
            }),
            valueContainer: (props) => ({
              ...props,
              ml: 1.5,
            }),
            inputContainer: (props) => ({
              ...props,
              paddingTop: 0,
              paddingBottom: 0,
            }),
          }}
        />
      ) : (
        <Flex
          w='full'
          gap='1'
          ml='5'
          align='center'
          _hover={{
            '& #edit-button': {
              opacity: 1,
            },
          }}
        >
          <Text
            cursor='default'
            color={value ? getLikelihoodColor(currentProbability) : 'gray.400'}
            onDoubleClick={() => setIsEditing(true)}
          >
            {value?.label ?? 'Not set'}
          </Text>
          <IconButton
            aria-label='erc'
            size='xs'
            borderRadius='md'
            minW='4'
            w='4'
            minH='4'
            h='4'
            opacity='0'
            variant='ghost'
            id='edit-button'
            onClick={() => setIsEditing(true)}
            icon={<Icons.Edit3 color='gray.500' boxSize='3' />}
          />
        </Flex>
      )}
      {currentProbability && (
        <Flex align='center'>
          {!previousProbability ? (
            <Icons.Dot boxSize='3' color='gray.500' />
          ) : isIncreased ? (
            <Icons.ArrowNarrowUpRight boxSize='3' color='gray.500' />
          ) : (
            <Icons.ArrowNarrowDownRight boxSize='3' color='gray.500' />
          )}
          <Text color='gray.500' ml='2'>
            {updatedAt
              ? `${formatDistanceToNow(new Date(updatedAt))} ago`
              : 'Not set'}
          </Text>
        </Flex>
      )}
    </Flex>
  );
};
