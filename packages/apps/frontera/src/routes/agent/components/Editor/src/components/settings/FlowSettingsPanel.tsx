import { useSearchParams } from 'react-router-dom';

import { useNodes } from '@xyflow/react';
import { FlowActionType } from '@store/Flows/types';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Settings03 } from '@ui/media/icons/Settings03';
import { UserPlus01 } from '@ui/media/icons/UserPlus01.tsx';

import { SenderSettings } from './SenderSettings';
import { NoEmailNodesPanel } from './NoEmailNodesPanel';

export const FlowSettingsPanel = ({
  id,
  onToggleSidePanel,
}: {
  id: string;
  onToggleSidePanel: (status: boolean) => void;
}) => {
  const store = useStore();
  const nodes = useNodes();
  const [searchParams] = useSearchParams();
  const showFinder = searchParams.get('show') === 'finder';
  const hasEmailNodes = nodes.some(
    (node) =>
      node.data?.action &&
      [FlowActionType.EMAIL_NEW, FlowActionType.EMAIL_REPLY].includes(
        node.data.action as FlowActionType,
      ),
  );
  const hasLinkedInNodes = nodes.some(
    (node) =>
      node.data?.action &&
      [
        FlowActionType.LINKEDIN_CONNECTION_REQUEST,
        FlowActionType.LINKEDIN_MESSAGE,
      ].includes(node.data.action as FlowActionType),
  );

  const showSenderSettings = hasEmailNodes || hasLinkedInNodes;

  return (
    <article
      onClick={(e) => e.stopPropagation()}
      className='fixed z-10 top-[0px] bottom-0 right-0 w-[400px] bg-white  border-l flex flex-col gap-4 animate-slideLeft'
    >
      {showFinder && (
        <div className='absolute left-[-138px] flex items-center top-[5px]'>
          <Button
            size='xs'
            variant='outline'
            colorScheme='gray'
            leftIcon={<UserPlus01 />}
            onClick={() =>
              store.ui.commandMenu.setOpen(true, {
                type: 'AddContactsToFlow',
                context: {
                  entity: 'Flow',
                  ids: [id],
                },
              })
            }
          >
            Add contacts
          </Button>
        </div>
      )}

      <div className='flex justify-between items-center border-b border-gray-200 p-4 y-2 h-[41px]'>
        <h1 className='font-medium'>Flow settings</h1>

        <div className='flex gap-2'>
          <IconButton
            size='xs'
            variant='outline'
            icon={<Settings03 />}
            aria-label={'Toggle Settings'}
            dataTest={'flow-toggle-settings'}
            onClick={() => onToggleSidePanel(false)}
          />
        </div>
      </div>
      <div className='px-4 gap-2 flex flex-col'>
        {showSenderSettings && (
          <SenderSettings
            id={id}
            hasEmailNodes={hasEmailNodes}
            hasLinkedInNodes={hasLinkedInNodes}
          />
        )}
        {!showSenderSettings && <NoEmailNodesPanel />}
      </div>
    </article>
  );
};
