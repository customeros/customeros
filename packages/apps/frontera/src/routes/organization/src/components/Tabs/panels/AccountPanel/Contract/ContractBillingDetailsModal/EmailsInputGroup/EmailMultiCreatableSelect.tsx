import { useParams } from 'react-router-dom';
import { SelectComponentsConfig } from 'react-select';
import { useMemo, useState, useEffect, forwardRef, useCallback } from 'react';
import {
  GroupBase,
  OptionProps,
  SelectInstance,
  OptionsOrGroups,
  MultiValueProps,
  components as reactSelectComponents,
} from 'react-select';

import { twMerge } from 'tailwind-merge';
import AsyncCreatableSelect from 'react-select/async-creatable';

import { cn } from '@ui/utils/cn.ts';
import { SelectOption } from '@ui/utils/types';
import { Copy01 } from '@ui/media/icons/Copy01';
import { getName } from '@utils/getParticipantsName';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Contact, ComparisonOperator } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { emailRegex } from '@organization/components/Timeline/PastZone/events/email/utils';
import {
  GetContactsEmailListDocument,
  useGetContactsEmailListQuery,
} from '@organization/graphql/getContactsEmailList.generated';

import { MultiValueWithActionMenu } from './MultiValueWithActionMenu.tsx';

type ExistingContact = { id: string; label: string; value?: string | null };
export const EmailMultiCreatableSelect = forwardRef<
  SelectInstance,
  {
    placeholder?: string;
    noOptionsMessage: () => null;
    value: SelectOption<string>[];
    navigateAfterAddingToPeople: boolean;
    onChange: (value: SelectOption<string>[]) => void;
  }
>(({ value, onChange, navigateAfterAddingToPeople, ...rest }) => {
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

  const [_, copyToClipboard] = useCopyToClipboard();

  const handleBlur = (stringVal: string) => {
    if (stringVal && emailRegex.test(stringVal)) {
      // onBlur([...value, { label: stringVal, value: stringVal }]);

      return;
    }
    // onBlur(value);
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

  return (
    <AsyncCreatableSelect
      cacheOptions
      closeMenuOnSelect={false}
      isMulti
      unstyled
      isClearable={false}
      tabSelectsValue={true}
      id={'test'}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      value={value?.map((e: { value: any; label: string | any[] }) => ({
        label: e.label.length > 1 ? e.label : e.value,
        value: e.value,
      }))}
      classNames={{
        container: ({ isFocused }) =>
          twMerge(
            'flex mt-1 cursor-pointer overflow-visible min-w-[300px] w-full focus-visible:border-0 focus:border-0',

            isFocused && 'border-primary-500',
          ),
        menu: ({ menuPlacement }) =>
          cn(
            menuPlacement === 'top'
              ? 'mb-2 animate-slideUpAndFade'
              : 'mt-2 animate-slideDownAndFade',
            'z-50',
          ),
        menuList: () =>
          'p-2 z-50  max-h-[12rem] border border-gray-200 bg-white rounded-lg shadow-lg overflow-y-auto overscroll-auto',
        option: ({ isFocused, isSelected }) =>
          cn(
            'my-[2px] px-3 py-1.5 rounded-md text-gray-700 line-clamp-1 text-sm transition ease-in-out delay-50 hover:bg-gray-50',
            isSelected && 'bg-gray-50 font-medium leading-normal',
            isFocused && 'ring-2 ring-gray-100',
          ),
        placeholder: () => 'text-gray-400 text-inherit',
        multiValue: () =>
          'flex p-0 gap-0 text-gray-700 text-inherit mr-1 cursor-default h-[auto]',

        multiValueRemove: () => 'hidden',
        groupHeading: () =>
          'text-gray-400 text-sm px-3 py-1.5 font-normal uppercase',
        control: () => 'overflow-visible',
        input: () => 'overflow-visible text-gray-500 leading-4',
        multiValueLabel: () =>
          'multiValueClass px-2 bg-transparent text-inherit shadow-md border font-semibold rounded-lg border-gray-200 max-h-[12rem] cursor-pointer z-50',
        valueContainer: () => 'w-full',
      }}
      onBlur={(e) => handleBlur(e.target.value)}
      // @ts-expect-error fix me later
      onChange={onChange}
      defaultMenuIsOpen
      components={
        components as SelectComponentsConfig<unknown, true, GroupBase<unknown>>
      }
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
