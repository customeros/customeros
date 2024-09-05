import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { FlowSequenceStatus } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { statusOptions } from './util';

interface FlowStatusCellProps {
  id: string;
}

export const FlowStatusCell = observer(({ id }: FlowStatusCellProps) => {
  const store = useStore();
  const [isEditing, setIsEditing] = useState(false);

  const flowSequence = store.flowSequences.value.get(id);

  const value = statusOptions.find(
    (option) => option.value === flowSequence?.value.status,
  );

  const handleSelect = (option: SelectOption<FlowSequenceStatus>) => {
    flowSequence?.update((seq) => {
      seq.status = option.value;

      return seq;
    });
    setIsEditing(false);
  };

  return (
    <div className='flex gap-1 items-center group/relationship'>
      <p
        onDoubleClick={() => setIsEditing(true)}
        data-test='organization-relationship-in-all-orgs-table'
        className={cn(
          'cursor-default text-gray-700',
          !value && 'text-gray-400',
        )}
      >
        {value?.label ?? 'No status'}
      </p>
      <Menu open={isEditing} onOpenChange={setIsEditing}>
        <MenuButton asChild>
          <IconButton
            size='xxs'
            variant='ghost'
            id='edit-button'
            aria-label='edit relationship'
            onClick={() => setIsEditing(true)}
            icon={<Edit03 className='text-gray-500' />}
            data-test='organization-relationship-button-in-all-orgs-table'
            className={cn(
              'rounded-md opacity-0 group-hover/relationship:opacity-100 min-w-5',
              isEditing && 'opacity-100',
            )}
          />
        </MenuButton>
        <MenuList>
          {statusOptions
            .filter((e) => e.value !== FlowSequenceStatus.Archived)
            .map((option) => (
              <MenuItem
                key={option.value.toString()}
                onClick={() => handleSelect(option)}
                data-test={`relationship-${option.value}`}
              >
                {option.label}
              </MenuItem>
            ))}
        </MenuList>
      </Menu>
    </div>
  );
});
