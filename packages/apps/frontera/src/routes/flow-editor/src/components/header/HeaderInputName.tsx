import { useRef, useState, useEffect } from 'react';
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
  const inputRef = useRef<HTMLInputElement | null>(null);

  const showFinder = searchParams.get('show') === 'finder';

  useEffect(() => {
    if (showFinder && inputRef.current) {
      inputRef.current?.blur();
      window.getSelection()?.removeAllRanges(); // needed to remove the highlighted text
    }
  }, [showFinder]);

  return (
    <>
      <ResizableInput
        ref={inputRef}
        variant='unstyled'
        readOnly={showFinder}
        data-test='flows-flow-name-input'
        value={store.flows.isLoading ? 'Loading flow…' : name}
        onChange={(e) => {
          setName(e.target.value);
        }}
        onFocus={(e) => {
          if (!showFinder) {
            e.target.select();
          }
        }}
        className={cn({
          'text-gray-500 cursor-pointer hover:text-gray-700': showFinder,
        })}
        onClick={(e) => {
          if (showFinder) {
            e.preventDefault();

            navigate(-1);
          }
        }}
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
