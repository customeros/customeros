import React from 'react';

import { ReactFlowProvider } from '@xyflow/react';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Header } from './src/Header';
import { FlowBuilder } from './src/FlowBuilder';

import '@xyflow/react/dist/style.css';

export const FlowEditor = () => {
  const allowExploration = useFeatureIsOn('flow-editor-poc');

  if (!allowExploration) {
    return null;
  }

  return (
    <ReactFlowProvider>
      <div className='flex h-full flex-col'>
        <Header />
        <FlowBuilder />
      </div>
    </ReactFlowProvider>
  );
};
