import { XCircle } from '@ui/media/icons/XCircle.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { InlineLoader } from '@ui/presentation/inline-loader';
import { CheckVerified01 } from '@ui/media/icons/CheckVerified01.tsx';

interface Props {
  isLoading: boolean;
  errorMessages?: Array<string>;
  showValidationMessage: boolean;
}

export const SimpleValidationIndicator = ({
  isLoading,
  errorMessages = [],
}: Props) => {
  if (isLoading) {
    return <InlineLoader color='#DB9E00' label='Validating' />;
  }

  if (!errorMessages.length) {
    return (
      <div className='flex items-center ml-2'>
        <CheckVerified01 className='text-success-600 size-3' />
      </div>
    );
  }

  return (
    <Tooltip label={errorMessages?.join(', ') ?? ''}>
      <div className='flex items-center ml-2'>
        <XCircle className='text-warning-500 size-3' />
      </div>
    </Tooltip>
  );
};
