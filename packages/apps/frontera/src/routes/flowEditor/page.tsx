import React from 'react';

import { ReactFlowProvider } from '@xyflow/react';

import { Header } from './src/Header.tsx';
import { MarketingFlowBuilder } from './src/FlowManager.tsx';

import '@xyflow/react/dist/style.css';

export const FlowEditor = () => {
  return (
    <ReactFlowProvider>
      <div className='flex h-full flex-col'>
        <Header />

        <MarketingFlowBuilder />
      </div>
    </ReactFlowProvider>
  );
};
