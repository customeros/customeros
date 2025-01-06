import { useState, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { ReactFlowProvider } from '@xyflow/react';
import { FinderTable } from '@finder/components/FinderTable';
import { FinderFilters } from '@finder/components/FinderFilters/FinderFilters';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ViewSettings } from '@shared/components/ViewSettings';
import { LoadingScreen } from '@shared/components/SplashScreen/components';
import {
  TableIdType,
  TableViewType,
} from '@shared/types/__generated__/graphql.types';

import { Header } from './src/Header';
import { FlowBuilder } from './src/FlowBuilder';
import { FlowSettingsPanel } from './src/components';

import '@xyflow/react/dist/style.css';

export const FlowEditor = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const { id } = useParams();
  const [isLoading, setIsLoading] = useState(() => !store.flows.isBootstrapped);

  const [hasNewChanges, setHasNewChanges] = useState(false);
  const [isSidePanelOpen, setIsSidePanelOpen] = useState<boolean>(false);

  if (typeof id === 'undefined') {
    navigate('/finder');

    return;
  }

  useEffect(() => {
    if (!store.flows.value.has(id) && store.flows.isLoading) {
      setIsLoading(true);
      store.flows.invalidateId({ id });
    }
  }, []);

  useEffect(() => {
    if (isLoading && store.flows.value.has(id)) {
      setIsLoading(false);
    }
  }, [store.flows.value.has(id)]);

  if (isLoading) {
    return <LoadingScreen hide={false} isLoaded={false} showSplash={true} />;
  }

  return (
    <ReactFlowProvider>
      <div className='flex h-full flex-col'>
        <Header
          hasChanges={hasNewChanges}
          isSidePanelOpen={isSidePanelOpen}
          onToggleHasChanges={setHasNewChanges}
          onToggleSidePanel={setIsSidePanelOpen}
        />
        <FlowContent
          hasNewChanges={hasNewChanges}
          isSidePanelOpen={isSidePanelOpen}
          setHasNewChanges={setHasNewChanges}
          setIsSidePanelOpen={setIsSidePanelOpen}
        />
      </div>
    </ReactFlowProvider>
  );
});

const FlowContent = observer(
  ({
    isSidePanelOpen,
    setHasNewChanges,
    setIsSidePanelOpen,
    hasNewChanges,
  }: {
    hasNewChanges: boolean;
    isSidePanelOpen: boolean;
    setHasNewChanges: (data: boolean) => void;
    setIsSidePanelOpen: (data: boolean) => void;
  }) => {
    const store = useStore();

    const [searchParams] = useSearchParams();
    const id = useParams().id as string;
    const preset = searchParams.get('preset');
    const showFinder = searchParams.get('show') === 'finder';
    const tableViewDef = store.tableViewDefs.getById(preset || '');
    const tableId = tableViewDef?.value.tableId;

    const tableType = tableViewDef?.value?.tableType;

    useEffect(() => {
      setIsSidePanelOpen(true);
    }, []);

    return (
      <>
        {showFinder && (
          <div className='flex justify-start flex-col bg-white h-full '>
            <div className='mt-2 mb-2 bg-white ml-1 flex items-center justify-between'>
              <FinderFilters
                tableId={tableId || TableIdType.Organizations}
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                type={tableType || (TableViewType.Organizations as any)}
              />
              <div className='flex items-center gap-2 mr-4'>
                {tableViewDef?.hasFilters() && (
                  <Button
                    size='xs'
                    variant='ghost'
                    className='mr-4'
                    onClick={() => tableViewDef?.removeFilters()}
                  >
                    Clear
                  </Button>
                )}
                {showFinder && <ViewSettings type={TableViewType.Contacts} />}
              </div>
            </div>

            <FinderTable />
          </div>
        )}
        <div
          className={cn('flex h-full flex-col', {
            hidden: showFinder,
          })}
        >
          <FlowBuilder
            showSidePanel={isSidePanelOpen}
            onToggleSidePanel={setIsSidePanelOpen}
            onHasNewChanges={() => setHasNewChanges(true)}
          />
        </div>

        {isSidePanelOpen && (
          <FlowSettingsPanel
            id={id}
            hasChanges={hasNewChanges}
            onToggleHasChanges={setHasNewChanges}
            onToggleSidePanel={setIsSidePanelOpen}
          />
        )}
      </>
    );
  },
);
