import { useSearchParams } from 'react-router-dom';

export const InvoiceNumberCell = ({
  value,
  invoiceId,
}: {
  value: string;
  invoiceId: string;
}) => {
  const [searchParams, setSearchParams] = useSearchParams();

  const handleClick = () => {
    const newSearchParams = new URLSearchParams(searchParams?.toString());
    newSearchParams.set('preview', invoiceId);
    setSearchParams(`?${newSearchParams.toString()}`);
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
