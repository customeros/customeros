import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FinderTable } from '@finder/components/FinderTable';
import { FinderFilters } from '@finder/components/FinderFilters/FinderFilters';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ViewSettings } from '@shared/components/ViewSettings';
import {
  TableIdType,
  TableViewType,
} from '@shared/types/__generated__/graphql.types';

export const List = observer(() => {
  const store = useStore();

  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset || '');
  const tableId = tableViewDef?.value.tableId;
  const tableType = tableViewDef?.value?.tableType;

  return (
    <>
      <div className='py-2 ml-1 flex items-center justify-between border-b '>
        <FinderFilters
          tableId={tableId || TableIdType.Organizations}
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          type={tableType || (TableViewType.Organizations as any)}
        />
        <div className='flex items-center gap-2 mr-4'>
          {tableViewDef?.hasFilters() && (
            <Button
              size='xs'
              variant='ghost'
              className='mr-4'
              onClick={() => tableViewDef?.removeFilters()}
            >
              Clear
            </Button>
          )}
          <ViewSettings type={TableViewType.Contacts} />
        </div>
      </div>
      <div className='flex justify-start flex-col bg-white h-full '>
        <FinderTable />
      </div>
    </>
  );
});
