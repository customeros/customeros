import { ViewSettings } from '@shared/components/ViewSettings';

import { Search } from './src/components/Search';
import { OrganizationsTable } from './src/components/OrganizationsTable';

export const OrganizationsPage = () => {
  return (
    <>
      <div className='flex items-center w-full justify-between'>
        <Search />
        <ViewSettings type='organizations' />
      </div>
      <OrganizationsTable />
    </>
  );
};
