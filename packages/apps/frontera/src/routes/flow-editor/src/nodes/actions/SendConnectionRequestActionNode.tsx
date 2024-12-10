import { MouseEventHandler } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn';
import { FlowStatus } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Trash01 } from '@ui/media/icons/Trash01';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';

export const SendConnectionRequestActionNode = observer(
  ({ id }: { id: string }) => {
    const { ui, flows } = useStore();
    const flowId = useParams()?.id as string;
    const flow = flows.value.get(flowId)?.value;
    const flowWasStarted =
      flow?.status === FlowStatus.On || flow?.firstStartedAt;

    const { deleteElements } = useReactFlow();

    const handleDelete: MouseEventHandler = (e) => {
      e.stopPropagation();

      return deleteElements({ nodes: [{ id }] });
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

            <span className='truncate'>Send connection request</span>
          </div>
          {!ui.flowActionSidePanel.isOpen && !flowWasStarted && (
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='Edit'
              icon={<Trash01 />}
              onClick={handleDelete}
              className={`ml-2 opacity-0 group-hover:opacity-100 pointer-events-all`}
            />
          )}
        </div>
      </>
    );
  },
);
