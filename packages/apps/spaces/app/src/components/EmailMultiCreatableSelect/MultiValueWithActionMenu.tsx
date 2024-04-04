import React, { FC } from 'react';
import { useField } from 'react-inverted-form';
import { useParams, useRouter, useSearchParams } from 'next/navigation';

import { useLocalStorage } from 'usehooks-ts';
import { MultiValueProps } from 'chakra-react-select';
import { useQueryClient } from '@tanstack/react-query';

import { SelectOption } from '@ui/utils';
import { Copy01 } from '@ui/media/icons/Copy01';
import { toastSuccess } from '@ui/presentation/Toast';
import { chakraComponents } from '@ui/form/SyncSelect';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { useContactCardMeta } from '@organization/src/state/ContactCardMeta.atom';
import { invalidateQuery } from '@organization/src/components/Tabs/panels/PeoplePanel/util';
import { useCreateContactMutation } from '@organization/src/graphql/createContact.generated';
import {
  Menu,
  MenuItem,
  MenuButton,
  MenuList as ChakraMenuList,
} from '@ui/overlay/Menu';
import { useAddOrganizationToContactMutation } from '@organization/src/graphql/addContactToOrganization.generated';

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
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const searchParams = useSearchParams();
  const router = useRouter();
  const organizationId = useParams()?.id as string;
  const [_d, setExpandedCardId] = useContactCardMeta();
  const { getInputProps } = useField(name, formId);
  const { onChange } = getInputProps();
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

  const isContactWithoutEmail = isContactInOrg && !rest?.data?.value;

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

    router.push(`?${urlSearchParams}`);
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
          const contactId = data.contact_Create.id;
          addContactToOrganization.mutate({
            input: { contactId, organizationId },
          });
        },
        onSettled: (data) => {
          if (navigateAfterAddingToPeople) {
            handleNavigateToContact(data?.contact_Create?.id as string, 'name');
          }
          toastSuccess(
            'Contact added to people list',
            data?.contact_Create?.id as string,
          );
        },
      },
    );
  };

  return (
    <Menu>
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
      <ChakraMenuList maxW={300} p={2}>
        {rest?.data?.value ? (
          <MenuItem
            display='flex'
            borderRadius='md'
            border='1px solid'
            borderColor='transparent'
            _hover={{
              bg: 'gray.50',
              borderColor: 'gray.100',
            }}
            _focus={{
              borderColor: 'gray.200',
            }}
            justifyContent='space-between'
            onClick={(e) => {
              e.stopPropagation();
              copyToClipboard(rest?.data?.value, 'Email copied');
            }}
          >
            {rest?.data?.value}
            <Copy01 boxSize={3} color='gray.500' ml={2} />
          </MenuItem>
        ) : (
          <MenuItem
            borderRadius='md'
            border='1px solid'
            borderColor='transparent'
            _hover={{
              bg: 'gray.50',
              borderColor: 'gray.100',
            }}
            _focus={{
              borderColor: 'gray.200',
            }}
            onClick={() =>
              isContactInOrg &&
              handleNavigateToContact(isContactInOrg.id, 'email')
            }
          >
            Add email in People list
          </MenuItem>
        )}

        <MenuItem
          borderRadius='md'
          border='1px solid'
          borderColor='transparent'
          _hover={{
            bg: 'gray.50',
            borderColor: 'gray.100',
          }}
          _focus={{
            borderColor: 'gray.200',
          }}
          onClick={() => {
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
            onClick={handleAddContact}
            borderRadius='md'
            border='1px solid'
            borderColor='transparent'
            _hover={{
              bg: 'gray.50',
              borderColor: 'gray.100',
            }}
            _focus={{
              borderColor: 'gray.200',
            }}
          >
            Add to people
          </MenuItem>
        )}
      </ChakraMenuList>
    </Menu>
  );
};
