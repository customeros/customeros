// import { IconButton } from '@ui/form/IconButton/IconButton';
// import { DotsVertical } from '@ui/media/icons/DotsVertical';

import { EditColumns } from './EditColumns';

interface ViewSettingsProps {
  type: 'invoices' | 'renewals';
}

export const ViewSettings = ({ type }: ViewSettingsProps) => {
  return (
    <div className='flex pr-2 gap-2 items-center'>
      <EditColumns type={type} />
      {/* <IconButton
        size='xs'
        variant='outline'
        icon={<DotsVertical />}
        aria-label='View options'
      /> */}
    </div>
  );
};
