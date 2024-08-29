import { TableIdType, TableViewType } from '@graphql/types';

import { EditColumns } from './EditColumns';

interface ViewSettingsProps {
  type: TableViewType;
  tableId?: TableIdType;
}

export const ViewSettings = ({ type, tableId }: ViewSettingsProps) => {
  return (
    <div className='flex items-center'>
      <EditColumns type={type} tableId={tableId} />
    </div>
  );
};
