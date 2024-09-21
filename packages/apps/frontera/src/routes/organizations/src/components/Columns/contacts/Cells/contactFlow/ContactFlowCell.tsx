import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store.ts';

import { cn } from '@ui/utils/cn.ts';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { TableCellTooltip } from '@ui/presentation/Table';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { Select, getContainerClassNames } from '@ui/form/Select';

interface ContactNameCellProps {
  contactId: string;
}

export const ContactFlowCell = observer(
  ({ contactId }: ContactNameCellProps) => {
    const store = useStore();
    const [isEditing, setIsEditing] = useState(false);

    const contactStore = store.contacts.value.get(contactId);
    const flowName = contactStore?.flow?.value?.name;
    const flowId = contactStore?.flow?.id;
    const itemRef = useRef<HTMLDivElement>(null);
    const flowOptions = store.flows.toComputedArray((arr) => arr);

    const [value, setValue] = useState(() =>
      flowName ? { label: flowName, value: flowId } : null,
    );

    const close = () => {
      setIsEditing(false);
      store.ui.setIsEditingTableCell(false);
    };

    const open = () => {
      setIsEditing(true);
      store.ui.setIsEditingTableCell(true);
    };

    if (!flowName && !isEditing)
      return (
        <div
          onDoubleClick={open}
          className={cn(
            'flex w-full gap-1 items-center [&_.edit-button]:hover:opacity-100',
          )}
        >
          <div className='text-gray-400'>None</div>
          <IconButton
            size='xxs'
            onClick={open}
            variant='ghost'
            id='edit-button'
            aria-label='edit owner'
            className='edit-button opacity-0'
            dataTest={`contact-flow-edit-${contactId}`}
            icon={<Edit03 className='text-gray-500 size-3' />}
          />
        </div>
      );

    if (!isEditing) {
      return (
        <TableCellTooltip
          hasArrow
          align='start'
          side='bottom'
          label={flowName}
          targetRef={itemRef}
        >
          <div
            onDoubleClick={open}
            className={cn(
              'cursor-default overflow-hidden overflow-ellipsis flex gap-1  [&_.edit-button]:hover:opacity-100',
            )}
          >
            <div ref={itemRef} className='flex overflow-hidden'>
              <div
                data-test='flow-name'
                className=' overflow-x-hidden overflow-ellipsis'
              >
                {flowName}
              </div>
            </div>
            <IconButton
              size='xxs'
              onClick={open}
              variant='ghost'
              id='edit-button'
              aria-label='edit owner'
              className='edit-button opacity-0'
              icon={<Edit03 className='text-gray-500 size-3' />}
            />
          </div>
        </TableCellTooltip>
      );
    }

    const handleSelect = (option: SelectOption) => {
      const targetFlow = store.flows.value.get(option?.value) as FlowStore;

      setValue(option);
      targetFlow.linkContact(contactId);
    };
    const filteredOptions = flowOptions
      ?.filter((flow) => flow.value.name)
      .map((flow) => ({
        value: flow.id,
        label: flow.value.name,
      }));

    return (
      <Select
        size='xs'
        autoFocus
        isClearable
        value={value}
        onBlur={close}
        defaultMenuIsOpen
        placeholder='Flow'
        dataTest='flow-name'
        backspaceRemovesValue
        openMenuOnClick={false}
        onChange={handleSelect}
        options={filteredOptions}
        menuPortalTarget={document.body}
        onKeyDown={(e) => {
          if (e.key === 'Escape') {
            close();
          }
        }}
        classNames={{
          container: ({ isFocused }) =>
            getContainerClassNames('border-0 w-[164px]', undefined, {
              isFocused,
              size: 'xs',
            }),
        }}
      />
    );
  },
);
