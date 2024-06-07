import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Social } from '@shared/types/__generated__/graphql.types';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02.tsx';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface SocialsCellProps {
  organizationId: string;
}

export const LinkedInCell = observer(({ organizationId }: SocialsCellProps) => {
  const store = useStore();
  const [isHovered, setIsHovered] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const organization = store.organizations.value.get(organizationId);
  const [metaKey, setMetaKey] = useState(false);

  useEffect(() => {
    if (isHovered && isEdit) {
      inputRef.current?.focus();
    }
  }, [isHovered, isEdit]);

  useEffect(() => {
    store.ui.setIsEditingTableCell(isEdit);
  }, [isEdit]);

  if (!organization?.value.socialMedia?.length)
    return (
      <div
        className='flex items-center'
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        onDoubleClick={() => setIsEdit(true)}
        onKeyDown={(e) => {
          if (e.metaKey) {
            setMetaKey(true);
          }
        }}
        onKeyUp={() => metaKey && setMetaKey(false)}
        onClick={(e) => {
          if (e.metaKey) {
            setIsEdit(true);
          }
        }}
        onBlur={() => setIsEdit(false)}
      >
        {!isEdit ? (
          <p className='text-gray-400'>Unknown</p>
        ) : (
          <Input
            size='xs'
            ref={inputRef}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                inputRef.current?.blur();
              }
            }}
            variant='unstyled'
            onBlur={(e) => {
              const value = e.target.value;
              if (!organization || value === 'Unknown' || value === '') return;
              organization.update((org) => {
                if (
                  value.includes('https://www') ||
                  value.includes('linkedin.com')
                ) {
                  const newUrl = getFormattedLink(value).replace(
                    /^linkedin\.com\//,
                    '',
                  );

                  org.socialMedia.push({
                    id: crypto.randomUUID(),
                    url: `linkedin.com/${newUrl}`,
                  } as Social);
                } else {
                  org.socialMedia.push({
                    id: crypto.randomUUID(),
                    url: `linkedin.com/in/${value}`,
                  } as Social);
                }

                return org;
              });
              setIsEdit(false);
            }}
          />
        )}
        {isHovered && !isEdit && (
          <IconButton
            className='ml-3 rounded-[5px]'
            variant='ghost'
            size='xxs'
            onClick={() => setIsEdit(!isEdit)}
            aria-label='edit'
            icon={<Edit03 className='text-gray-500' />}
          />
        )}
      </div>
    );
  const linkedIn = organization?.value.socialMedia.find((social) =>
    social?.url?.includes('linkedin'),
  );

  if (!linkedIn?.url) return;

  const formattedLink = getFormattedLink(linkedIn.url).replace(
    /^linkedin\.com\/in\//,
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
          size='xs'
          ref={inputRef}
          variant='unstyled'
          value={formattedLink}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              inputRef.current?.blur();
            }
          }}
          onChange={(e) => {
            const value = e.target.value;

            organization.update((org) => {
              const idx = organization?.value.socialMedia.findIndex(
                (s) => s.id === linkedinId,
              );

              if (idx !== -1) {
                if (
                  value.includes('https://www.') ||
                  value.includes('linkedin.com')
                ) {
                  const newUrl = getFormattedLink(value).replace(
                    /^linkedin\.com\//,
                    '',
                  );

                  org.socialMedia[idx].url = `linkedin.com/${newUrl}`;
                } else {
                  org.socialMedia[idx].url = `linkedin.com/in/${value}`;
                }
              }

              if (value === '') {
                org.socialMedia.splice(idx, 1);
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
          onKeyDown={(e) => {
            if (e.metaKey) {
              setMetaKey(true);
            }
          }}
          onKeyUp={() => metaKey && setMetaKey(false)}
          onClick={(e) => {
            if (e.metaKey) setIsEdit(true);
          }}
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
});
