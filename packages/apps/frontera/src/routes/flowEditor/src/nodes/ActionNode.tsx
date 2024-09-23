import { useState, ReactElement } from 'react';

import { NodeProps, useReactFlow } from '@xyflow/react';

import { IconButton } from '@ui/form/IconButton';
import { Check } from '@ui/media/icons/Check.tsx';
import { Button } from '@ui/form/Button/Button.tsx';
import { Mail01 } from '@ui/media/icons/Mail01.tsx';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { XSquare } from '@ui/media/icons/XSquare.tsx';
import { Input, ResizableInput } from '@ui/form/Input';
import { Textarea } from '@ui/form/Textarea/Textarea.tsx';
import { MailReply } from '@ui/media/icons/MailReply.tsx';
import { Hourglass02 } from '@ui/media/icons/Hourglass02.tsx';
import { Hourglass01 } from '@ui/media/icons/Hourglass01.tsx';
import {
  Modal,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

import { Handle } from '../components';

const iconMap: Record<string, ReactElement> = {
  SendEmail: <Mail01 className='text-inherit' />,
  ReplyToEmail: <MailReply className='text-inherit' />,
  Wait: <Hourglass01 className='text-inherit' />,
};

const colorMap: Record<string, string> = {
  SendEmail: 'blue',
  ReplyToEmail: 'blue',
  Wait: 'gray',
};

export const ActionNode = ({
  id,
  data,
}: NodeProps & { data: Record<string, string | number> }) => {
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const [isFocused, setFocused] = useState(false);
  const { setNodes } = useReactFlow();

  if (data.stepType === 'Wait') {
    const handleDurationChange = (newValue: string) => {
      setNodes((nds) =>
        nds.map((node) => {
          if (node.id === id) {
            return {
              ...node,
              data: {
                ...node.data,
                waitDuration: parseInt(newValue, 10),
              },
            };
          }

          return node;
        }),
      );
    };

    return (
      <div
        className={`h-[56px] w-[150px] bg-white border-2 border-grayModern-300 p-3 rounded-lg group`}
      >
        <div className='truncate  text-sm flex items-center '>
          <div className='truncate text-sm flex items-center'>
            <div
              className={`size-6 mr-2 bg-gray-100 text-gray-500 border-gray-50  rounded flex items-center justify-center`}
            >
              <Hourglass02 />
            </div>

            {isFocused ? (
              <div className='flex mr-1 items-baseline'>
                <ResizableInput
                  min={1}
                  size='xs'
                  autoFocus
                  type='number'
                  value={data.waitDuration || 0}
                  onFocus={(e) => e.target.select()}
                  className=' min-w-2.5 min-h-0 max-h-4'
                  onChange={(e) => handleDurationChange(e.target.value)}
                />
                <span className='ml-1'>
                  {data.waitDuration === 1 ? 'day' : 'days'}
                </span>
              </div>
            ) : (
              <span className='truncate'>
                {data.waitDuration || 0}{' '}
                {data.waitDuration === 1 ? 'day' : 'days'}
              </span>
            )}
          </div>

          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Edit'
            onClick={() => setFocused(!isFocused)}
            icon={!isFocused ? <Edit03 /> : <XSquare />}
            className='ml-2 w-0 group-hover:w-auto opacity-0 group-hover:opacity-100 pointer-events-all'
          />
        </div>
        <Handle type='target' />
        <Handle type='source' />
      </div>
    );
  }

  const color = colorMap?.[data.stepType];

  const handleEmailDataChange = ({
    newValue,
    property,
  }: {
    newValue: string;
    property: string;
  }) => {
    setNodes((nds) =>
      nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              [property]: newValue,
            },
          };
        }

        return node;
      }),
    );
  };

  return (
    <>
      <div
        className={`aspect-[9/1] max-w-[300px] bg-white border-2 border-grayModern-300 p-3 rounded-lg group`}
      >
        <div className='truncate  text-sm flex items-center '>
          <div className='truncate text-sm flex items-center'>
            <div
              className={`size-6 mr-2 bg-${color}-100 text-${color}-500 border-${color}-50  rounded flex items-center justify-center`}
            >
              {iconMap?.[data.stepType]}
            </div>

            <span className='truncate font-medium'>
              {data.subject ? (
                data.subject
              ) : data.content ? (
                data.content
              ) : (
                <span className='text-gray-400 font-normal'>
                  Write an email that wows them
                </span>
              )}
            </span>
          </div>

          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Edit'
            icon={<Edit03 />}
            onClick={() => setIsEditorOpen(true)}
            className='ml-2 opacity-0 group-hover:opacity-100 pointer-events-all'
          />
        </div>
        <Modal open={isEditorOpen} onOpenChange={setIsEditorOpen}>
          <ModalPortal>
            <ModalOverlay className='z-50'>
              {/* width and height of A4 */}
              <ModalContent className='w-full h-full flex justify-center max-w-full top-0'>
                <div className='w-[570px]'>
                  <div className='flex justify-between mt-4'>
                    <div className=''></div>
                    <Button
                      variant='ghost'
                      leftIcon={<Check />}
                      onClick={() => setIsEditorOpen(false)}
                    >
                      Done
                    </Button>
                  </div>

                  <Input
                    variant='unstyled'
                    placeholder='Subject'
                    value={data?.subject ?? ''}
                    onChange={(e) =>
                      handleEmailDataChange({
                        newValue: e.target.value,
                        property: 'subject',
                      })
                    }
                  />

                  <Textarea
                    value={data?.content ?? ''}
                    onChange={(e) =>
                      handleEmailDataChange({
                        newValue: e.target.value,
                        property: 'content',
                      })
                    }
                  />
                </div>
              </ModalContent>
            </ModalOverlay>
          </ModalPortal>
        </Modal>
        <Handle type='target' />
        <Handle type='source' />
      </div>
    </>
  );
};
