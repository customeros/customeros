import { MouseEventHandler } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { useDownloadCsv } from '@finder/components/TableViewMenu/useDownloadTableViewAsCSV.ts';

import { TableViewType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive.tsx';
import { Download02 } from '@ui/media/icons/Download02.tsx';
import { LayersTwo01 } from '@ui/media/icons/LayersTwo01.tsx';
import { DotsVertical } from '@ui/media/icons/DotsVertical.tsx';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

export const TableViewMenu = observer(() => {
  const { downloadCSV } = useDownloadCsv();
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset') ?? '1';
  const flag = useFeatureIsOn('filters-v2');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset);
  const tableType = tableViewDef?.value?.tableType;

  const isPreset = tableViewDef?.value?.isPreset;

  const handleAddToMyViews: MouseEventHandler<HTMLDivElement> = (e) => {
    e.stopPropagation();

    if (!preset) {
      store.ui.toastError(
        `We were unable to add this view to favorites`,
        'dup-view-error',
      );

      return;
    }
    store.ui.commandMenu.toggle('DuplicateView');
    store.ui.commandMenu.setContext({
      ids: [preset],
      entity: 'TableViewDef',
    });
  };

  if (
    !isPreset ||
    (tableType &&
      [TableViewType.Invoices, TableViewType.Flow].includes(tableType))
  )
    return <div className='ml-2'></div>;

  return flag ? (
    <Download02 onClick={downloadCSV} className='text=gray-500 mr-3.5 ml-2' />
  ) : (
    <Menu>
      <MenuButton className='w-6 h-6 mr-2 outline-none focus:outline-none text-gray-400 hover:text-gray-500'>
        <DotsVertical />
      </MenuButton>
      <MenuList align='end' side='bottom'>
        <MenuItem onClick={handleAddToMyViews}>
          <LayersTwo01 className='text-gray-500' />
          Save as...
        </MenuItem>
        {tableType &&
          ![TableViewType.Invoices, TableViewType.Flow].includes(tableType) && (
            <MenuItem className='py-1.5' onClick={downloadCSV}>
              <Download02 className='text-gray-500' />
              Export view as CSV
            </MenuItem>
          )}

        {!isPreset && (
          <MenuItem
            onClick={() => {
              store.ui.commandMenu.setContext({
                ids: [preset],
                entity: 'TableViewDef',
              });
              store.ui.commandMenu.setType('DeleteConfirmationModal');
              store.ui.commandMenu.setOpen(true);
            }}
          >
            <Archive className='text-gray-500' />
            Archive view
          </MenuItem>
        )}
      </MenuList>
    </Menu>
  );
});
