import React, { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { ReactFlowProvider } from '@xyflow/react';
import { FinderTable } from '@finder/components/FinderTable';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Header } from './src/Header';
import { FlowBuilder } from './src/FlowBuilder';

import '@xyflow/react/dist/style.css';

export const FlowEditor = () => {
  const [searchParams] = useSearchParams();
  const [hasNewChanges, setHasNewChanges] = useState(false);
  const [isSidePanelOpen, setIsSidePanelOpen] = useState<boolean>(false);

  const allowExploration = useFeatureIsOn('flow-editor-poc');

  const showFinder = searchParams.get('show') === 'finder';

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
        {showFinder ? (
          <FinderTable isSidePanelOpen={false} />
        ) : (
          <FlowBuilder
            showSidePanel={isSidePanelOpen}
            onToggleSidePanel={setIsSidePanelOpen}
            onHasNewChanges={() => setHasNewChanges(true)}
          />
        )}
      </div>
    </ReactFlowProvider>
  );
};
