import { useMemo, useState, ReactElement } from 'react';

import { FlowActionType } from '@store/Flows/types.ts';
import { NodeProps, useReactFlow } from '@xyflow/react';

import { Input } from '@ui/form/Input';
import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { Editor } from '@ui/form/Editor/Editor.tsx';
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
  [FlowActionType.EMAIL_NEW]: <Mail01 className='text-inherit' />,
  [FlowActionType.EMAIL_REPLY]: <MailReply className='text-inherit' />,
  WAIT: <Hourglass01 className='text-inherit' />,
};

const colorMap: Record<string, string> = {
  [FlowActionType.EMAIL_NEW]: 'blue',
  [FlowActionType.EMAIL_REPLY]: 'blue',
  WAIT: 'gray',
};

export const ActionNode = ({
  id,
  data,
}: NodeProps & { data: Record<string, string | number> }) => {
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const { setNodes } = useReactFlow();

  const color = colorMap?.[data.action];
  const placeholder = useMemo(() => getRandomEmailPrompt(), []);

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
              {iconMap?.[data.action]}
            </div>

            <span className='truncate font-medium'>
              {data.subject ? (
                data.subject
              ) : data.bodyTemplate ? (
                data.bodyTemplate
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
                    size='lg'
                    variant='unstyled'
                    placeholder='Subject'
                    className='font-medium'
                    value={data?.subject ?? ''}
                    onChange={(e) =>
                      handleEmailDataChange({
                        newValue: e.target.value,
                        property: 'subject',
                      })
                    }
                  />

                  <Editor
                    placeholder={placeholder}
                    className='mb-10 text-base'
                    dataTest='flow-email-editor'
                    namespace='flow-email-editor'
                    onChange={(e) => {
                      handleEmailDataChange({
                        newValue: e,
                        property: 'bodyTemplate',
                      });
                    }}
                  ></Editor>
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

const emailPrompts = [
  "Write something they'll want to share with their boss",
  'Craft an email that makes them say "Wow!"',
  "Compose an email they'll quote in their presentation",
  "Make them feel like they've discovered hidden treasure",
  'Write an email that makes them rethink their strategy',
  "Write something they can't get from a Google search",
  'Compose the email that ends their decision paralysis',
  "Write an email they can't ignore",
  'Turn this blank canvas into a sales masterpiece',
  'Write something that makes them feel stupid for not replying',
  'Write something that makes them say, "Yes, this is what we need!"',
  "Show them what they're missing—start typing...",
  'Type an email that helps them win',
  "Write something they'll remember",
  'Make your email impossible to ignore',
  'Start an email that stands out',
];

function getRandomEmailPrompt(): string {
  const randomIndex = Math.floor(Math.random() * emailPrompts.length);

  return emailPrompts[randomIndex];
}
