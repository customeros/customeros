import { EditColumns } from './EditColumns';

interface ViewSettingsProps {
  type: 'invoices' | 'renewals' | 'organizations';
}

export const ViewSettings = ({ type }: ViewSettingsProps) => {
  return (
    <div className='flex pr-2 gap-2 items-center'>
      <EditColumns type={type} />
    </div>
  );
};
