import { useState, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { useUnmount } from 'usehooks-ts';
import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';
import { FlowStore } from '@store/Flows/Flow.store';

import { cn } from '@ui/utils/cn';
import { FlowStatus } from '@graphql/types';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { User01 } from '@ui/media/icons/User01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Settings03 } from '@ui/media/icons/Settings03';
import { UserPlus01 } from '@ui/media/icons/UserPlus01';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

import { HeaderInputName } from './components/header';
import { FlowStatusMenu, FlowMoreActionsMenu } from './components';

import '@xyflow/react/dist/style.css';

export const Header = observer(
  ({
    hasChanges,
    onToggleHasChanges,
    onToggleSidePanel,
    isSidePanelOpen,
  }: {
    hasChanges: boolean;
    isSidePanelOpen: boolean;
    onToggleSidePanel: (status: boolean) => void;
    onToggleHasChanges: (status: boolean) => void;
  }) => {
    const id = useParams().id as string;
    const store = useStore();
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const { getNodes, getEdges } = useReactFlow();

    const flow = store.flows.value.get(id) as FlowStore;

    const contactsStore = store.contacts;
    const showFinder = searchParams.get('show') === 'finder';
    const viewDefPreset = flow?.value?.tableViewDefId;

    const canSave = hasChanges && flow.value.status === FlowStatus.Off;
    const [showPublishChangesButton, setShowPublishChangesButton] =
      useState(false);

    useEffect(() => {
      if (!store.ui.commandMenu.isOpen) {
        store.ui.commandMenu.setType('FlowCommands');
        store.ui.commandMenu.setContext({
          entity: 'Flow',
          ids: [id],
        });
      }
    }, [store.ui.commandMenu.isOpen, id]);

    // todo remove when we move to save on node by node basis
    const handleSaveChanges = () => {
      const nodes = getNodes();
      const edges = getEdges();

      // this should never happen
      if (nodes.length === 0 && edges.length === 0) return;
      onToggleHasChanges(false);

      flow?.updateFlow({
        nodes: JSON.stringify(nodes),
        edges: JSON.stringify(edges),
      });
    };

    useEffect(() => {
      if (canSave) {
        handleSaveChanges();
      }
    }, [showFinder]);
    useUnmount(() => {
      if (canSave) {
        handleSaveChanges();
      }
    });
    useEffect(() => {
      const handleBeforeUnload = (event: BeforeUnloadEvent) => {
        if (hasChanges) {
          event.preventDefault();
          setShowPublishChangesButton(true);

          return `Changes you’ve made will NOT be published.`;
        }
      };

      window.addEventListener('beforeunload', handleBeforeUnload);

      return () => {
        window.removeEventListener('beforeunload', handleBeforeUnload);
      };
    }, [hasChanges]);

    const handleSave = () => {
      const nodes = getNodes();
      const edges = getEdges();

      flow?.updateFlow(
        {
          nodes: JSON.stringify(nodes),
          edges: JSON.stringify(edges),
        },
        {
          onSuccess: () => {
            setTimeout(() => {
              onToggleHasChanges(false);
            }, 0);
          },
          onError: () => {
            onToggleHasChanges(true);
          },
        },
      );
    };

    return (
      <div>
        <div
          className={cn(
            'bg-white px-4 pl-2 h-[41px] flex items-center text-base font-bold justify-between',
            'border-b',
          )}
        >
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

              <HeaderInputName />
              {!showFinder && <FlowMoreActionsMenu id={id} />}
              {showFinder ? (
                <>
                  <ChevronRight className='text-gray-400' />
                  <span className='font-medium cursor-default'>
                    {`${flow?.value?.participants?.length} ${
                      flow?.value?.participants?.length === 1
                        ? 'Contact'
                        : 'Contacts'
                    }`}
                  </span>
                </>
              ) : (
                <Tooltip
                  label={
                    contactsStore.isLoading || store.flows.isLoading
                      ? `We're loading your contacts`
                      : ''
                  }
                >
                  <Button
                    size='xxs'
                    variant='outline'
                    colorScheme='gray'
                    leftIcon={<User01 />}
                    dataTest='flow-contacts'
                    isLoading={store.flows.isLoading}
                    onClick={() => {
                      navigate(`?show=finder&preset=${viewDefPreset}`);
                    }}
                    leftSpinner={
                      <Spinner
                        size='sm'
                        label='adding'
                        className='text-gray-300 fill-gray-400'
                      />
                    }
                  >
                    {flow?.value?.participants?.length}
                  </Button>
                </Tooltip>
              )}
            </div>
          </div>
          <div className='flex gap-2'>
            {showFinder && (
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
            )}
            {showPublishChangesButton && canSave && (
              <Button
                size='xs'
                variant='outline'
                colorScheme='gray'
                dataTest='save-flow'
                onClick={handleSave}
              >
                Publish changes
              </Button>
            )}
            <FlowStatusMenu
              id={id}
              hasUnsavedChanges={hasChanges}
              onToggleHasChanges={onToggleHasChanges}
              handleOpenSettingsPanel={() => onToggleSidePanel(true)}
            />
            <IconButton
              size='xs'
              variant='outline'
              icon={<Settings03 />}
              aria-label={'Toggle Settings'}
              dataTest={'flow-toggle-settings'}
              onClick={() => onToggleSidePanel(!isSidePanelOpen)}
            />
          </div>
        </div>
      </div>
    );
  },
);
