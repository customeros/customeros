import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getFormattedLink } from '@utils/getExternalLink';
import { Social } from '@shared/types/__generated__/graphql.types';

import { LinkedInInput } from './LinkedInInput.tsx';
import { LinkedInDisplay } from './LinkedInDisplay.tsx';

interface SocialsCellProps {
  organizationId: string;
}

export const OrganizationLinkedInCell = observer(
  ({ organizationId }: SocialsCellProps) => {
    const store = useStore();
    const [isHovered, setIsHovered] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const organization = store.organizations.value.get(organizationId);
    const [metaKey, setMetaKey] = useState(false);

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit, store.ui]);

    const handleAddSocial = (url: string) => {
      if (!organization || url === 'Unknown' || url === '') return;
      organization.update((org) => {
        const formattedValue =
          url.includes('https://www') || url.includes('linkedin.com')
            ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
            : `in/${url}`;
        org.socialMedia.push({
          id: crypto.randomUUID(),
          url: `linkedin.com/${formattedValue}`,
        } as Social);

        return org;
      });
      setIsEdit(false);
    };

    const handleUpdateSocial = (url: string) => {
      const linkedinId = organization?.value.socialMedia.find((social) =>
        social.url.includes('linkedin'),
      )?.id;
      if (!linkedinId) return;

      organization.update((org) => {
        const idx = org.socialMedia.findIndex((s) => s.id === linkedinId);
        if (idx !== -1) {
          const formattedValue =
            url.includes('https://www') || url.includes('linkedin.com')
              ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
              : `in/${url}`;
          org.socialMedia[idx].url = `linkedin.com/${formattedValue}`;
        }

        if (url === '') {
          org.socialMedia.splice(idx, 1);
        }

        return org;
      });
    };

    const toggleEditMode = () => setIsEdit(!isEdit);

    if (!organization?.value.socialMedia?.length) {
      return (
        <LinkedInInput
          isHovered={isHovered}
          isEdit={isEdit}
          setIsHovered={setIsHovered}
          setIsEdit={setIsEdit}
          handleAddSocial={handleAddSocial}
          metaKey={metaKey}
          setMetaKey={setMetaKey}
          type='company'
        />
      );
    }

    const linkedIn = organization.value.socialMedia.find((social) =>
      social.url.includes('linkedin'),
    );
    if (!linkedIn?.url) return null;

    return (
      <LinkedInDisplay
        isHovered={isHovered}
        isEdit={isEdit}
        type='company'
        setIsHovered={setIsHovered}
        setIsEdit={setIsEdit}
        link={linkedIn.url}
        handleUpdateSocial={handleUpdateSocial}
        metaKey={metaKey}
        setMetaKey={setMetaKey}
        toggleEditMode={toggleEditMode}
      />
    );
  },
);
