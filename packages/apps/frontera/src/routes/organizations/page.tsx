import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { FinderTable } from '@organizations/components/FinderTable';

import { Search } from './src/components/Search';

export const FinderPage = observer(() => {
  const store = useStore();
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const defaultPreset = store.tableViewDefs.defaultPreset;

  useEffect(() => {
    if (!preset && defaultPreset) {
      setSearchParams(`?preset=${defaultPreset}`);
    }
  }, [preset, setSearchParams, defaultPreset]);

  return (
    <>
      <Search />
      <FinderTable />
    </>
  );
});
