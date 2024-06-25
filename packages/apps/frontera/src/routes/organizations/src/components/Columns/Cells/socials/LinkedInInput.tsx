import { useRef, FocusEvent, KeyboardEvent } from 'react';

import { Input } from '@ui/form/Input';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton/IconButton';

interface LinkedInInputProps {
  isEdit: boolean;
  metaKey: boolean;
  isHovered: boolean;
  setIsEdit: (value: boolean) => void;
  setMetaKey: (value: boolean) => void;
  setIsHovered: (value: boolean) => void;
  handleAddSocial: (url: string) => void;
}

export const LinkedInInput = ({
  isHovered,
  isEdit,
  setIsHovered,
  setIsEdit,
  handleAddSocial,
  metaKey,
  setMetaKey,
}: LinkedInInputProps) => {
  const inputRef = useRef<HTMLInputElement>(null);

  const handleBlur = (e: FocusEvent<HTMLInputElement>) => {
    handleAddSocial(e.target.value);
  };

  const handleKeyEvents = (e: KeyboardEvent) => {
    if (e.key === 'Enter') {
      inputRef.current?.blur();
    }
    if (e.key === 'Escape') {
      setIsEdit(false);
    }
  };

  return (
    <div
      className='flex items-center'
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onDoubleClick={() => setIsEdit(true)}
      onKeyDown={(e) => e.metaKey && setMetaKey(true)}
      onKeyUp={() => metaKey && setMetaKey(false)}
      onClick={() => metaKey && setIsEdit(true)}
      onBlur={() => setIsEdit(false)}
    >
      {!isEdit ? (
        <p className='text-gray-400'>Unknown</p>
      ) : (
        <Input
          size='xs'
          ref={inputRef}
          variant='unstyled'
          onKeyDown={handleKeyEvents}
          onBlur={handleBlur}
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
};
