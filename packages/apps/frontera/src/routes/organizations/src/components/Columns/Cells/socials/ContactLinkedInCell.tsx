import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getFormattedLink } from '@utils/getExternalLink';
import { Social } from '@shared/types/__generated__/graphql.types';

import { LinkedInInput } from './LinkedInInput.tsx';
import { LinkedInDisplay } from './LinkedInDisplay.tsx';

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

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit, store.ui]);

    const handleAddSocial = (url: string) => {
      if (!contact || url === 'Unknown' || url === '') return;
      contact.update((org) => {
        const formattedValue =
          url.includes('https://www') || url.includes('linkedin.com')
            ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
            : `company/${url}`;
        org.socials.push({
          id: crypto.randomUUID(),
          url: `linkedin.com/${formattedValue}`,
        } as Social);

        return org;
      });
      setIsEdit(false);
    };

    const handleUpdateSocial = (url: string) => {
      const linkedinId = contact?.value.socials.find((social) =>
        social.url.includes('linkedin'),
      )?.id;
      if (!linkedinId) return;

      contact.update((org) => {
        const idx = org.socials.findIndex((s) => s.id === linkedinId);
        if (idx !== -1) {
          const formattedValue =
            url.includes('https://www') || url.includes('linkedin.com')
              ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
              : `in/${url}`;
          org.socials[idx].url = `linkedin.com/${formattedValue}`;
        }

        if (url === '') {
          org.socials.splice(idx, 1);
        }

        return org;
      });
    };

    const toggleEditMode = () => setIsEdit(!isEdit);

    const linkedIn = contact?.value.socials.find((social) =>
      social.url.includes('linkedin'),
    );
    if (!contact?.value.socials?.length || !linkedIn) {
      return (
        <LinkedInInput
          isHovered={isHovered}
          isEdit={isEdit}
          setIsHovered={setIsHovered}
          setIsEdit={setIsEdit}
          handleAddSocial={handleAddSocial}
          metaKey={metaKey}
          type='in'
          setMetaKey={setMetaKey}
        />
      );
    }

    return (
      <LinkedInDisplay
        isHovered={isHovered}
        isEdit={isEdit}
        setIsHovered={setIsHovered}
        setIsEdit={setIsEdit}
        link={linkedIn.url}
        alias={linkedIn.alias}
        handleUpdateSocial={handleUpdateSocial}
        metaKey={metaKey}
        type={'in'}
        setMetaKey={setMetaKey}
        toggleEditMode={toggleEditMode}
      />
    );
  },
);
