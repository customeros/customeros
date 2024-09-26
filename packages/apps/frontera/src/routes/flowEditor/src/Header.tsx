import React, { useEffect } from 'react';

import { Button } from '@ui/form/Button/Button.tsx';
// import { Users01 } from '@ui/media/icons/Users01.tsx';
// import { PieChart03 } from '@ui/media/icons/PieChart03.tsx';
// import { Dataflow03 } from '@ui/media/icons/Dataflow03.tsx';
import { ChevronRight } from '@ui/media/icons/ChevronRight.tsx';
// import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';

import { useParams, useNavigate } from 'react-router-dom';

import { useUnmount } from 'usehooks-ts';
import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';
import { FlowStore } from '@store/Flows/Flow.store.ts';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { User01 } from '@ui/media/icons/User01.tsx';

import { FlowStatusMenu } from './components';

import '@xyflow/react/dist/style.css';

export const Header = observer(() => {
  const id = useParams().id as string;
  const store = useStore();
  const navigate = useNavigate();
  const { getEdges, getNodes } = useReactFlow();

  const flow = store.flows.value.get(id) as FlowStore;

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
    const nodes = getNodes();
    const edges = getEdges();

    flow?.updateFlow({
      nodes: JSON.stringify(nodes),
      edges: JSON.stringify(edges),
    });
  });

  const handleSave = () => {
    const nodes = getNodes();
    const edges = getEdges();

    flow?.updateFlow({
      nodes: JSON.stringify(nodes),
      edges: JSON.stringify(edges),
    });
  };

  return (
    <div>
      <div className='bg-white px-10 h-14 border-b flex items-center text-base font-bold justify-between'>
        <div className='flex items-center'>
          <span
            role='button'
            onClick={() => navigate(-1)}
            className='font-medium text-gray-500'
          >
            Flows
          </span>
          <ChevronRight className='size-4 mx-1 text-gray-500' />
          <span className='mr-2'>{flow?.value?.name || 'Unnamed'}</span>

          <Button
            size='xxs'
            variant='outline'
            colorScheme='gray'
            onClick={handleSave}
            leftIcon={<User01 />}
            className='font-medium'
          >
            {flow?.value?.contacts?.length}
          </Button>
        </div>
        <div className='flex gap-2'>
          <FlowStatusMenu id={id} />
          <Button
            size='xs'
            variant='outline'
            colorScheme='gray'
            onClick={handleSave}
            leftIcon={<Check />}
          >
            Save
          </Button>
        </div>
      </div>
      {/*/!* HEADER L2 *!/*/}
      {/*<div className='bg-white px-10 border-b flex items-center text-base font-bold gap-2 py-2'>*/}
      {/*  <Button size='xs' variant='outline' leftIcon={<Dataflow03 />}>*/}
      {/*    Editor*/}
      {/*  </Button>*/}
      {/*  <Button size='xs' variant='ghost' leftIcon={<PieChart03 />}>*/}
      {/*    Report*/}
      {/*  </Button>*/}
      {/*  <Button size='xs' variant='ghost' leftIcon={<Users01 />}>*/}
      {/*    Users*/}
      {/*  </Button>*/}
      {/*</div>*/}
    </div>
  );
});
