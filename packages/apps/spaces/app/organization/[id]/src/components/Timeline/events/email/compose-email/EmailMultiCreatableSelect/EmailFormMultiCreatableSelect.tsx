import { useParams } from 'next/navigation';
import { OptionsOrGroups } from 'react-select';
import { useField } from 'react-inverted-form';
import React, {
  useMemo,
  useState,
  useEffect,
  forwardRef,
  useCallback,
} from 'react';

import { GroupBase, OptionProps, MultiValueProps } from 'chakra-react-select';

import { SelectOption } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton';
import { chakraComponents } from '@ui/form/SyncSelect';
import { getName } from '@spaces/utils/getParticipantsName';
import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { Contact, ComparisonOperator } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { multiCreatableSelectStyles } from '@ui/form/MultiCreatableSelect/styles';
import { emailRegex } from '@organization/src/components/Timeline/events/email/utils';
import {
  FormSelectProps,
  MultiCreatableSelect,
} from '@ui/form/MultiCreatableSelect';
import {
  GetContactsEmailListDocument,
  useGetContactsEmailListQuery,
} from '@organization/src/graphql/getContactsEmailList.generated';
import { MultiValueWithActionMenu } from '@organization/src/components/Timeline/events/email/compose-email/EmailMultiCreatableSelect/MultiValueWithActionMenu';

type ExistingContact = { id: string; label: string; value?: string | null };
export const EmailFormMultiCreatableSelect = forwardRef<
  SelectInstance,
  FormSelectProps
>(({ name, formId, ...rest }, ref) => {
  const client = getGraphQLClient();
  const organizationId = useParams()?.id as string;
  const [existingContacts, setExistingContacts] = useState<
    Array<ExistingContact>
  >([]);

  const { data } = useGetContactsEmailListQuery(client, {
    id: organizationId,
    pagination: {
      page: 1,
      limit: 100,
    },
  });

  useEffect(() => {
    if (data?.organization?.contacts?.content?.length) {
      const organizationContacts = (
        (data?.organization?.contacts?.content || []) as Array<Contact>
      )
        .map((e: Contact) => {
          if (e.emails.some((e) => !!e.email)) {
            return e.emails.map((email) => ({
              id: e.id,
              value: email.email,
              label: `${e.firstName} ${e.lastName}`,
            }));
          }

          return [
            {
              id: e.id,
              label: getName(e),
              value: '',
            },
          ];
        })
        .flat();
      setExistingContacts(organizationContacts);
    }
  }, [data]);

  const { getInputProps } = useField(name, formId);
  const { id, onChange, onBlur, value } = getInputProps();
  const [_, copyToClipboard] = useCopyToClipboard();

  const handleBlur = (stringVal: string) => {
    if (stringVal && emailRegex.test(stringVal)) {
      onBlur([...value, { label: stringVal, value: stringVal }]);

      return;
    }
    onBlur(value);
  };

  const getFilteredSuggestions = async (
    filterString: string,
    callback: (options: OptionsOrGroups<unknown, GroupBase<unknown>>) => void,
  ) => {
    try {
      const results = await client.request<{
        organization: {
          contacts: { content: Contact[] };
        };
      }>(GetContactsEmailListDocument, {
        id: organizationId,
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
        results?.organization?.contacts?.content || []
      )
        .map((e: Contact) => {
          if (e.emails.some((e) => !!e.email)) {
            return e.emails.map((email) => ({
              value: email.email,
              label: `${e.firstName} ${e.lastName}`,
            }));
          }

          return [
            {
              label: getName(e),
              value: '',
            },
          ];
        })
        .flat();

      callback(options);
    } catch (error) {
      callback([]);
    }
  };

  const Option = useCallback((rest: OptionProps<SelectOption>) => {
    const fullLabel =
      rest?.data?.label &&
      rest?.data?.value &&
      `${rest.data.label} - ${rest.data.value}`;
    const emailOnly =
      !rest?.data?.label && rest?.data?.value && `${rest.data.value}`;
    const noEmail = rest?.data?.label && !rest?.data?.value && (
      <Text>
        {rest.data.label} -
        <Text as='span' color='gray.500' ml={1}>
          [No email for this contact]
        </Text>
      </Text>
    );

    return (
      <chakraComponents.Option {...rest}>
        {fullLabel || emailOnly || noEmail}
        {rest?.isFocused && (
          <IconButton
            aria-label='Copy'
            size='xs'
            p={0}
            height={5}
            variant='ghost'
            icon={<Copy01 boxSize={3} color='gray.500' />}
            onClick={(e) => {
              e.stopPropagation();
              copyToClipboard(rest.data.value, 'Email copied!');
            }}
          />
        )}
      </chakraComponents.Option>
    );
  }, []);

  const components = useMemo(
    () => ({
      MultiValueRemove: () => null,
      MultiValue: (multiValueProps: MultiValueProps<SelectOption>) => (
        <MultiValueWithActionMenu
          {...multiValueProps}
          name={name}
          formId={formId}
          existingContacts={existingContacts}
        />
      ),
    }),
    [existingContacts, name, formId],
  );

  return (
    <MultiCreatableSelect
      ref={ref}
      id={id}
      formId={formId}
      name={name}
      value={value}
      onBlur={(e) => handleBlur(e.target.value)}
      onChange={onChange}
      Option={Option}
      customStyles={multiCreatableSelectStyles}
      components={components}
      loadOptions={(inputValue: string, callback) => {
        getFilteredSuggestions(inputValue, callback);
      }}
      {...rest}
    />
  );
});
