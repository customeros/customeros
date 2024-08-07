import React, { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';

import Fuse from 'fuse.js';
import { useCommandState } from 'cmdk';
import { observer } from 'mobx-react-lite';

import { Avatar } from '@ui/media/Avatar';
import { useStore } from '@shared/hooks/useStore';
import { Command } from '@ui/overlay/CommandMenu';
import { User03 } from '@ui/media/icons/User03.tsx';

export const GlobalSearchResultNavigationCommands = observer(() => {
  const search = useCommandState((state) => state.search);
  const navigate = useNavigate();

  const { contacts, organizations, ui } = useStore();
  const contactsList = useMemo(
    () => contacts.toArray(),
    [contacts.totalElements],
  );
  const orgList = useMemo(
    () => organizations.toArray(),
    [organizations.totalElements],
  );
  const fuseContact = useMemo(
    () =>
      new Fuse(contactsList, {
        keys: ['name'],
        threshold: 0.3,
      }),
    [],
  );
  const fuseOrg = useMemo(
    () =>
      new Fuse(orgList, {
        keys: ['value.name'],
        threshold: 0.3,
      }),
    [],
  );
  const filteredContacts = useMemo(() => {
    if (!search) return [];

    return fuseContact.search(search).map((result) => result.item);
  }, [search, fuseContact]);

  const filteredOrgs = useMemo(() => {
    if (!search) return [];

    return fuseOrg.search(search).map((result) => result.item.value);
  }, [search, fuseOrg]);

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
              <span className='ml-2'>{contactStore.value.name}</span>

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
                size='xs'
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
