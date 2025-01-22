import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

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
        {enrichingStatus ? 'Enriching...' : 'Not set'}
      </p>
    );

  return (
    <Tooltip label={value}>
      <p className='text-gray-700 cursor-default truncate'>{value}</p>
    </Tooltip>
  );
};
