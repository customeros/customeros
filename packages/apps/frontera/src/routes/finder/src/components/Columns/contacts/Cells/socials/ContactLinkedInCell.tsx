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

    const enrichingStatus = contact?.isEnriching;

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit, store.ui]);

    const handleAddSocial = (url: string) => {
      if (!contact || url === 'Unknown' || url === '') return;
      const formattedValue =
        url.includes('https://www') || url.includes('linkedin.com')
          ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
          : `company/${url}`;

      contact.value.linkedInUrl = formattedValue;

      contact.commit();
      setIsEdit(false);
    };

    const toggleEditMode = () => setIsEdit(!isEdit);

    const linkedIn = contact?.value.linkedInUrl;

    if (!linkedIn) {
      return (
        <LinkedInInput
          type='in'
          isEdit={isEdit}
          metaKey={metaKey}
          isHovered={isHovered}
          setIsEdit={setIsEdit}
          enrichedStatus={false}
          setMetaKey={setMetaKey}
          setIsHovered={setIsHovered}
          handleAddSocial={handleAddSocial}
        />
      );
    }

    return (
      <LinkedInDisplay
        type={'in'}
        isEdit={isEdit}
        link={linkedIn}
        metaKey={metaKey}
        isHovered={isHovered}
        setIsEdit={setIsEdit}
        setMetaKey={setMetaKey}
        setIsHovered={setIsHovered}
        toggleEditMode={toggleEditMode}
        // handleUpdateSocial={handleUpdateSocial}
        alias={contact.value.linkedInAlias || ''}
      />
    );
  },
);
