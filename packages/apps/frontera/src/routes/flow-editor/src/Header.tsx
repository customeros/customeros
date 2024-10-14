import React, { useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { useUnmount } from 'usehooks-ts';
import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';
import { FlowStore } from '@store/Flows/Flow.store.ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { User01 } from '@ui/media/icons/User01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Settings03 } from '@ui/media/icons/Settings03.tsx';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { ViewSettings } from '@shared/components/ViewSettings';
import { TableViewType } from '@shared/types/__generated__/graphql.types';

import { FlowStatusMenu } from './components';

import '@xyflow/react/dist/style.css';

export const Header = observer(
  ({
    hasChanges,
    onToggleHasChanges,
    onToggleSidePanel,
  }: {
    hasChanges: boolean;
    onToggleSidePanel: () => void;
    onToggleHasChanges: (status: boolean) => void;
  }) => {
    const id = useParams().id as string;
    const store = useStore();
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const { getNodes, getEdges } = useReactFlow();

    const saveFlag = useFeatureIsOn('flow-editor-save-button_1');

    const flow = store.flows.value.get(id) as FlowStore;
    const contactsStore = store.contacts;
    const showFinder = searchParams.get('show') === 'finder';
    const flowContactsPreset = store.tableViewDefs.flowContactsPreset;

    useEffect(() => {
      if (!store.ui.commandMenu.isOpen) {
        store.ui.commandMenu.setType('FlowCommands');
        store.ui.commandMenu.setContext({
          entity: 'Flow',
          ids: [id],
        });
      }
    }, [store.ui.commandMenu.isOpen, id]);

    useUnmount(() => {
      if (saveFlag && !showFinder) {
        const nodes = getNodes();
        const edges = getEdges();

        // this should never happen
        if (nodes.length === 0 && edges.length === 0) return;

        flow?.updateFlow({
          nodes: JSON.stringify(nodes),
          edges: JSON.stringify(edges),
        });
      }
    });

    const handleSave = () => {
      const nodes = getNodes();
      const edges = getEdges();

      onToggleHasChanges(false);

      flow?.updateFlow(
        {
          nodes: JSON.stringify(nodes),
          edges: JSON.stringify(edges),
        },
        {
          onError: () => {
            onToggleHasChanges(true);
          },
        },
      );
    };

    return (
      <div>
        <div className='bg-white px-4 pl-2 h-[41px] border-b flex items-center text-base font-bold justify-between'>
          <div className='flex items-center'>
            <div className='flex items-center gap-1 font-medium'>
              <span
                role='button'
                data-test='navigate-to-flows'
                onClick={() => navigate(showFinder ? -2 : -1)}
                className='font-medium text-gray-500 hover:text-gray-700'
              >
                Flows
              </span>
              <ChevronRight className='text-gray-400' />
              <span
                data-test='flows-flow-name'
                onClick={() => (showFinder ? navigate(-1) : null)}
                className={cn({
                  'text-gray-500 cursor-pointer hover:text-gray-700':
                    showFinder,
                })}
              >
                {store.flows.isLoading
                  ? 'Loading flow…'
                  : flow?.value?.name || 'Unnamed'}
              </span>
              {showFinder ? (
                <>
                  <ChevronRight className='text-gray-400' />
                  <span className='font-medium cursor-default'>
                    {`${flow?.value?.contacts?.length} ${
                      flow?.value?.contacts?.length === 1
                        ? 'Contact'
                        : 'Contacts'
                    }`}
                  </span>
                </>
              ) : (
                <Button
                  size='xxs'
                  className='ml-2'
                  variant='outline'
                  colorScheme='gray'
                  leftIcon={<User01 />}
                  dataTest='flow-contacts'
                  isLoading={contactsStore.isLoading || store.flows.isLoading}
                  onClick={() => {
                    navigate(`?show=finder&preset=${flowContactsPreset}`);

                    if (saveFlag) handleSave();
                  }}
                  leftSpinner={
                    <Spinner
                      size='sm'
                      label='adding'
                      className='text-gray-300 fill-gray-400'
                    />
                  }
                >
                  {flow?.value?.contacts?.length}
                </Button>
              )}
            </div>
          </div>
          <div className='flex gap-2'>
            {showFinder && <ViewSettings type={TableViewType.Contacts} />}
            {!showFinder && saveFlag && hasChanges && (
              <Button
                size='xs'
                variant='outline'
                colorScheme='gray'
                dataTest='save-flow'
                onClick={handleSave}
              >
                Save
              </Button>
            )}
            <FlowStatusMenu id={id} />
            <IconButton
              size='xs'
              variant='ghost'
              icon={<Settings03 />}
              onClick={onToggleSidePanel}
              aria-label={'Toggle Settings'}
            />
          </div>
        </div>
      </div>
    );
  },
);
