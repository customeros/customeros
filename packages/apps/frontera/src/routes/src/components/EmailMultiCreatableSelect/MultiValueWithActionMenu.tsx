import { components } from 'react-select';
import { useField } from 'react-inverted-form';
import { MultiValueProps } from 'react-select';
import { FC, useState, ReactEventHandler } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';
import { useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input/Input';
import { SelectOption } from '@ui/utils/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { Copy01 } from '@ui/media/icons/Copy01';
import { toastSuccess } from '@ui/presentation/Toast';
import { validateEmail } from '@shared/util/emailValidation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { useContactCardMeta } from '@organization/state/ContactCardMeta.atom';
import { useCreateContactMutation } from '@organization/graphql/createContact.generated';
import { useAddOrganizationToContactMutation } from '@organization/graphql/addContactToOrganization.generated';

import { invalidateQuery } from './util';

interface MultiValueWithActionMenuProps extends MultiValueProps<SelectOption> {
  name: string;
  formId: string;
  navigateAfterAddingToPeople: boolean;
  existingContacts: Array<{ id: string; label: string; value?: string | null }>;
}

export const MultiValueWithActionMenu: FC<MultiValueWithActionMenuProps> = ({
  existingContacts,
  name,
  formId,
  navigateAfterAddingToPeople,
  ...rest
}) => {
  const [editInput, setEditInput] = useState(false);

  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [searchParams, setSearchParams] = useSearchParams();
  const organizationId = useParams()?.id as string;
  const [_d, setExpandedCardId] = useContactCardMeta();
  const { getInputProps } = useField(name, formId);
  const { onChange, value } = getInputProps();
  const createContact = useCreateContactMutation(client);
  const addContactToOrganization = useAddOrganizationToContactMutation(client, {
    onSuccess: () => invalidateQuery(queryClient, organizationId),
  });
  const [_, copyToClipboard] = useCopyToClipboard();
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [organizationId as string]: 'tab=about' },
  );
  const isContactInOrg = existingContacts.find(
    (data: SelectOption | unknown) => {
      return rest?.data?.value
        ? (data as SelectOption)?.value === rest.data.value
        : rest.data.label?.trim() === (data as SelectOption)?.label?.trim();
    },
  );

  const validationMessage = validateEmail(rest?.data?.value);

  const isContactWithoutEmail =
    (isContactInOrg && !rest?.data?.value) || validationMessage;

  const handleNavigateToContact = (
    contactId: string,
    initialFocusedField: 'name' | 'email',
  ) => {
    const urlSearchParams = new URLSearchParams(searchParams?.toString());
    urlSearchParams.set('tab', 'people');
    setLastActivePosition({
      ...lastActivePosition,
      [organizationId as string]: urlSearchParams.toString(),
    });

    setSearchParams(urlSearchParams);
    setExpandedCardId({
      expandedId: contactId,
      initialFocusedField,
    });
  };
  const handleAddContact = () => {
    const name =
      rest?.data?.label !== rest?.data?.value
        ? rest?.data?.label
        : rest?.data?.label
            ?.split('@')?.[0]
            ?.split('.')
            .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ');
    createContact.mutate(
      {
        input: {
          name,
          email: { email: rest?.data?.value },
        },
      },
      {
        onSuccess: (data) => {
          const contactId = data.contact_Create;
          addContactToOrganization.mutate({
            input: { contactId, organizationId },
          });
        },
        onSettled: (data) => {
          if (navigateAfterAddingToPeople) {
            handleNavigateToContact(data?.contact_Create as string, 'name');
          }
          toastSuccess(
            'Contact added to people list',
            data?.contact_Create as string,
          );
        },
      },
    );
  };

  if (editInput) {
    const handleChangeValue: ReactEventHandler<HTMLElement> = (event) => {
      const newValue = value.map((e: SelectOption<string>) =>
        e.value === rest?.data?.value
          ? {
              label: (event?.target as HTMLInputElement)?.value,
              value: (event?.target as HTMLInputElement)?.value,
            }
          : e,
      );
      onChange(newValue);
      setEditInput(false);
    };

    return (
      <Input
        className='w-auto inline text-warning-700 h-8'
        variant='unstyled'
        onBlur={(e) => {
          handleChangeValue(e);
        }}
        onKeyDown={(e) => {
          e.stopPropagation();
          if (e.key === 'Enter') {
            handleChangeValue(e);
          }
        }}
        defaultValue={rest?.data?.value}
      />
    );
  }
  const handleEditInput = () => {
    setEditInput(true);
  };

  return (
    <Menu>
      <MenuButton
        className={cn(
          isContactWithoutEmail
            ? '[&_.multiValueClass]:data-[state=closed]:bg-warning-50 [&_.multiValueClass]:data-[state=closed]:text-warning-700 [&_.multiValueClass]:data-[state=closed]:border-warning-200 [&_.multiValueClass]:data-[state=open]:bg-warning-50 [&_.multiValueClass]:data-[state=open]:text-warning-700 [&_.multiValueClass]:data-[state=open]:border-warning-200'
            : '[&_.multiValueClass]:data-[state=closed]:bg-gray-50 [&_.multiValueClass]:data-[state=closed]:text-gray-700 [&_.multiValueClass]:data-[state=closed]:border-gray-200 [&_.multiValueClass]:data-[state=open]:bg-primary-50 [&_.multiValueClass]:data-[state=open]:text-primary-700 [&_.multiValueClass]:data-[state=open]:last:border-primary-200',
        )}
      >
        <components.MultiValue {...rest}>{rest.children}</components.MultiValue>
      </MenuButton>
      <div onPointerDown={(e) => e.stopPropagation()}>
        <MenuList className='max-w-[300px] p-2' side='bottom' align='start'>
          <MenuItem
            className='flex justify-between items-center rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
            onPointerDown={() => {
              handleEditInput();
            }}
          >
            Edit address
            <Edit03 className='size-3 text-gray-500 ml-2' />
          </MenuItem>
          {rest?.data?.value ? (
            <MenuItem
              className='flex justify-between items-center rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
              onPointerDown={() => {
                copyToClipboard(rest?.data?.value, 'Email copied');
              }}
            >
              {rest?.data?.value}
              <Copy01 className='size-3 text-gray-500 ml-2' />
            </MenuItem>
          ) : (
            <MenuItem
              className='rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
              onPointerDown={() => {
                isContactInOrg &&
                  handleNavigateToContact(isContactInOrg.id, 'email');
              }}
            >
              Add email in People list
            </MenuItem>
          )}

          <MenuItem
            className='rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
            onPointerDown={() => {
              const newValue = (
                (rest?.selectProps?.value as Array<SelectOption>) ?? []
              )?.filter((e: SelectOption) => e.value !== rest?.data?.value);
              onChange(newValue);
            }}
          >
            Remove address
          </MenuItem>
          {!isContactInOrg && (
            <MenuItem
              className='rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
              onPointerDown={() => {
                handleAddContact();
              }}
            >
              Add to people
            </MenuItem>
          )}
        </MenuList>
      </div>
    </Menu>
  );
};
