import { useState } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store';

import { cn } from '@ui/utils/cn';
import { ResizableInput } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';

export const HeaderInputName = observer(() => {
  const id = useParams().id as string;
  const store = useStore();
  const flow = store.flows.value.get(id) as FlowStore;

  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [name, setName] = useState(flow?.value?.name ?? '');

  const showFinder = searchParams.get('show') === 'finder';

  return (
    <>
      <ResizableInput
        autoFocus
        variant='unstyled'
        readOnly={showFinder}
        data-test='flows-flow-name-input'
        onClick={(e) => (showFinder ? navigate(-1) : null)}
        value={store.flows.isLoading ? 'Loading flow…' : name}
        onChange={(e) => {
          setName(e.target.value);
        }}
        className={cn({
          'text-gray-500 cursor-pointer hover:text-gray-700': showFinder,
        })}
        onBlur={() => {
          flow?.update((value) => {
            value.name = name;

            return value;
          });
        }}
      />
    </>
  );
});
