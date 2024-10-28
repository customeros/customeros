import React, { useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';

import { FinderTable } from '@finder/components/FinderTable';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { useReactFlow, ReactFlowProvider } from '@xyflow/react';

import { cn } from '@ui/utils/cn';

import { Header } from './src/Header';
import { FlowBuilder } from './src/FlowBuilder';
import { FlowSettingsPanel } from './src/components';

import '@xyflow/react/dist/style.css';

export const FlowEditor = () => {
  const [hasNewChanges, setHasNewChanges] = useState(false);
  const [isSidePanelOpen, setIsSidePanelOpen] = useState<boolean>(false);

  const allowExploration = useFeatureIsOn('flow-editor-poc');

  if (!allowExploration) {
    return null;
  }

  return (
    <ReactFlowProvider>
      <div className='flex h-full flex-col'>
        <Header
          hasChanges={hasNewChanges}
          onToggleHasChanges={setHasNewChanges}
          onToggleSidePanel={() => setIsSidePanelOpen(!isSidePanelOpen)}
        />
        <FlowContent
          showSidePanel={isSidePanelOpen}
          setHasNewChanges={setHasNewChanges}
        />
      </div>
    </ReactFlowProvider>
  );
};

const FlowContent = ({
  showSidePanel,
  setHasNewChanges,
}: {
  showSidePanel: boolean;
  setHasNewChanges: (data: boolean) => void;
}) => {
  const [searchParams] = useSearchParams();
  const [isSidePanelOpen, setIsSidePanelOpen] = useState<boolean>(false);
  const { getNodes } = useReactFlow();
  const id = useParams().id as string;

  const showFinder = searchParams.get('show') === 'finder';

  return (
    <>
      {showFinder && <FinderTable isSidePanelOpen={false} />}
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
};
