import { IconButton } from '@ui/form/IconButton';
import { InlineLoader } from '@ui/presentation/inline-loader';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

import CheckCircle from './assets/CheckCircle';

import styles from './simple-validation-indicator.module.scss';
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
    return <InlineLoader label='Validating' color='#DB9E00' />;
  }

  if (!errorMessages.length) {
    return (
      <div className={styles.validEntry}>
        <CheckCircle className='text-success-600 size-4' />
      </div>
    );
  }

  if (!errorMessages.length) {
    return null;
  }

  return (
    <Popover>
      <PopoverTrigger>
        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='Show validationresults'
          icon={<div className='w-2 bg-warning-400 h-2 rounded-full' />}
        />
      </PopoverTrigger>
      <PopoverContent>
        {errorMessages.map((data) => (
          <p className='text-gray-600' key={data.split(' ').join('-')}>
            {data}
          </p>
        ))}
      </PopoverContent>
    </Popover>
  );
};
