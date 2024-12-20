import { useSearchParams } from 'react-router-dom';

import { Eye } from '@ui/media/icons/Eye';

export const InvoicePreviewCell = ({ value }: { value: string }) => {
  const [searchParams, setSearchParams] = useSearchParams();

  const handleClick = () => {
    const newSearchParams = new URLSearchParams(searchParams?.toString());

    newSearchParams.set('preview', value);
    setSearchParams(newSearchParams.toString());
  };

  return (
    <div
      tabIndex={0}
      role='button'
      onClick={handleClick}
      className='font-medium cursor-pointer hover:text-gray-900 transition-colors group flex gap-1 items-center mr-2'
    >
      <span className='truncate'>{value}</span>
      <div>
        <Eye className='opacity-0 group-hover:opacity-100 transition-opacity text-gray-400 size-4' />
      </div>
    </div>
  );
};
