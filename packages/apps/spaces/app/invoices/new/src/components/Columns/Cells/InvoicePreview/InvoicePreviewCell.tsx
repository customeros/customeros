import { useRouter, useSearchParams } from 'next/navigation';

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
    <span
      className='font-medium cursor-pointer hover:text-gray-900 transition-colors'
      onClick={handleClick}
    >
      {value}
    </span>
  );
};
