import { TableIdType, TableViewType } from '@graphql/types';

import { EditColumns } from './EditColumns';

interface ViewSettingsProps {
  type: TableViewType;
  tableId?: TableIdType;
}

export const ViewSettings = ({ type, tableId }: ViewSettingsProps) => {
  return (
    <div className='flex'>
      <EditColumns type={type} tableId={tableId} />
    </div>
  );
};
