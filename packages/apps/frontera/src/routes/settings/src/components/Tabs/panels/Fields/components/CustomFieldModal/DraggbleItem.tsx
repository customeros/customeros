import { Draggable } from '@hello-pangea/dnd';

import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { HandleDrag } from '@ui/media/icons/HandleDrag';

interface Option {
  id: string;
  value: string;
  label: string;
}

interface DraggableItemProps {
  index: number;
  option: Option;
  newOption: Option[];
  isHovered: string | null;
  setIsHovered: (id: string | null) => void;
  setnewOptions: (options: Option[]) => void;
}

export const DraggableItem = ({
  index,
  option,
  isHovered,
  setIsHovered,
  newOption,
  setnewOptions,
}: DraggableItemProps) => {
  return (
    <Draggable index={index} key={option.id} draggableId={option.id}>
      {(provided) => (
        <div
          key={option.id}
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
          className='flex relative mt-1'
          onMouseLeave={() => setIsHovered(null)}
          onMouseEnter={() => setIsHovered(option.id)}
        >
          <HandleDrag className='absolute bottom-2.5 left-[7px]' />
          <Input
            size='sm'
            variant='outline'
            value={option.label}
            placeholder='Option'
            id={`option-${index}`}
            className='my-0.5 pl-8'
            onChange={(e) => {
              const newOptions = [...newOption];

              newOptions[index] = {
                id: option.id,
                value: e.target.value,
                label: e.target.value,
              };
              setnewOptions(newOptions);
            }}
          />
          {isHovered === option.id && (
            <IconButton
              size='xxs'
              icon={<X />}
              variant='ghost'
              aria-label='delete option'
              className='absolute right-2 transform translate-y-[43%]'
              onClick={() => {
                newOption.splice(index, 1);
                setnewOptions([...newOption]);
              }}
            />
          )}
        </div>
      )}
    </Draggable>
  );
};
