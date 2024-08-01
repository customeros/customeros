import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
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
  const currentPreset = store.tableViewDefs
    ?.toArray()
    .find((e) => e.value.id === preset)?.value?.name;

  useEffect(() => {
    if (!preset && defaultPreset) {
      setSearchParams(`?preset=${defaultPreset}`);
    }
  }, [preset, setSearchParams, defaultPreset]);

  useEffect(() => {
    match(currentPreset)
      .with('All orgs', () => {
        store.ui.commandMenu.setType('OrganizationHub');
      })

      .otherwise(() => {
        store.ui.commandMenu.setType('GlobalHub');
      });
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
