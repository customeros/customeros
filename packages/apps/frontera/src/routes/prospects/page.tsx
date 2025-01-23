import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FinderTable } from '@finder/components/FinderTable';
import { FinderFilters } from '@finder/components/FinderFilters/FinderFilters';

import { cn } from '@ui/utils/cn.ts';
import { Button } from '@ui/form/Button/Button';
import { Menu01 } from '@ui/media/icons/Menu01';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Columns03 } from '@ui/media/icons/Columns03';
import { TableIdType, TableViewType } from '@graphql/types';
import { ViewSettings } from '@shared/components/ViewSettings';

import { Search } from './src/components/Search';
import { ProspectsBoard } from './src/components/ProspectsBoard';

export const ProspectsBoardPage = observer(() => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const store = useStore();
  const opportunitiesView = store.tableViewDefs.getById(
    store.tableViewDefs.opportunitiesTablePreset ?? '',
  );

  const showFinder = searchParams.get('show') === 'finder';

  return (
    <div className='flex flex-col text-gray-700 overflow-auto bg-white'>
      <div className='flex justify-between pr-4 border-b border-b-gray-200 bg-gray-25'>
        <Search />

        <div className='flex items-center'>
          <ButtonGroup className='flex items-center w-[252px]'>
            <Button
              size='xs'
              className={cn('px-4 w-full flex-1', {
                selected: showFinder,
              })}
              onClick={() =>
                navigate(`?show=finder&preset=${opportunitiesView?.id}`)
              }
            >
              <Menu01 />
              List
            </Button>
            <Button
              size='xs'
              onClick={() => navigate('/prospects')}
              className={cn('px-4 w-full flex-1', {
                selected: !showFinder,
              })}
            >
              <Columns03 />
              Board
            </Button>
          </ButtonGroup>
        </div>
      </div>
      <div
        className={cn(
          'flex justify-between items-center py-2 border-b-gray-200 pr-4',
          {
            'border-b ': !showFinder,
          },
        )}
      >
        <div>
          {showFinder && (
            <FinderFilters
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              type={TableViewType.Opportunities as any}
              tableId={TableIdType.OpportunitiesRecords}
            />
          )}
        </div>

        <div
          className={cn({
            'my-[1px]': !showFinder,
          })}
        >
          <ViewSettings
            type={TableViewType.Opportunities}
            tableId={
              showFinder
                ? TableIdType.OpportunitiesRecords
                : TableIdType.Opportunities
            }
          />
        </div>
      </div>

      {showFinder && <FinderTable />}
      {!showFinder && <ProspectsBoard />}
    </div>
  );
});
