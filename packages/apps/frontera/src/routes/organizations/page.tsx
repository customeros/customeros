import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { Preview } from '@invoices/components/Preview';

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

  useEffect(() => {
    // should be replaced by OrganizationsHub(or appropiate hub) when ready
    store.ui.commandMenu.setType('GlobalHub');
  }, [preset]);

  return (
    <div className='flex w-full items-start'>
      <div className='w-[100%] '>
        <Search open={open} onOpen={onOpen} onClose={onClose} />
        <FinderTable isSidePanelOpen={open} />
        <Preview />
      </div>
    </div>
  );
});
