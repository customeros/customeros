import { ReactNode } from 'react';
import { useParams } from 'react-router-dom';

import { FlagWrongFieldUsecase } from '@domain/usecases/organization-about-panel/flag-wrong-field.usecase';

import { FlagWrongFields } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { ThumbsDown } from '@ui/media/icons/ThumbsDown';

interface FieldMarkerProps {
  icon: ReactNode;
  dataTest?: string;
  placeholder?: string;
  value?: string | null;
  field: FlagWrongFields;
}

const flagWrongFieldUsecase = new FlagWrongFieldUsecase();

export const AboutTabField = ({
  field,
  icon,
  value,
  dataTest,
  placeholder,
}: FieldMarkerProps) => {
  const id = useParams()?.id as string;
  const label = fieldLabels[field];

  return (
    <div className='flex group'>
      <Tooltip align='start' label={label}>
        <p className='text-sm flex items-center cursor-default '>
          {icon}
          {value ? (
            <span>{value}</span>
          ) : (
            <span data-test={dataTest} className={'text-gray-400'}>
              {placeholder}
            </span>
          )}
        </p>
      </Tooltip>
      {value && (
        <Tooltip label={`This ${label.toLowerCase()} is incorrect`}>
          <div>
            <IconButton
              size='xxs'
              variant='ghost'
              icon={<ThumbsDown />}
              aria-label={`Mark this ${label} as incorrect`}
              className='opacity-0 group-hover:opacity-100 ml-2'
              onClick={() => flagWrongFieldUsecase.flagWrongField(id, field)}
            />
          </div>
        </Tooltip>
      )}
    </div>
  );
};

const fieldLabels: Record<FlagWrongFields, string> = {
  [FlagWrongFields.OrganizationIndustry]: 'Industry',
};
