import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { Preview } from '@invoices/components/Preview';
import { FinderTable } from '@finder/components/FinderTable';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { FinderFilters } from '@finder/components/FinderFilters/FinderFilters';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ViewSettings } from '@shared/components/ViewSettings';
import {
  TableIdType,
  TableViewType,
} from '@shared/types/__generated__/graphql.types';

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

  const tableViewDef = store.tableViewDefs.getById(preset || '');
  const tableId = tableViewDef?.value.tableId;
  const tableViewType = tableViewDef?.value.tableType;

  const tableType = tableViewDef?.value?.tableType;
  const flag = useFeatureIsOn('filters-v2');

  return (
    <div className='flex w-full items-start'>
      <div className='w-[100%] bg-white'>
        <Search open={open} onOpen={onOpen} onClose={onClose} />
        {flag && (
          <div className='flex justify-between mx-4 my-2 items-start'>
            <FinderFilters
              tableId={tableId || TableIdType.Organizations}
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              type={tableType || (TableViewType.Organizations as any)}
            />
            <div className='flex items-center gap-2'>
              <Button
                size='xs'
                variant='ghost'
                onClick={() => tableViewDef?.removeFilters()}
              >
                Clear
              </Button>
              {tableViewType && (
                <ViewSettings tableId={tableId} type={tableViewType} />
              )}
            </div>
          </div>
        )}
        <FinderTable isSidePanelOpen={open} />
        <Preview />
      </div>
    </div>
  );
});
