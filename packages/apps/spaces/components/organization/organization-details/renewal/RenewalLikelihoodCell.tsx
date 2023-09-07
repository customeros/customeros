import { useState } from 'react';
import capitalize from 'lodash/capitalize';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';
import { RenewalLikelihoodProbability } from '@graphql/types';
import { useUpdateRenewalLikelihoodMutation } from '@spaces/graphql';

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
  const [isEditing, setIsEditing] = useState(false);
  const [updateRenewalLikelihood] = useUpdateRenewalLikelihoodMutation();

  const isIncreased = isLikelihoodIncreased(
    currentProbability,
    previousProbability,
  );
  const value = currentProbability
    ? { label: capitalize(currentProbability), value: currentProbability }
    : undefined;

  const handleChange = async (
    newValue: SelectOption<RenewalLikelihoodProbability>,
  ) => {
    setIsEditing(false);
    await updateRenewalLikelihood({
      variables: {
        input: {
          id: organizationId,
          probability: newValue.value,
        },
      },
      update: (cache) => {
        const normalizedId = cache.identify({
          id: organizationId,
          __typename: 'Organization',
        });

        cache.modify({
          id: normalizedId,
          fields: {
            accountDetails() {
              return {
                __typename: 'OrgAccountDetails',
                renewalLikelihood: {
                  __typename: 'RenewalLikelihood',
                  probability: newValue.value,
                  previousProbability: currentProbability,
                  updatedAt: new Date().toISOString(),
                },
              };
            },
          },
        });
        cache.gc();
      },
    });
  };

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
