import { components } from 'react-select';
import { MultiValueProps } from 'react-select';
import { useParams, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn.ts';
import { DataSource } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@ui/utils/types.ts';
import { Copy01 } from '@ui/media/icons/Copy01.tsx';
import { validateEmail } from '@shared/util/emailValidation.ts';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { useContactCardMeta } from '@organization/state/ContactCardMeta.atom.ts';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

interface MultiValueWithActionMenuProps extends MultiValueProps<SelectOption> {
  value: Array<SelectOption>;
  navigateAfterAddingToPeople: boolean;
  onChange: (newValue: Array<SelectOption<string>>) => void;
  existingContacts: Array<{ id: string; label: string; value?: string | null }>;
}

export const MultiValueWithActionMenu = observer(
  ({
    existingContacts,
    navigateAfterAddingToPeople,
    onChange,
    value,
    ...rest
  }: MultiValueWithActionMenuProps) => {
    const store = useStore();
    const [searchParams, setSearchParams] = useSearchParams();
    const organizationId = useParams()?.id as string;
    const [_d, setExpandedCardId] = useContactCardMeta();

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
              .map(
                (word: string) => word.charAt(0).toUpperCase() + word.slice(1),
              )
              .join(' ');

      store.contacts.create(organizationId, {
        onSuccess: (newContactId) => {
          const contact = store.contacts.value.get(newContactId);

          contact?.update((d) => {
            d.name = name;
            d.emails = [
              {
                email: rest?.data?.value,
                appSource: '',
                contacts: [],
                createdAt: undefined,
                emailValidationDetails: {
                  __typename: undefined,
                  verified: false,
                  verifyingCheckAll: false,
                  isValidSyntax: undefined,
                  isRisky: undefined,
                  isFirewalled: undefined,
                  provider: undefined,
                  firewall: undefined,
                  isCatchAll: undefined,
                  canConnectSmtp: undefined,
                  isDeliverable: undefined,
                  isMailboxFull: undefined,
                  isRoleAccount: undefined,
                  isFreeAccount: undefined,
                  smtpSuccess: undefined,
                },
                id: '',
                organizations: [],
                primary: false,
                source: DataSource.Openline,
                sourceOfTruth: DataSource.Openline,
                updatedAt: undefined,
                users: [],
              },
            ];

            return d;
          });
        },
      });
    };

    return (
      <Menu>
        <MenuButton
          className={cn(
            isContactWithoutEmail
              ? 'text-base [&_.multiValueClass]:data-[state=closed]:bg-warning-50 [&_.multiValueClass]:data-[state=closed]:text-warning-700 [&_.multiValueClass]:data-[state=closed]:border-warning-200 [&_.multiValueClass]:data-[state=open]:bg-warning-50 [&_.multiValueClass]:data-[state=open]:text-warning-700 [&_.multiValueClass]:data-[state=open]:border-warning-200'
              : 'text-base [&_.multiValueClass]:data-[state=closed]:bg-gray-50 [&_.multiValueClass]:data-[state=closed]:text-gray-700 [&_.multiValueClass]:data-[state=closed]:border-gray-200 [&_.multiValueClass]:data-[state=open]:bg-primary-50 [&_.multiValueClass]:data-[state=open]:text-primary-700 [&_.multiValueClass]:data-[state=open]:last:border-primary-200',
          )}
        >
          <components.MultiValue {...rest} className='text-base'>
            {rest.children}
          </components.MultiValue>
        </MenuButton>
        <div onPointerDown={(e) => e.stopPropagation()}>
          <MenuList side='bottom' align='start' className='max-w-[300px] p-2'>
            {rest?.data?.value ? (
              <MenuItem
                onPointerDown={() => {
                  copyToClipboard(rest?.data?.value, 'Email copied');
                }}
                className='flex justify-between items-center rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
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
                onPointerDown={() => {
                  handleAddContact();
                }}
                className='rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
              >
                Add to people
              </MenuItem>
            )}
          </MenuList>
        </div>
      </Menu>
    );
  },
);
