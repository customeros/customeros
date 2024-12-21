import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

import { LinkedInDisplay } from '../../../shared/Filters/abstract/LinkedIn';

interface SocialsCellProps {
  contactId: string;
}

export const ContactLinkedInCell = observer(
  ({ contactId }: SocialsCellProps) => {
    const store = useStore();
    const [isHovered, setIsHovered] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const contact = store.contacts.value.get(contactId);
    const [metaKey, setMetaKey] = useState(false);

    if (!contact) return null;

    const toggleEditMode = () => setIsEdit(!isEdit);

    const linkedIn = contact?.value.linkedInUrl;

    if (!linkedIn) {
      return (
        <p
          className='text-sm text-gray-400 cursor-pointer'
          onClick={() => {
            store.ui.commandMenu.setType('AddLinkedin');
            store.ui.commandMenu.setOpen(true);
          }}
        >
          Not set
        </p>
      );
    }

    return (
      <LinkedInDisplay
        type={'in'}
        isEdit={isEdit}
        metaKey={metaKey}
        link={linkedIn || ''}
        isHovered={isHovered}
        setIsEdit={setIsEdit}
        setMetaKey={setMetaKey}
        setIsHovered={setIsHovered}
        toggleEditMode={toggleEditMode}
        alias={contact.value.linkedInAlias || ''}
      />
    );
  },
);
