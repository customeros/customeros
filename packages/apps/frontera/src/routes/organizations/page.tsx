import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { FinderTable } from '@organizations/components/FinderTable';

import { Search } from './src/components/Search';

export const FinderPage = observer(() => {
  const store = useStore();
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const defaultPreset = store.tableViewDefs.defaultPreset;
  const { open, onOpen, onClose } = useDisclosure({ id: 'flow-finder' });

  useEffect(() => {
    if (!preset && defaultPreset) {
      setSearchParams(`?preset=${defaultPreset}`);
    }
  }, [preset, setSearchParams, defaultPreset]);

  return (
    <div className='flex w-full items-start'>
      <div className='w-[100%] '>
        <Search onOpen={onOpen} onClose={onClose} open={open} />
        <FinderTable isSidePanelOpen={open} />
      </div>
    </div>
  );
});
