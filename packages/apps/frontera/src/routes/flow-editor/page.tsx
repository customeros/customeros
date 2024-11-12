import { useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FinderTable } from '@finder/components/FinderTable';
import { useReactFlow, ReactFlowProvider } from '@xyflow/react';
import { FinderFilters } from '@finder/components/FinderFilters/FinderFilters';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ViewSettings } from '@shared/components/ViewSettings';
import {
  TableIdType,
  TableViewType,
} from '@shared/types/__generated__/graphql.types';

import { Header } from './src/Header';
import { FlowBuilder } from './src/FlowBuilder';
import { FlowSettingsPanel } from './src/components';

import '@xyflow/react/dist/style.css';

export const FlowEditor = () => {
  const [hasNewChanges, setHasNewChanges] = useState(false);
  const [isSidePanelOpen, setIsSidePanelOpen] = useState<boolean>(false);

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
          showSidePanel={isSidePanelOpen}
          setHasNewChanges={setHasNewChanges}
        />
      </div>
    </ReactFlowProvider>
  );
};

const FlowContent = observer(
  ({
    showSidePanel,
    setHasNewChanges,
  }: {
    showSidePanel: boolean;
    setHasNewChanges: (data: boolean) => void;
  }) => {
    const store = useStore();

    const [searchParams] = useSearchParams();
    const [isSidePanelOpen, setIsSidePanelOpen] = useState<boolean>(false);
    const { getNodes } = useReactFlow();
    const id = useParams().id as string;
    const preset = searchParams.get('preset');
    const showFinder = searchParams.get('show') === 'finder';
    const tableViewDef = store.tableViewDefs.getById(preset || '');
    const tableId = tableViewDef?.value.tableId;

    const tableType = tableViewDef?.value?.tableType;

    return (
      <>
        {showFinder && (
          <div className='flex justify-start flex-col bg-white '>
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

            <FinderTable isSidePanelOpen={false} />
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

        {showSidePanel && <FlowSettingsPanel id={id} nodes={getNodes()} />}
      </>
    );
  },
);
