import React from 'react';

import { Play } from '@ui/media/icons/Play.tsx';
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

import { useStore } from '@shared/hooks/useStore';

import '@xyflow/react/dist/style.css';

export const Header = observer(() => {
  const id = useParams().id as string;
  const store = useStore();
  const navigate = useNavigate();
  const { getEdges, getNodes } = useReactFlow();

  const flow = store.flows.value.get(id) as FlowStore;

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
          {/*<Tag*/}
          {/*  size='md'*/}
          {/*  variant='outline'*/}
          {/*  colorScheme='gray'*/}
          {/*  onClick={() => {}}*/}
          {/*  className='bg-transparent py-0.5 cursor-pointer text-gray-500'*/}
          {/*>*/}
          {/*  <TagLeftIcon>*/}
          {/*    <Users01 />*/}
          {/*  </TagLeftIcon>*/}
          {/*  <TagLabel>4 contacts</TagLabel>*/}
          {/*</Tag>*/}
        </div>
        <div className='flex gap-2'>
          <Button
            size='xs'
            isDisabled
            variant='outline'
            leftIcon={<Play />}
            colorScheme='primary'
          >
            Start flow
          </Button>
          <Button
            size='xs'
            variant='outline'
            colorScheme='gray'
            onClick={handleSave}
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
