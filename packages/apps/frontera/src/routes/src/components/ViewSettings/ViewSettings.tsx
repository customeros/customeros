import { TableViewType } from '@graphql/types';

import { EditColumns } from './EditColumns';

interface ViewSettingsProps {
  type: TableViewType;
}

export const ViewSettings = ({ type }: ViewSettingsProps) => {
  return (
    <div className='flex pr-2 gap-2 items-center'>
      <EditColumns type={type} />
    </div>
  );
};
