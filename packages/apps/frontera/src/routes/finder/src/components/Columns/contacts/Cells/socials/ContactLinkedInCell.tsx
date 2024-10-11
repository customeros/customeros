import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getFormattedLink } from '@utils/getExternalLink';
import { Social } from '@shared/types/__generated__/graphql.types';

import {
  LinkedInInput,
  LinkedInDisplay,
} from '../../../shared/Filters/abstract/LinkedIn';

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

    const enrichedContact = contact?.value.enrichDetails;
    const enrichingStatus =
      enrichedContact?.requestedAt &&
      !enrichedContact?.failedAt &&
      !enrichedContact?.enrichedAt;

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit, store.ui]);

    const handleAddSocial = (url: string) => {
      if (!contact || url === 'Unknown' || url === '') return;

      contact.update((contactData) => {
        const formattedValue =
          url.includes('https://www') || url.includes('linkedin.com')
            ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
            : `company/${url}`;

        contactData.socials.push({
          id: crypto.randomUUID(),
          url: `linkedin.com/${formattedValue}`,
        } as Social);

        return contactData;
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
          type='in'
          isEdit={isEdit}
          metaKey={metaKey}
          isHovered={isHovered}
          setIsEdit={setIsEdit}
          setMetaKey={setMetaKey}
          setIsHovered={setIsHovered}
          enrichedStatus={enrichingStatus}
          handleAddSocial={handleAddSocial}
        />
      );
    }

    return (
      <LinkedInDisplay
        type={'in'}
        isEdit={isEdit}
        metaKey={metaKey}
        link={linkedIn.url}
        isHovered={isHovered}
        setIsEdit={setIsEdit}
        alias={linkedIn.alias}
        setMetaKey={setMetaKey}
        setIsHovered={setIsHovered}
        toggleEditMode={toggleEditMode}
        handleUpdateSocial={handleUpdateSocial}
      />
    );
  },
);
