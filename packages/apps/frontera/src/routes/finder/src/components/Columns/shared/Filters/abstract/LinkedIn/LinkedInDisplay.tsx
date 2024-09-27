import { useRef, useEffect, KeyboardEvent } from 'react';

import { Input } from '@ui/form/Input';
import { Edit03 } from '@ui/media/icons/Edit03';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface LinkedInDisplayProps {
  link: string;
  type: string;
  alias?: string;
  isEdit: boolean;
  metaKey: boolean;
  isHovered: boolean;
  toggleEditMode: () => void;
  setIsEdit: (value: boolean) => void;
  setMetaKey: (value: boolean) => void;
  setIsHovered: (value: boolean) => void;
  handleUpdateSocial: (url: string) => void;
}

export const LinkedInDisplay = ({
  isHovered,
  alias,
  isEdit,
  setIsHovered,
  setIsEdit,
  handleUpdateSocial,
  metaKey,
  link,
  setMetaKey,
  toggleEditMode,
  type,
}: LinkedInDisplayProps) => {
  const inputRef = useRef<HTMLInputElement>(null);

  useOutsideClick({
    ref: inputRef,
    handler: () => {
      setIsEdit(false);
    },
  });

  const handleKeyEvents = (e: KeyboardEvent) => {
    if (e.key === 'Enter') {
      inputRef.current?.blur();
      setIsEdit(false);
    }

    if (e.key === 'Escape') {
      setIsEdit(false);
    }
  };

  useEffect(() => {
    if (isEdit) {
      inputRef?.current?.focus();
    }
  }, [isEdit]);

  const handleBlur = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleUpdateSocial(e.target.value);
  };
  const formattedLink = getFormattedLink(link).replace(
    /^linkedin\.com\/(?:in\/|company\/)?/,
    '/',
  );

  const displayLink = alias ? `/${alias}` : formattedLink;
  const url = formattedLink
    ? link.includes('linkedin')
      ? getExternalUrl(`https://linkedin.com/${type}${displayLink}`)
      : getExternalUrl(link)
    : '';

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
          value={link || ''}
          onChange={handleBlur}
          onKeyDown={handleKeyEvents}
          onBlur={() => setIsEdit(false)}
          onFocus={(e) => {
            displayLink
              ? handleUpdateSocial(`linkedin.com/${type}${displayLink}`)
              : handleUpdateSocial('');
            e.target.focus();
            setTimeout(() => e.target.select(), 0);
          }}
        />
      ) : (
        <Tooltip label={url ?? ''}>
          <p
            onDoubleClick={toggleEditMode}
            onClick={() => metaKey && toggleEditMode()}
            onKeyUp={() => metaKey && setMetaKey(false)}
            onKeyDown={(e) => e.metaKey && setMetaKey(true)}
            className='text-gray-700 cursor-default truncate'
          >
            {displayLink}
          </p>
        </Tooltip>
      )}
      {isHovered && !isEdit && (
        <>
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label='edit'
            onClick={toggleEditMode}
            className='ml-3 rounded-[5px]'
            icon={<Edit03 className='text-gray-500' />}
          />
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label='contact website'
            className='ml-1 rounded-[5px]'
            icon={<LinkExternal02 className='text-gray-500' />}
            onClick={() => window.open(url, '_blank', 'noopener')}
          />
        </>
      )}
    </div>
  );
};
