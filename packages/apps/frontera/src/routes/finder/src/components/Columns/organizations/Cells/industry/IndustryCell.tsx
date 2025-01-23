export const IndustryCell = ({
  value,
  enrichingStatus,
}: {
  value?: string;
  enrichingStatus: boolean;
}) => {
  if (!value)
    return (
      <p className='text-gray-400'>
        {enrichingStatus ? 'Enriching...' : 'Not found yet'}
      </p>
    );

  return (
    <p title={value} className='text-gray-700 cursor-default truncate'>
      {value}
    </p>
  );
};
