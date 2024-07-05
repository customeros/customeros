import { TableViewType } from '@graphql/types';

import { EditColumns } from './EditColumns';

interface ViewSettingsProps {
  type: TableViewType;
}

export const ViewSettings = ({ type }: ViewSettingsProps) => {
  return (
    <div className='flex items-center'>
      <EditColumns type={type} />
    </div>
  );
};
