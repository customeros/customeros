import { TableViewType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { ViewSettings } from '@shared/components/ViewSettings';

import { Search } from './src/components/Search';
import { ProspectsBoard } from './src/components/ProspectsBoard';

export const ProspectsBoardPage = () => {
  const store = useStore();
  const numberOfColumns =
    store.settings.tenant.value?.opportunityStages.length ?? 0;

  return (
    <div className='flex flex-col text-gray-700 overflow-auto bg-white'>
      <div
        style={{ minWidth: `${numberOfColumns * 150}px` }}
        className='flex justify-between pr-4 border-b border-b-gray-200 bg-gray-25 sticky top-0 z-50 '
      >
        <Search />
        <ViewSettings type={TableViewType.Opportunities} />
      </div>
      <ProspectsBoard />
    </div>
  );
};
