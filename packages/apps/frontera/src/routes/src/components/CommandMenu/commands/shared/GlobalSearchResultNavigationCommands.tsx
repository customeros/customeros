import React, { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';

import Fuse from 'fuse.js';
import { useCommandState } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { Contact } from '@store/Contacts/Contact.dto';
import { Organization } from '@store/Organizations/Organization.dto';

import { Avatar } from '@ui/media/Avatar';
import { useStore } from '@shared/hooks/useStore';
import { Command } from '@ui/overlay/CommandMenu';
import { User03 } from '@ui/media/icons/User03.tsx';

const isContact = (
  item: { type: string; item: Contact } | { type: string; item: Organization },
): item is { item: Contact; type: 'contact' } => {
  return item.type === 'contact';
};

const isOrganization = (
  item: { type: string; item: Contact } | { type: string; item: Organization },
): item is { item: Organization; type: 'organization' } => {
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
      filteredContacts: Contact[];
      filteredOrgs: Organization['value'][];
    }>(
      (acc, result) => {
        if (isContact(result.item)) {
          acc.filteredContacts.push(result.item.item);
        }

        if (isOrganization(result.item)) {
          acc.filteredOrgs.push(
            result.item.item.value as Organization['value'],
          );
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
      <Command.Group heading={filteredContacts.length > 0 && 'Contacts'}>
        {filteredContacts?.map((contactStore) => (
          <Command.Item
            key={contactStore.value.id}
            onSelect={() =>
              handleGoTo(
                contactStore.value.primaryOrganizationId || '',
                'people',
              )
            }
            value={
              contactStore.name || contactStore.value.id
                ? `${contactStore.name} ${contactStore.value.id}`
                : ''
            }
          >
            <div className='flex items-center'>
              <Avatar
                size='xxs'
                textSize='xxs'
                name={contactStore.name}
                icon={<User03 className='text-primary-700  ' />}
                src={
                  contactStore?.value?.profilePhotoUrl
                    ? contactStore.value.profilePhotoUrl
                    : undefined
                }
              />
              <span className='ml-2 capitalize line-clamp-1 '>
                {contactStore.name}
              </span>

              <span className='ml-1.5 text-gray-500 line-clamp-1 max-w-[250px]'>
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

      <Command.Group heading={filteredOrgs.length > 0 && 'Organizations'}>
        {filteredOrgs?.map((org) => (
          <Command.Item
            key={org.id}
            value={`${org.name}-${org.id}`}
            onSelect={() => handleGoTo(org.id, 'about')}
          >
            <div className='flex items-center'>
              <Avatar
                size='xxs'
                textSize='xs'
                name={org.name}
                className='mr-2'
                variant='roundedSquare'
                src={org.iconUrl || org.logoUrl || undefined}
              />

              {org.name}
            </div>
          </Command.Item>
        ))}
      </Command.Group>
    </>
  );
});
