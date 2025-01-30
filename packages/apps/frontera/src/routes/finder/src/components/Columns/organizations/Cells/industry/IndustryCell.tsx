import { FlagWrongFieldUsecase } from '@domain/usecases/organization-industry-field/flag-wrong-field.usecase';

import { FlagWrongFields } from '@graphql/types';
import { ThumbsDown } from '@ui/media/icons/ThumbsDown.tsx';

const flagWrongFieldUsecase = new FlagWrongFieldUsecase();

export const IndustryCell = ({
  value,
  id,
  enrichingStatus,
  flaggedAsIncorrect,
}: {
  id: string;
  value?: string;
  enrichingStatus: boolean;
  flaggedAsIncorrect: boolean;
}) => {
  if (!value)
    return (
      <p className='text-gray-400'>
        {enrichingStatus ? 'Enriching...' : 'Not found yet'}
      </p>
    );

  return (
    <div className='flex items-center gap-2 group/industry'>
      <p title={value} className='text-gray-700 cursor-default truncate group'>
        {value}
      </p>

      <div
        tabIndex={0}
        role={'button'}
        title={
          flaggedAsIncorrect
            ? 'Reported as incorrect before'
            : `This industry is incorrect`
        }
        onClick={() =>
          flagWrongFieldUsecase.flagWrongField(
            id,
            FlagWrongFields.OrganizationIndustry,
          )
        }
      >
        <ThumbsDown className='text-gray-500 hover:text-gray-700 opacity-0 group-hover/industry:opacity-100 block size-3' />
      </div>
    </div>
  );
};
