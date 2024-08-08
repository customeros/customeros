import React from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { TableViewType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Star01 } from '@ui/media/icons/Star01.tsx';
import { Download02 } from '@ui/media/icons/Download02.tsx';
import { DotsVertical } from '@ui/media/icons/DotsVertical.tsx';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';
import { useDownloadCsv } from '@organizations/components/TableViewMenu/useDownloadTableViewAsCSV.ts';

interface TableViewMenuProps {}

export const TableViewMenu: React.FC<TableViewMenuProps> = observer(() => {
  const { downloadCSV } = useDownloadCsv();
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;

  return (
    <Menu>
      <MenuButton className='w-6 h-6 mr-2 outline-none focus:outline-none text-gray-400 hover:text-gray-500'>
        <DotsVertical />
      </MenuButton>
      <MenuList align='end' side='bottom'>
        <MenuItem
          className='py-2.5'
          onClick={() => store.tableViewDefs.createFavorite(preset ?? '')}
        >
          <Star01 className='text-gray-500' />
          Save to favorites
        </MenuItem>
        {tableType !== TableViewType.Invoices && (
          <MenuItem className='py-1.5' onClick={downloadCSV}>
            <Download02 className='text-gray-500' />
            Export view as CSV
          </MenuItem>
        )}
        {/*<MenuItem onClick={todo}>*/}
        {/*  <Archive />*/}
        {/*  Archive view*/}
        {/*</MenuItem>*/}
      </MenuList>
    </Menu>
  );
});
