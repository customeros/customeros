import { useRouter, useSearchParams } from 'next/navigation';

import { Eye } from '@ui/media/icons/Eye';

export const InvoicePreviewCell = ({
  value,
  invoiceId,
}: {
  value: string;
  invoiceId: string;
}) => {
  const router = useRouter();
  const searchParams = useSearchParams();

  const handleClick = () => {
    const newSearchParams = new URLSearchParams(searchParams?.toString());
    newSearchParams.set('preview', invoiceId);
    router.push(`?${newSearchParams.toString()}`);
  };

  return (
    <div className='flex gap-1 items-center'>
      <span
        onClick={handleClick}
        className='font-medium cursor-pointer hover:text-gray-900 transition-colors peer'
      >
        {value}
      </span>
      <Eye className='opacity-0 peer-hover:opacity-100 transition-opacity text-gray-400 size-4' />
    </div>
  );
};
