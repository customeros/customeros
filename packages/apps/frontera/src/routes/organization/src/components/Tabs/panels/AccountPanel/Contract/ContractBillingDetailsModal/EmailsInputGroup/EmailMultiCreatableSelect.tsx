import { useParams } from 'react-router-dom';
import { useMemo, useState, useEffect, forwardRef, useCallback } from 'react';
import {
  GroupBase,
  OptionProps,
  SelectInstance,
  OptionsOrGroups,
  MultiValueProps,
  components as reactSelectComponents,
} from 'react-select';

import merge from 'lodash/merge';
import AsyncCreatableSelect from 'react-select/async-creatable';

import { SelectOption } from '@ui/utils/types';
import { Copy01 } from '@ui/media/icons/Copy01';
import { getName } from '@utils/getParticipantsName';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Contact, ComparisonOperator } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import {
  getDefaultClassNames,
  getMultiValueLabelClassNames,
} from '@ui/form/CreatableSelect';
import {
  GetContactsEmailListDocument,
  useGetContactsEmailListQuery,
} from '@organization/graphql/getContactsEmailList.generated';

import { MultiValueWithActionMenu } from './MultiValueWithActionMenu.tsx';

type ExistingContact = { id: string; label: string; value?: string | null };
export const EmailMultiCreatableSelect = forwardRef<
  SelectInstance,
  {
    isMulti: boolean;
    placeholder?: string;
    noOptionsMessage: () => null;
    value: SelectOption<string>[];
    navigateAfterAddingToPeople: boolean;
    onChange: (value: SelectOption<string>[]) => void;
  }
>(({ value, onChange, navigateAfterAddingToPeople, isMulti, ...rest }) => {
  const client = getGraphQLClient();
  const organizationId = useParams()?.id as string;
  const [existingContacts, setExistingContacts] = useState<
    Array<ExistingContact>
  >([]);

  const [isFocused, setIsFocused] = useState<boolean>(false);

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

  const [_, copyToClipboard] = useCopyToClipboard();

  const getFilteredSuggestions = async (
    filterString: string,
    callback: (
      options: OptionsOrGroups<SelectOption, GroupBase<SelectOption>>,
    ) => void,
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
      const options = (results?.organization?.contacts?.content || [])
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
        .flat() as OptionsOrGroups<SelectOption, GroupBase<SelectOption>>;

      callback(options);
    } catch (error) {
      callback([]);
    }
  };

  const Option = useCallback((rest: OptionProps<SelectOption>) => {
    const fullLabel =
      rest?.data?.label.length > 1 &&
      rest?.data?.value.length > 1 &&
      `${rest.data.label}  ${rest.data.value}`;

    const emailOnly = rest?.data?.value.length > 1 && `${rest.data.value}`;

    const noEmail = rest?.data?.label && !rest?.data?.value && (
      <p>
        {rest.data.label} -
        <span className='text-gray-500 ml-1'>[No email for this contact]</span>
      </p>
    );

    return (
      <reactSelectComponents.Option {...rest}>
        {fullLabel || emailOnly || noEmail}
        {rest?.isFocused && (
          <IconButton
            className='h-5 p-0 self-end float-end'
            aria-label='Copy'
            size='xs'
            variant='ghost'
            icon={<Copy01 className='size-3 text-gray-500' />}
            onClick={(e) => {
              e.stopPropagation();
              copyToClipboard(rest.data.value, 'Email copied');
            }}
          />
        )}
      </reactSelectComponents.Option>
    );
  }, []);

  const MultiValue = useCallback(
    (multiValueProps: MultiValueProps<SelectOption>) => {
      return (
        <MultiValueWithActionMenu
          {...multiValueProps}
          navigateAfterAddingToPeople={navigateAfterAddingToPeople}
          existingContacts={existingContacts}
          value={value}
          onChange={onChange}
        />
      );
    },
    [navigateAfterAddingToPeople],
  );

  const components = useMemo(
    () => ({
      MultiValueRemove: () => null,
      LoadingIndicator: () => null,
      DropdownIndicator: () => null,
      Option,
      MultiValue,
    }),
    [MultiValue, Option],
  );
  const defaultClassNames = useMemo(
    () => merge(getDefaultClassNames({ size: 'md' })),
    [],
  );

  return (
    <AsyncCreatableSelect
      cacheOptions
      closeMenuOnSelect={false}
      isMulti={isMulti}
      unstyled
      isClearable={false}
      menuIsOpen
      onFocus={() => setIsFocused(true)}
      tabSelectsValue={true}
      id={'email-multi-creatable-select'}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      value={value?.map((e: { value: any; label: string | any[] }) => ({
        label: e.label.length > 1 ? e.label : e.value,
        value: e.value,
      }))}
      classNames={{
        ...defaultClassNames,
        singleValue: () =>
          isFocused ? getMultiValueLabelClassNames('', 'md') : 'text-gray-500',
      }}
      onBlur={() => setIsFocused(false)}
      // @ts-expect-error fix me later
      onChange={onChange}
      defaultMenuIsOpen
      components={components}
      loadOptions={(inputValue: string, callback) => {
        getFilteredSuggestions(inputValue, callback);
      }}
      formatCreateLabel={(input: string) => {
        return input;
      }}
      {...rest}
    />
  );
});
