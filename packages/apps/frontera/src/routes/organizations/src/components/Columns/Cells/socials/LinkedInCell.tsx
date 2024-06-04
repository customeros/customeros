import { useRef, useState } from 'react';

import { Input } from '@ui/form/Input';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02.tsx';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface SocialsCellProps {
  organizationId: string;
}

export const LinkedInCell = ({ organizationId }: SocialsCellProps) => {
  const store = useStore();
  const [isHovered, setIsHovered] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const organization = store.organizations.value.get(organizationId);

  if (!organization?.value.socialMedia?.length)
    return <p className='text-gray-400'>Unknown</p>;
  const linkedIn = organization?.value.socialMedia.find((social) =>
    social?.url?.includes('linkedin'),
  );

  if (!linkedIn?.url) return;

  const formattedLink = getFormattedLink(linkedIn.url).replace(
    /^linkedin\.com\//,
    '',
  );
  const linkedinId = linkedIn?.id;

  return (
    <div
      className='flex items-center'
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {isEdit ? (
        <Input
          size='sm'
          ref={inputRef}
          variant='unstyled'
          value={formattedLink}
          onChange={(e) => {
            const value = e.target.value;
            if (!organization) return null;
            organization.update((org) => {
              const idx = organization?.value.socialMedia.findIndex(
                (s) => s.id === linkedinId,
              );

              if (idx !== -1) {
                org.socialMedia[idx].url = `linkedin.com/${value}`;
              }

              return org;
            });
          }}
          onBlur={() => setIsEdit(false)}
        />
      ) : (
        <p
          className='text-gray-700 cursor-default truncate'
          onDoubleClick={() => setIsEdit(true)}
        >
          {formattedLink}
        </p>
      )}
      {isHovered && !isEdit && (
        <>
          <IconButton
            className='ml-3 rounded-[5px]'
            variant='ghost'
            size='xxs'
            onClick={() => setIsEdit(!isEdit)}
            aria-label='edit'
            icon={<Edit03 className='text-gray-500' />}
          />
          <IconButton
            className='ml-1 rounded-[5px]'
            variant='ghost'
            size='xxs'
            onClick={() =>
              window.open(
                getExternalUrl(linkedIn.url ?? '/'),
                '_blank',
                'noopener',
              )
            }
            aria-label='organization website'
            icon={<LinkExternal02 className='text-gray-500' />}
          />
        </>
      )}
    </div>
  );
};
