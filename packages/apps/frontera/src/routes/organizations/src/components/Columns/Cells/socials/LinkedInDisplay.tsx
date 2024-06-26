import { useRef, KeyboardEvent } from 'react';

import { Input } from '@ui/form/Input';
import { Edit03 } from '@ui/media/icons/Edit03';
import { getExternalUrl } from '@utils/getExternalLink';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';

interface LinkedInDisplayProps {
  isEdit: boolean;
  metaKey: boolean;
  isHovered: boolean;
  formattedLink: string;
  toggleEditMode: () => void;
  setIsEdit: (value: boolean) => void;
  setMetaKey: (value: boolean) => void;
  setIsHovered: (value: boolean) => void;
  handleUpdateSocial: (url: string) => void;
}

export const LinkedInDisplay = ({
  isHovered,
  isEdit,
  setIsHovered,
  setIsEdit,
  formattedLink,
  handleUpdateSocial,
  metaKey,
  setMetaKey,
  toggleEditMode,
}: LinkedInDisplayProps) => {
  const inputRef = useRef<HTMLInputElement>(null);

  const handleKeyEvents = (e: KeyboardEvent) => {
    if (e.key === 'Enter') {
      inputRef.current?.blur();
    }
    if (e.key === 'Escape') {
      setIsEdit(false);
    }
  };

  const handleBlur = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleUpdateSocial(e.target.value);
  };

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
          onKeyDown={handleKeyEvents}
          onBlur={() => setIsEdit(false)}
          onChange={handleBlur}
        />
      ) : (
        <p
          className='text-gray-700 cursor-default truncate'
          onDoubleClick={toggleEditMode}
          onKeyDown={(e) => e.metaKey && setMetaKey(true)}
          onKeyUp={() => metaKey && setMetaKey(false)}
          onClick={() => metaKey && toggleEditMode()}
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
            onClick={toggleEditMode}
            aria-label='edit'
            icon={<Edit03 className='text-gray-500' />}
          />
          <IconButton
            className='ml-1 rounded-[5px]'
            variant='ghost'
            size='xxs'
            onClick={() =>
              window.open(
                getExternalUrl(`https://linkedin.com/${formattedLink}`),
                '_blank',
                'noopener',
              )
            }
            aria-label='contact website'
            icon={<LinkExternal02 className='text-gray-500' />}
          />
        </>
      )}
    </div>
  );
};
