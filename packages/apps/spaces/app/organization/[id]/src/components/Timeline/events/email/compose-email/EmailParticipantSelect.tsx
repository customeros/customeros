'use client';
import React, { FC, useCallback } from 'react';
import { OptionsOrGroups } from 'react-select';

import { GroupBase, OptionProps, MultiValueProps } from 'chakra-react-select';
import {
  Menu,
  MenuItem,
  MenuButton,
  MenuList as ChakraMenuList,
} from '@chakra-ui/menu';

import { Flex } from '@ui/layout/Flex';
import { SelectOption } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton';
import { chakraComponents } from '@ui/form/SyncSelect';
import { Contact, ComparisonOperator } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { emailRegex } from '@organization/src/components/Timeline/events/email/utils';
import { GetContactsEmailListDocument } from '@organization/src/graphql/getContactsEmailList.generated';
import { EmailFormMultiCreatableSelect } from '@organization/src/components/Timeline/events/email/compose-email/EmailFormMultiCreatableSelect';
interface EmailParticipantSelect {
  formId: string;
  entryType: string;
  fieldName: string;
  autofocus: boolean;
}

export const EmailParticipantSelect: FC<EmailParticipantSelect> = ({
  entryType,
  fieldName,
  formId,
  autofocus = false,
}) => {
  const client = getGraphQLClient();

  const getFilteredSuggestions = async (
    filterString: string,
    callback: (options: OptionsOrGroups<unknown, GroupBase<unknown>>) => void,
  ) => {
    try {
      const results = await client.request<{
        contacts: { content: Contact[] };
      }>(GetContactsEmailListDocument, {
        pagination: {
          page: 1,
          limit: 5,
        },
        where: {
          OR: [
            {
              filter: {
                property: 'FIRST_NAME',
                value: filterString,
                operation: ComparisonOperator.Contains,
              },
            },
            {
              filter: {
                property: 'LAST_NAME',
                value: filterString,
                operation: ComparisonOperator.Contains,
              },
            },
            {
              filter: {
                property: 'NAME',
                value: filterString,
                operation: ComparisonOperator.Contains,
              },
            },
          ],
        },
      });

      const options: OptionsOrGroups<unknown, GroupBase<unknown>> = (
        results?.contacts?.content || []
      )
        .filter((e: Contact) => e.emails.length)
        .map((e: Contact) =>
          e.emails.map((email) => ({
            value: email.email,
            label: `${e.firstName} ${e.lastName}`,
          })),
        )
        .flat();

      callback(options);
    } catch (error) {
      callback([]);
    }
  };

  return (
    <Flex
      alignItems='baseline'
      marginBottom={-1}
      marginTop={0}
      flex={1}
      overflow='visible'
      // maxH='86px'
    >
      <Text as={'span'} color='gray.700' fontWeight={600} mr={1}>
        {entryType}:
      </Text>
      <EmailFormMultiCreatableSelect
        autoFocus={autofocus}
        name={fieldName}
        formId={formId}
        placeholder='Enter name or email...'
        noOptionsMessage={() => null}
        loadOptions={(inputValue: string, callback) => {
          getFilteredSuggestions(inputValue, callback);
        }}
        allowCreateWhileLoading={false}
        formatCreateLabel={(input) => {
          return input;
        }}
        isValidNewOption={(input) => emailRegex.test(input)}
        getOptionLabel={(d) => {
          if (d?.__isNew__) {
            return `${d.label}`;
          }

          return `${d.label} - ${d.value}`;
        }}
      />
    </Flex>
  );
};
