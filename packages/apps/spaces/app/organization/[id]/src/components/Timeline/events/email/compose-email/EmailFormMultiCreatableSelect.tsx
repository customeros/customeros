import { OptionsOrGroups } from 'react-select';
import { useField } from 'react-inverted-form';
import React, { useMemo, forwardRef, useCallback } from 'react';
import { useParams, useRouter, useSearchParams } from 'next/navigation';

import { useLocalStorage } from 'usehooks-ts';
import { useQueryClient } from '@tanstack/react-query';
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
import { invalidateQuery } from '@organization/src/components/Tabs/panels/PeoplePanel/util';
import { useCreateContactMutation } from '@organization/src/graphql/createContact.generated';
import {
  Menu,
  MenuItem,
  MenuButton,
  MenuList as ChakraMenuList,
} from '@ui/overlay/Menu';
import { useAddOrganizationToContactMutation } from '@organization/src/graphql/addContactToOrganization.generated';
import {
  GetContactsEmailListDocument,
  useGetContactsEmailListQuery,
} from '@organization/src/graphql/getContactsEmailList.generated';

export const EmailFormMultiCreatableSelect = forwardRef<
  SelectInstance,
  FormSelectProps
>(({ name, formId, ...rest }, ref) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const organizationId = useParams()?.id as string;
  const searchParams = useSearchParams();
  const router = useRouter();
  const createContact = useCreateContactMutation(client);
  const addContactToOrganization = useAddOrganizationToContactMutation(client, {
    onSuccess: () => invalidateQuery(queryClient, organizationId),
  });
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [organizationId as string]: 'tab=about' },
  );

  const { data } = useGetContactsEmailListQuery(client, {
    id: organizationId,
    pagination: {
      page: 1,
      limit: 100,
    },
  });

  const handleNavigateToContact = () => {
    const urlSearchParams = new URLSearchParams(searchParams?.toString());
    urlSearchParams.set('tab', 'people');
    setLastActivePosition({
      ...lastActivePosition,
      [organizationId as string]: urlSearchParams.toString(),
    });

    router.push(`?${urlSearchParams}`);
  };
  const handleAddContact = ({
    name,
    email,
  }: {
    name: string;
    email: string;
  }) => {
    createContact.mutate(
      {
        input: {
          name,
          email: { email },
        },
      },
      {
        onSuccess: (data) => {
          const contactId = data.contact_Create.id;
          addContactToOrganization.mutate({
            input: { contactId, organizationId },
          });
        },
        onSettled: handleNavigateToContact,
      },
    );
  };

  const organizationContacts: OptionsOrGroups<unknown, GroupBase<unknown>> = (
    (data?.organization?.contacts?.content || []) as Array<Contact>
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

  const handleAddToPeople = () => {};
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
  const MultiValue = useCallback(
    (rest: MultiValueProps<SelectOption>) => {
      const isContactInOrg = organizationContacts.findIndex(
        (data: SelectOption | unknown) => {
          return rest?.data?.value
            ? (data as SelectOption)?.value === rest?.data?.value
            : rest?.data?.label === (data as SelectOption)?.label;
        },
      );
      const isContactWithoutEmail = isContactInOrg && !rest?.data?.value;
      const name =
        rest?.data?.label !== rest?.data?.value
          ? rest?.data?.label
          : rest?.data?.label
              ?.split('@')?.[0]
              ?.split('.')
              .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
              .join(' ');

      return (
        <Menu isLazy>
          <MenuButton
            sx={{
              '&[aria-expanded="true"] > span > span': {
                bg: isContactWithoutEmail
                  ? 'warning.50 !important'
                  : 'primary.50 !important',
                color: isContactWithoutEmail
                  ? 'warning.700 !important'
                  : 'primary.700 !important',
                borderColor: isContactWithoutEmail
                  ? 'warning.200 !important'
                  : 'primary.200 !important',
              },
            }}
          >
            <chakraComponents.MultiValue {...rest}>
              {rest.children}
            </chakraComponents.MultiValue>
          </MenuButton>
          <ChakraMenuList maxW={300}>
            {rest?.data?.value ? (
              <MenuItem
                display='flex'
                justifyContent='space-between'
                onClick={(e) => {
                  e.stopPropagation();
                  copyToClipboard(rest?.data?.value, 'Email copied!');
                }}
              >
                {rest?.data?.value}
                <Copy01 boxSize={3} color='gray.500' ml={2} />
              </MenuItem>
            ) : (
              <MenuItem onClick={handleNavigateToContact}>
                Add email in People list
              </MenuItem>
            )}

            <MenuItem
              onClick={() => {
                const newValue = (
                  (rest?.selectProps?.value as Array<SelectOption>) ?? []
                )?.filter((e: SelectOption) => e.value !== rest?.data?.value);
                onChange(newValue);
              }}
            >
              Remove address
            </MenuItem>
            {isContactInOrg < 0 && (
              <MenuItem
                onClick={() => {
                  handleAddContact({
                    name,
                    email: rest?.data?.value,
                  });
                }}
              >
                Add to people
              </MenuItem>
            )}
          </ChakraMenuList>
        </Menu>
      );
    },
    [organizationContacts, searchParams, handleAddToPeople],
  );

  const components = useMemo(
    () => ({
      MultiValueRemove: () => null,
    }),
    [],
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
      options={organizationContacts}
      MultiValue={MultiValue}
      customStyles={multiCreatableSelectStyles}
      components={components}
      loadOptions={(inputValue: string, callback) => {
        getFilteredSuggestions(inputValue, callback);
      }}
      {...rest}
    />
  );
});
