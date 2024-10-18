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
  organizationId: string;
}

export const OrganizationLinkedInCell = observer(
  ({ organizationId }: SocialsCellProps) => {
    const store = useStore();
    const [isHovered, setIsHovered] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const [metaKey, setMetaKey] = useState(false);
    const organization = store.organizations.value.get(organizationId);

    const enrichedOrganizations = organization?.value.enrichDetails;

    const enrichingStatus =
      !enrichedOrganizations?.enrichedAt &&
      enrichedOrganizations?.requestedAt &&
      !enrichedOrganizations?.failedAt;

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit, store.ui]);

    const handleAddSocial = (url: string) => {
      if (!organization || url === 'Unknown' || url === '') return;

      const formattedValue =
        url.includes('https://www') || url.includes('linkedin.com')
          ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
          : `in/${url}`;

      organization.value.socialMedia.push({
        id: crypto.randomUUID(),
        url: `linkedin.com/${formattedValue}`,
      } as Social);

      organization.commit();

      setIsEdit(false);
    };

    const handleUpdateSocial = (url: string) => {
      const linkedinId = organization?.value.socialMedia.find((social) =>
        social.url.includes('linkedin'),
      )?.id;

      if (!linkedinId) return;

      const idx = organization.value.socialMedia.findIndex(
        (s) => s.id === linkedinId,
      );

      if (idx !== -1) {
        const formattedValue =
          url.includes('https://www') || url.includes('linkedin.com')
            ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
            : `in/${url}`;

        organization.value.socialMedia[
          idx
        ].url = `linkedin.com/${formattedValue}`;
      }

      if (url === '') {
        organization.value.socialMedia.splice(idx, 1);
      }

      organization.commit();
    };

    const toggleEditMode = () => setIsEdit(!isEdit);
    const linkedIn = organization?.value?.socialMedia.find((social) =>
      social.url.includes('linkedin'),
    );

    if (!organization?.value.socialMedia?.length || !linkedIn) {
      return (
        <LinkedInInput
          type='company'
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
        type='company'
        isEdit={isEdit}
        metaKey={metaKey}
        link={linkedIn.url}
        isHovered={isHovered}
        setIsEdit={setIsEdit}
        setMetaKey={setMetaKey}
        setIsHovered={setIsHovered}
        toggleEditMode={toggleEditMode}
        handleUpdateSocial={handleUpdateSocial}
      />
    );
  },
);
