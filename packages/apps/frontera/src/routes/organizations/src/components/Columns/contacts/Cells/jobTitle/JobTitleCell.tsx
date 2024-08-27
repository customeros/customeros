import { useRef, useState, useEffect, KeyboardEvent } from 'react';

import set from 'lodash/set';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';

interface JobTitleCellProps {
  contactId: string;
}

export const JobTitleCell = observer(({ contactId }: JobTitleCellProps) => {
  const store = useStore();
  const jobTitleInputRef = useRef<HTMLInputElement | null>(null);
  const [isHovered, setIsHovered] = useState(false);

  const contactStore = store.contacts.value.get(contactId);
  const jobTitle = contactStore?.value.jobRoles?.[0]?.jobTitle ?? '';

  const [isEdit, setIsEdit] = useState(false);
  const ref = useRef(null);

  useOutsideClick({
    ref: ref,
    handler: () => {
      setIsEdit(false);
    },
  });

  useEffect(() => {
    if (isHovered && isEdit) {
      jobTitleInputRef.current?.focus();
    }
  }, [isHovered, isEdit]);

  useEffect(() => {
    store.ui.setIsEditingTableCell(isEdit);
  }, [isEdit]);

  const handleEscape = (e: KeyboardEvent<HTMLDivElement>) => {
    if (e.key === 'Escape' || e.key === 'Enter') {
      jobTitleInputRef?.current?.blur();
      setIsEdit(false);
    }
  };

  return (
    <div
      ref={ref}
      onKeyDown={handleEscape}
      className='flex justify-between'
      onDoubleClick={() => setIsEdit(true)}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <div className='flex ' style={{ width: `calc(100% - 1rem)` }}>
        {!isEdit && !jobTitle && <p className='text-gray-400'>Unknown</p>}
        {!isEdit && jobTitle && (
          <p className='overflow-ellipsis overflow-hidden'>{jobTitle}</p>
        )}
        {isEdit && (
          <Input
            size='xs'
            variant='unstyled'
            ref={jobTitleInputRef}
            value={jobTitle ?? ''}
            placeholder='Job Title'
            onFocus={(e) => e.target.select()}
            onBlur={() => {
              if (!contactStore?.value?.jobRoles?.[0]?.id) {
                contactStore?.addJobRole();
              } else {
                contactStore?.updateJobRole();
              }
            }}
            onChange={(e) => {
              contactStore?.update(
                (value) => {
                  set(value, 'jobRoles[0].jobTitle', e.target.value);

                  return value;
                },
                { mutate: false },
              );
            }}
          />
        )}
        {isHovered && !isEdit && (
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label='edit'
            className='ml-3 rounded-[5px]'
            onClick={() => setIsEdit(!isEdit)}
            icon={<Edit03 className='text-gray-500' />}
          />
        )}
      </div>
    </div>
  );
});
