import { useState, ReactElement } from 'react';

import { NodeProps, useReactFlow } from '@xyflow/react';

import { Input } from '@ui/form/Input';
import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { Textarea } from '@ui/form/Textarea/Textarea';
import { MailReply } from '@ui/media/icons/MailReply';
import { Hourglass01 } from '@ui/media/icons/Hourglass01';
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
  const { setNodes } = useReactFlow();

  const color = colorMap?.[data.stepType];

  const handleEmailDataChange = ({
    newValue,
    property,
  }: {
    newValue: string;
    property: 'subject' | 'bodyTemplate';
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
        className={`aspect-[9/1] max-w-[300px] bg-white border border-grayModern-300 p-3 rounded-lg group`}
      >
        <div className='truncate  text-sm flex items-center '>
          <div className='truncate text-sm flex items-center'>
            <div
              className={`size-6 mr-2 bg-${color}-50 text-${color}-500 border border-${color}-100  rounded flex items-center justify-center`}
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
                    value={data?.bodyTemplate ?? ''}
                    onChange={(e) =>
                      handleEmailDataChange({
                        newValue: e.target.value,
                        property: 'bodyTemplate',
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
