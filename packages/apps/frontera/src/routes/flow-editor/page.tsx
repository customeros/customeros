import React from 'react';
import { useSearchParams } from 'react-router-dom';

import { ReactFlowProvider } from '@xyflow/react';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Header } from './src/Header';
import { FlowBuilder } from './src/FlowBuilder';
import { SubjectsTable } from './src/components';

import '@xyflow/react/dist/style.css';

export const FlowEditor = () => {
  const [searchParams] = useSearchParams();
  const allowExploration = useFeatureIsOn('flow-editor-poc');

  const showSubjects = searchParams.get('show') === 'subjects';

  if (!allowExploration) {
    return null;
  }

  return (
    <ReactFlowProvider>
      <div className='flex h-full flex-col'>
        <Header />
        {showSubjects ? <SubjectsTable /> : <FlowBuilder />}
      </div>
    </ReactFlowProvider>
  );
};
