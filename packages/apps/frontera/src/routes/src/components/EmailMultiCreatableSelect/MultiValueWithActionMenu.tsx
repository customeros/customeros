import { FC, useState } from 'react';
import { components } from 'react-select';
import { MultiValueProps } from 'react-select';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';
import { EditEmailCase } from '@domain/usecases/command-menu/edit-email.usecase';

import { cn } from '@ui/utils/cn';
import { validateEmail } from '@utils/email';
import { SelectOption } from '@ui/utils/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { Copy01 } from '@ui/media/icons/Copy01';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { useContactCardMeta } from '@organization/state/ContactCardMeta.atom';

interface MultiValueWithActionMenuProps extends MultiValueProps<SelectOption> {
  name: string;
  navigateAfterAddingToPeople: boolean;
  removeOption: (value: string) => void;
}

const editEmailUseCase = EditEmailCase.getInstance();

export const MultiValueWithActionMenu: FC<MultiValueWithActionMenuProps> =
  observer(({ name, navigateAfterAddingToPeople, removeOption, ...rest }) => {
    const [isOpen, setIsOpen] = useState(false);
    const store = useStore();
    const navigate = useNavigate();

    const [searchParams, setSearchParams] = useSearchParams();
    const organizationId = useParams()?.id as string;
    const existingContacts =
      store.organizations.getById(organizationId)?.contacts;

    const [_d, setExpandedCardId] = useContactCardMeta();

    const [_, copyToClipboard] = useCopyToClipboard();
    const [lastActivePosition, setLastActivePosition] = useLocalStorage(
      `customeros-player-last-position`,
      { [organizationId as string]: 'tab=about' },
    );
    const isContactInOrg = existingContacts?.find((data) => {
      return rest?.data?.value
        ? data.emails.some((e) => e.email === rest.data.value)
        : false;
    });

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

    return (
      <Menu
        onOpenChange={(newState) => {
          setIsOpen(newState);
        }}
      >
        <MenuButton
          className={cn(
            isContactWithoutEmail
              ? '[&_.multiValueClass]:data-[state=closed]:bg-warning-50 [&_.multiValueClass]:data-[state=closed]:text-warning-700 [&_.multiValueClass]:data-[state=closed]:border-warning-200 [&_.multiValueClass]:data-[state=open]:bg-warning-50 [&_.multiValueClass]:data-[state=open]:text-warning-700 [&_.multiValueClass]:data-[state=open]:border-warning-200'
              : 'hover:bg-grayModern-100  rounded-sm px-[2px] [&_.multiValueClass]:data-[state=closed]:bg-gray-50 [&_.multiValueClass]:data-[state=closed]:text-gray-700 [&_.multiValueClass]:data-[state=closed]:border-gray-200 [&_.multiValueClass]:data-[state=open]:bg-primary-50 [&_.multiValueClass]:data-[state=open]:text-primary-700 [&_.multiValueClass]:data-[state=open]:last:border-primary-200',
            {
              'bg-grayModern-100 ': isOpen,
              'bg-transparent hover:bg-grayModern-100 focus:bg-grayModern-100':
                !isOpen,
            },
          )}
        >
          <components.MultiValue {...rest} className={'rounded-md'}>
            {rest.data.label ? (
              <Tooltip label={rest.data.value}>
                <div>{rest.data.label}</div>
              </Tooltip>
            ) : (
              rest.data.value
            )}
          </components.MultiValue>
        </MenuButton>
        <div onPointerDown={(e) => e.stopPropagation()}>
          <MenuList side='bottom' align='start' className='max-w-[300px] p-2'>
            {isContactInOrg?.id && (
              <MenuItem
                className='flex justify-between items-center rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
                onPointerDown={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  editEmailUseCase.setEmail(rest.data.value);
                  store.ui.commandMenu.setType('EditEmail');

                  store.ui.commandMenu.setContext({
                    ids: [isContactInOrg.id],
                    entity: 'Contact',
                    property: 'email',
                  });
                  store.ui.commandMenu.setOpen(true);
                }}
              >
                Edit address
                <Edit03 className='size-3 text-gray-500 ml-2' />
              </MenuItem>
            )}

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
              onPointerDown={() => {
                removeOption(rest.data.value);
              }}
              className='rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
            >
              Remove address
            </MenuItem>
            {!isContactInOrg && (
              <MenuItem
                className='rounded-md border border-transparent hover:bg-gray-50 hover:border-gray-100 focus:border-gray-200'
                onPointerDown={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  store.ui.commandMenu.setType('AddSingleContact');
                  store.ui.commandMenu.setContext({
                    ids: [organizationId],
                    entity: 'Organization',
                    property: 'contacts',
                    meta: {
                      email: rest.data.value,
                      callback: () => {
                        if (navigateAfterAddingToPeople) {
                          navigate(`?tab=people`);
                        }
                      },
                    },
                  });
                  store.ui.commandMenu.setOpen(true);
                }}
              >
                Add to people
              </MenuItem>
            )}
          </MenuList>
        </div>
      </Menu>
    );
  });
