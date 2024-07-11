import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface WebsiteCellProps {
  organizationId: string;
}

export const WebsiteCell = observer(({ organizationId }: WebsiteCellProps) => {
  const store = useStore();
  const [isHovered, setIsHovered] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [metaKey, setMetaKey] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const organization = store.organizations.value.get(organizationId);

  useEffect(() => {
    if (isHovered && isEdit) {
      inputRef.current?.focus();
    }
  }, [isHovered, isEdit]);

  useEffect(() => {
    store.ui.setIsEditingTableCell(isEdit);
  }, [isEdit]);

  if (!organization?.value.website?.length)
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
          if (e.metaKey) setIsEdit(true);
        }}
        onBlur={() => setIsEdit(false)}
      >
        {!isEdit ? (
          <p
            className='text-gray-400'
            data-test='organization-website-in-all-orgs-table'
          >
            Unknown
          </p>
        ) : (
          <Input
            size='xs'
            ref={inputRef}
            placeholder='Unknown'
            variant='unstyled'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                inputRef.current?.blur();
              }
              if (e.key === 'Escape') {
                inputRef.current?.blur();
              }
            }}
            onBlur={(e) => {
              const value = e.target.value;
              if (!organization || value === 'Unknown' || value === '') return;
              organization.update((org) => {
                if (value.includes('https://www')) {
                  const newUrl = getFormattedLink(value);
                  org.website = newUrl;
                }
                org.website = value;

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
  const website = organization?.value.website;

  const formattedLink = getFormattedLink(website);

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
          placeholder='Unknown'
          variant='unstyled'
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              inputRef.current?.blur();
            }
            if (e.key === 'Escape') {
              inputRef.current?.blur();
            }
          }}
          value={formattedLink}
          onChange={(e) => {
            const value = e.target.value;
            organization.update((org) => {
              if (value.includes('https://www')) {
                const newUrl = getFormattedLink(value);
                org.website = newUrl;
              }
              org.website = value;

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
          {formattedLink || 'Unknown'}
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
              window.open(getExternalUrl(website ?? '/'), '_blank', 'noopener')
            }
            aria-label='organization website'
            icon={<LinkExternal02 className='text-gray-500' />}
          />
        </>
      )}
    </div>
  );
});
