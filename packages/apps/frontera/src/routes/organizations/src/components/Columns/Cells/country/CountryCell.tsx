import React from 'react';

import { observer } from 'mobx-react-lite';
import { ContactsStore } from '@store/Contacts/Contacts.store';
import { OrganizationsStore } from '@store/Organizations/Organizations.store';

import { flags } from '@ui/media/flags';
import { useStore } from '@shared/hooks/useStore';

interface ContactNameCellProps {
  id: string;
  type?: 'contact' | 'organization';
}

export const CountryCell: React.FC<ContactNameCellProps> = observer(
  ({ id, type }) => {
    const { organizations, contacts } = useStore();
    const store: ContactsStore | OrganizationsStore =
      type === 'contact' ? contacts : organizations;
    const itemStore = store.value.get(id);
    const country = itemStore?.country;

    if (!country) {
      return <div className='text-gray-400'>Unknown</div>;
    }
    const alpha2 = itemStore?.value?.locations?.[0]?.countryCodeA2;

    return (
      <div className='flex items-center'>
        <div className='flex items-center'>{alpha2 && flags[alpha2]}</div>
        <span className='ml-2 overflow-hidden overflow-ellipsis whitespace-nowrap'>
          {country}
        </span>
      </div>
    );
  },
);
