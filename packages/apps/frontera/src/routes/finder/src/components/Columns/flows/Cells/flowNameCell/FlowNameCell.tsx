import { useNavigate } from 'react-router-dom';
import { useRef, useState, useEffect, KeyboardEvent } from 'react';

import set from 'lodash/set';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { TableCellTooltip } from '@ui/presentation/Table';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';

interface FlowNameCellProps {
  id: string;
}

export const FlowNameCell = observer(({ id }: FlowNameCellProps) => {
  const store = useStore();
  const ref = useRef<HTMLDivElement | null>(null);
  const nameInputRef = useRef<HTMLInputElement | null>(null);

  const [isEdit, setIsEdit] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const navigate = useNavigate();

  const flowStore = store.flows.value.get(id);
  const flowName = flowStore?.value?.name;

  const itemRef = useRef<HTMLDivElement>(null);

  useOutsideClick({
    ref: ref,
    handler: () => {
      setIsEdit(false);
    },
  });

  useEffect(() => {
    if (isHovered && isEdit) {
      nameInputRef.current?.focus();
    }
  }, [isHovered, isEdit]);

  useEffect(() => {
    store.ui.setIsEditingTableCell(isEdit);
  }, [isEdit]);

  const handleEscape = (e: KeyboardEvent<HTMLDivElement>) => {
    if (e.key === 'Escape' || e.key === 'Enter') {
      e.stopPropagation();
      nameInputRef?.current?.blur();
      setIsEdit(false);
    }
  };

  const handleNavigate = () => {
    navigate(`/flow-editor/${id}`);
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
        {!isEdit && !flowName && <p className='text-gray-400'>Unnamed</p>}
        {!isEdit && flowName && (
          <TableCellTooltip
            hasArrow
            align='start'
            side='bottom'
            label={flowName}
            targetRef={itemRef}
          >
            <div
              ref={itemRef}
              role='button'
              onClick={handleNavigate}
              className='overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer'
            >
              <div
                data-test='flow-name-in-flows-table'
                className=' overflow-x-hidden overflow-ellipsis font-medium'
              >
                {flowName}
              </div>
            </div>
          </TableCellTooltip>
        )}
        {isEdit && (
          <Input
            size='xs'
            variant='unstyled'
            ref={nameInputRef}
            className='min-h-5'
            placeholder='Flow name'
            onKeyDown={handleEscape}
            onFocus={(e) => e.target.select()}
            value={flowStore?.value?.name ?? ''}
            onChange={(e) => {
              flowStore?.update((value) => {
                set(value, 'name', e.target.value);

                return value;
              });
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
