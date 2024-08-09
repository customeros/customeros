import React, { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';

import Fuse from 'fuse.js';
import { useCommandState } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { Avatar } from '@ui/media/Avatar';
import { Organization } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Command } from '@ui/overlay/CommandMenu';
import { User03 } from '@ui/media/icons/User03.tsx';

const isContact = (
  item:
    | { type: string; item: ContactStore }
    | { type: string; item: OrganizationStore },
): item is { type: 'contact'; item: ContactStore } => {
  return item.type === 'contact';
};

const isOrganization = (
  item:
    | { type: string; item: ContactStore }
    | { type: string; item: OrganizationStore },
): item is { type: 'organization'; item: OrganizationStore } => {
  return item.type === 'organization';
};
export const GlobalSearchResultNavigationCommands = observer(() => {
  const search = useCommandState((state) => state.search);
  const navigate = useNavigate();

  const { contacts, organizations, ui } = useStore();
  const combinedList = useMemo(() => {
    return [
      ...contacts
        .toArray()
        .map((contact) => ({ type: 'contact', item: contact })),
      ...organizations
        .toArray()
        .map((org) => ({ type: 'organization', item: org })),
    ];
  }, [contacts.totalElements, organizations.totalElements]);

  const fuseCombined = useMemo(
    () =>
      new Fuse(combinedList, {
        keys: ['item.name', 'item.value.name'],
        threshold: 0.3,
      }),
    [combinedList],
  );

  const { filteredContacts, filteredOrgs } = useMemo(() => {
    if (!search) return { filteredContacts: [], filteredOrgs: [] };

    const results = fuseCombined.search(search, { limit: 10 });
    const { filteredContacts, filteredOrgs } = results.reduce<{
      filteredContacts: ContactStore[];
      filteredOrgs: OrganizationStore['value'][];
    }>(
      (acc, result) => {
        if (isContact(result.item)) {
          acc.filteredContacts.push(result.item.item);
        }

        if (isOrganization(result.item)) {
          acc.filteredOrgs.push(result.item.item.value as Organization);
        }

        return acc;
      },
      { filteredContacts: [], filteredOrgs: [] },
    );

    return {
      filteredContacts,
      filteredOrgs,
    };
  }, [search, fuseCombined]);

  const handleGoTo = (id: string, tab: string) => {
    navigate(`/organization/${id}?tab=${tab}`);
    ui.commandMenu.setOpen(false);
  };

  return (
    <>
      <Command.Group>
        {filteredContacts?.map((contactStore) => (
          <Command.Item
            key={contactStore.value.metadata.id}
            onSelect={() => handleGoTo(contactStore.organizationId, 'people')}
            value={
              `${contactStore.name} ${contactStore.value.metadata.id}` ?? ''
            }
          >
            <div className='flex items-center'>
              <Avatar
                size='xs'
                textSize='xs'
                name={contactStore.name}
                icon={<User03 className='text-primary-700  ' />}
                src={
                  contactStore?.value?.profilePhotoUrl
                    ? contactStore.value.profilePhotoUrl
                    : undefined
                }
              />
              <span className='ml-2 capitalize'>{contactStore.name}</span>

              <span className='ml-1.5 text-gray-500'>
                ·{' '}
                {contactStore.organizationId
                  ? organizations.value.get(contactStore.organizationId)?.value
                      .name
                  : ''}
              </span>
            </div>
          </Command.Item>
        ))}
      </Command.Group>
      <Command.Group>
        {filteredOrgs?.map((org) => (
          <Command.Item
            value={org.name}
            key={org.metadata.id}
            onSelect={() => handleGoTo(org.metadata.id, 'about')}
          >
            <div className='flex items-center'>
              <Avatar
                size='xxs'
                textSize='xs'
                name={org.name}
                className='mr-2'
                variant='roundedSquare'
                src={org.icon || org.logo || undefined}
              />

              {org.name}
            </div>
          </Command.Item>
        ))}
      </Command.Group>
    </>
  );
});
