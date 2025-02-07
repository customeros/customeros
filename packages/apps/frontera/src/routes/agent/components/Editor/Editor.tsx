import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { ReactFlowProvider } from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';
import { LoadingScreen } from '@shared/components/SplashScreen/components';

import { FlowBuilder } from './src/FlowBuilder.tsx';

import '@xyflow/react/dist/style.css';

export const Editor = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const { id } = useParams();
  const [isLoading, setIsLoading] = useState(() => !store.flows.isBootstrapped);

  const [_hasNewChanges, setHasNewChanges] = useState(false);
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
      <div className='h-full ' style={{ height: `calc(100vh - 95px)` }}>
        <FlowBuilder
          showSidePanel={isSidePanelOpen}
          onToggleSidePanel={setIsSidePanelOpen}
          onHasNewChanges={() => setHasNewChanges(true)}
        />
      </div>
    </ReactFlowProvider>
  );
});
