import { useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { useTablePlaceholder } from '@finder/hooks/useTablePlaceholder.tsx';

import { TableViewType } from '@graphql/types';
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
    const tableView = store.tableViewDefs.getById(preset || '');
    const tableViewName = tableView?.value.name;

    const { multi: multiResultPlaceholder, single: singleResultPlaceholder } =
      useTablePlaceholder(tableViewName);

    const totalResults = match(tableView?.value.tableType)
      .returnType<number>()
      .with(
        TableViewType.Organizations,
        () => store.organizations.availableCounts.get(preset ?? '') ?? 0,
      )
      .with(
        TableViewType.Contacts,
        () => store.contacts.availableCounts.get(preset ?? '') ?? 0,
      )
      .otherwise(() => store.ui.searchCount);

    const tableName =
      totalResults === 1 ? singleResultPlaceholder : multiResultPlaceholder;

    const hideSearch =
      tableView?.value?.tableType === TableViewType.Organizations ||
      tableView?.value?.tableType === TableViewType.Contacts;

    return (
      <div className='flex flex-row items-center gap-1'>
        <SearchSm className='size-5' />
        <div
          data-test={dataTest ? dataTest : ''}
          className={'font-medium flex items-center gap-1 break-keep w-max '}
        >
          {`${totalResults} ${tableName}` + (hideSearch ? '' : ':')}
        </div>
      </div>
    );
  },
);
