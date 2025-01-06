import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

import { LinkedInDisplay } from '../../../shared/Filters/abstract/LinkedIn';

interface SocialsCellProps {
  contactId: string;
}

export const ContactLinkedInCell = observer(
  ({ contactId }: SocialsCellProps) => {
    const store = useStore();
    const contact = store.contacts.value.get(contactId);

    if (!contact) return null;

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
        link={linkedIn || ''}
        alias={contact.value.linkedInAlias || ''}
      />
    );
  },
);
