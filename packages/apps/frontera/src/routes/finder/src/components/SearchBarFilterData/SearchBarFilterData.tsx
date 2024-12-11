import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useTablePlaceholder } from '@finder/hooks/useTablePlaceholder.tsx';

import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';

interface SearchBarFilterDataProps {
  dataTest?: string;
}

export const SearchBarFilterData = observer(
  ({ dataTest }: SearchBarFilterDataProps) => {
    const store = useStore();
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');
    const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;

    const { multi: multiResultPlaceholder, single: singleResultPlaceholder } =
      useTablePlaceholder(tableViewName);

    const totalResults = store.ui.searchCount;

    const tableName =
      totalResults === 1 ? singleResultPlaceholder : multiResultPlaceholder;

    return (
      <div className='flex flex-row items-center gap-1'>
        <SearchSm className='size-5' />
        <div
          data-test={dataTest ? dataTest : ''}
          className={'font-medium flex items-center gap-1 break-keep w-max '}
        >
          {totalResults} {tableName}:
        </div>
      </div>
    );
  },
);
