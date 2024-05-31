import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { useStore } from '@shared/hooks/useStore';

import { Search } from './src/components/Search';
import { OrganizationsTable } from './src/components/OrganizationsTable';

export const OrganizationsPage = () => {
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
      <OrganizationsTable />
    </>
  );
};
