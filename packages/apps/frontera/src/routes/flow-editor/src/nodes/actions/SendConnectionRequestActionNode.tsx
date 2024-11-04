import { FlowActionType } from '@store/Flows/types.ts';
import { NodeProps, useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn';
import { IconButton } from '@ui/form/IconButton';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';

import { LinkedInMessageEditorModal } from '../../components/email';

export const SendConnectionRequestActionNode = ({
  id,
  data,
}: NodeProps & {
  data: {
    isEditing?: boolean;
    action: FlowActionType;
    messageTemplate: string;
  };
}) => {
  const { setNodes } = useReactFlow();

  const handleEmailDataChange = ({
    messageTemplate,
  }: {
    messageTemplate: string;
  }) => {
    setNodes((nds) =>
      nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              messageTemplate,
              isEditing: false,
            },
          };
        }

        return node;
      }),
    );
  };

  const handleCancel = () => {
    setNodes((nds) =>
      nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              isEditing: false,
            },
          };
        }

        return node;
      }),
    );
  };

  const toggleEditing = () => {
    setNodes((nds) =>
      nds.map((node) =>
        node.id === id
          ? { ...node, data: { ...node.data, isEditing: true } }
          : node,
      ),
    );
  };

  return (
    <>
      <div className='text-sm flex items-center justify-between overflow-hidden w-full'>
        <div className='truncate text-sm flex items-center'>
          <div
            className={cn(
              `size-6 min-w-6 mr-2 bg-blue-50 text-blue-500 border border-blue-100 rounded flex items-center justify-center`,
            )}
          >
            <LinkedinOutline className='size-4 text-blue-500' />
          </div>

          <span className='truncate'>
            {data?.messageTemplate
              ? data.messageTemplate
              : 'Send connection request'}
          </span>
        </div>
        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='Edit'
          icon={<Edit03 />}
          onClick={toggleEditing}
          className='ml-2 opacity-0 group-hover:opacity-100 pointer-events-all'
        />
      </div>
      <LinkedInMessageEditorModal
        data={data}
        handleCancel={handleCancel}
        isEditorOpen={data?.isEditing || false}
        handleEmailDataChange={handleEmailDataChange}
      />
    </>
  );
};
