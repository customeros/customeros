import { useState } from 'react';
import { ValueContainerProps } from 'react-select';

import { Droppable, DragDropContext } from '@hello-pangea/dnd';

import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { Type01 } from '@ui/media/icons/Type01';
import { Hash02 } from '@ui/media/icons/Hash02';
import { IconButton } from '@ui/form/IconButton';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { RadioButton } from '@ui/media/icons/RadioButton';
import { ListBulleted } from '@ui/media/icons/ListBulleted';
import {
  Select,
  components,
  OptionProps,
  getContainerClassNames,
} from '@ui/form/Select';
import {
  Modal,
  ModalBody,
  ModalClose,
  ModalPortal,
  ModalFooter,
  ModalOverlay,
  ModalCloseButton,
  ModalFeaturedHeader,
  ModalFeaturedContent,
} from '@ui/overlay/Modal';

interface NewCustomFieldModalProps {
  title: string;
  isOpen: boolean;
  onOpenChange: (value: boolean) => void;
}

const options = [
  {
    label: 'Text',
    id: 'text',
    icon: <Type01 />,
  },
  {
    label: 'Number',
    id: 'number',
    icon: <Hash02 />,
  },
  {
    label: 'Single select',
    id: 'single-select',
    icon: <RadioButton />,
  },
  {
    label: 'Multi select',
    id: 'multi-select',
    icon: <ListBulleted />,
  },
];

export const NewCustomFieldModal = ({
  isOpen,
  title,
  onOpenChange,
}: NewCustomFieldModalProps) => {
  const [selectedOption, setSelectedOption] = useState(options[0]);
  const [isHovered, setIsHovered] = useState<string | null>(null);

  const [newOption, setnewOptions] = useState<
    { id: string; value: string; label: string }[]
  >([]);

  const Option = ({ children, ...props }: OptionProps) => {
    const option = props.data as {
      id: string;
      label: string;
      icon: JSX.Element;
    };

    return (
      <components.Option {...props}>
        <div className='flex items-center gap-2'>
          {option.icon}
          {children}
        </div>
      </components.Option>
    );
  };

  const ValueContainer = ({ children, ...props }: ValueContainerProps) => {
    const selectedValue = props?.getValue()[0] as unknown as {
      id: string;
      label: string;
      icon: JSX.Element;
    };

    const icon = options.find(
      (option) => option?.id === selectedValue?.id,
    )?.icon;

    return (
      <components.ValueContainer {...props}>
        <div className='flex items-center gap-2'>
          {icon}
          {children}
        </div>
      </components.ValueContainer>
    );
  };

  return (
    <Modal open={isOpen} onOpenChange={(value) => onOpenChange(value)}>
      <ModalPortal>
        <ModalOverlay className='z-[999]' />
        <ModalFeaturedContent className='z-[9999]'>
          <ModalFeaturedHeader>
            <p className='text-lg font-semibold mb-1'>{title}</p>
            <ModalCloseButton asChild />
          </ModalFeaturedHeader>
          <ModalCloseButton />
          <ModalBody className='flex flex-col gap-4'>
            <div className='flex flex-col gap-2'>
              <div>
                <label htmlFor='type' className='font-medium'>
                  Type
                </label>
                <Select
                  id='type'
                  autoFocus={false}
                  defaultValue={options[0]}
                  components={{ Option, ValueContainer }}
                  onChange={(value) => {
                    setSelectedOption(value as typeof selectedOption);
                  }}
                  options={options.map((option) => ({
                    ...option,
                    value: option.id,
                  }))}
                  classNames={{
                    container: (props) =>
                      getContainerClassNames('', 'outline', {
                        ...props,
                        size: 'sm',
                      }),
                  }}
                />
              </div>
              <div>
                <label htmlFor='name' className='font-medium'>
                  Name
                </label>
                <Input
                  id='name'
                  size='sm'
                  variant='outline'
                  placeholder='Custom field name'
                />
              </div>
              {selectedOption.id === 'single-select' && (
                <div className='mt-2 flex flex-col'>
                  <label htmlFor='options' className='font-medium'>
                    Options
                  </label>
                  <DragDropContext onDragEnd={() => {}}>
                    <Droppable>
                      {(provided) => (
                        <>
                          {newOption.map((option, index) => (
                            <div
                              key={index}
                              className='flex relative'
                              onMouseLeave={() => setIsHovered(null)}
                              onMouseEnter={() => setIsHovered(option.id)}
                            >
                              <Input
                                size='sm'
                                className='mb-1'
                                variant='outline'
                                value={option.label}
                                placeholder='Option'
                                id={`option-${index}`}
                                onChange={(e) => {
                                  const newOptions = [...newOption];

                                  newOptions[index] = {
                                    id: `option-${crypto.randomUUID()}`,
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
                                  className='absolute right-2 transform translate-y-1/3'
                                  onClick={() => {
                                    newOption.splice(index, 1);
                                    setnewOptions([...newOption]);
                                  }}
                                />
                              )}
                            </div>
                          ))}
                        </>
                      )}
                    </Droppable>
                  </DragDropContext>

                  <Button
                    size='xs'
                    variant='ghost'
                    className='w-fit'
                    leftIcon={<PlusCircle />}
                    onClick={() =>
                      setnewOptions([
                        ...newOption,
                        {
                          value: '',
                          label: '',
                          id: `option-${crypto.randomUUID()}`,
                        },
                      ])
                    }
                  >
                    Add option
                  </Button>
                </div>
              )}
            </div>
          </ModalBody>
          <ModalFooter className='flex gap-3'>
            <ModalClose className='w-full'>
              <Button className='w-full'>Cancel</Button>
            </ModalClose>
            <Button className='w-full' colorScheme='primary'>
              Create field
            </Button>
          </ModalFooter>
        </ModalFeaturedContent>
      </ModalPortal>
    </Modal>
  );
};
