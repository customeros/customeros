'use client';
import React, { FC } from 'react';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { ComparisonOperator, Contact } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { GetContactsEmailListDocument } from '@organization/graphql/getContactsEmailList.generated';
import { emailRegex } from '@organization/components/Timeline/events/email/utils';
import { OptionsOrGroups } from 'react-select';
import { EmailFormMultiCreatableSelect } from '@organization/components/Timeline/events/email/compose-email/EmailFormMultiCreatableSelect';

interface EmailParticipantSelect {
  entryType: string;
  fieldName: string;
  formId: string;
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
    callback: (options: OptionsOrGroups<any, any>) => void,
  ) => {
    try {
      const results = await client.request<any>(GetContactsEmailListDocument, {
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
      const options: OptionsOrGroups<string, any> = (
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
      maxH='86px'
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
