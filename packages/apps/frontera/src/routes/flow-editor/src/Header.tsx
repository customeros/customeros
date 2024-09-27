import { useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { Button } from '@ui/form/Button/Button.tsx';
// import { Users01 } from '@ui/media/icons/Users01.tsx';
// import { PieChart03 } from '@ui/media/icons/PieChart03.tsx';
// import { Dataflow03 } from '@ui/media/icons/Dataflow03.tsx';
import { ChevronRight } from '@ui/media/icons/ChevronRight.tsx';
// import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';

import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';
import { FlowStore } from '@store/Flows/Flow.store.ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { User01 } from '@ui/media/icons/User01.tsx';
import { ViewSettings } from '@shared/components/ViewSettings';
import { TableViewType } from '@shared/types/__generated__/graphql.types';

import { FlowStatusMenu } from './components';

import '@xyflow/react/dist/style.css';

export const Header = observer(() => {
  const id = useParams().id as string;
  const store = useStore();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { getEdges, getNodes } = useReactFlow();
  const saveFlag = useFeatureIsOn('flow-editor-save-button_1');

  const flow = store.flows.value.get(id) as FlowStore;
  const showFinder = searchParams.get('show') === 'finder';
  const contactsPreset = store.tableViewDefs.contactsPreset;

  useEffect(() => {
    if (!store.ui.commandMenu.isOpen) {
      store.ui.commandMenu.setType('FlowCommands');
      store.ui.commandMenu.setContext({
        entity: 'Flow',
        ids: [id],
      });
    }
  }, [store.ui.commandMenu.isOpen, id]);

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
            leftIcon={<User01 />}
            className='font-medium'
            onClick={() => navigate(`?show=finder&preset=${contactsPreset}`)}
          >
            {flow?.value?.contacts?.length}
          </Button>
        </div>

        {showFinder ? (
          <ViewSettings type={TableViewType.Contacts} />
        ) : (
          <div className='flex gap-2'>
            <FlowStatusMenu id={id} />
            {saveFlag && (
              <Button
                size='xs'
                variant='outline'
                colorScheme='gray'
                onClick={handleSave}
                leftIcon={<Check />}
              >
                Save
              </Button>
            )}
          </div>
        )}
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
