import React from 'react';
import styles from './simple-validation-indicator.module.scss';
import { InlineLoader } from '@ui/presentation/inline-loader';
import {
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
  Portal,
} from '@ui/overlay/Popover';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';

import CheckCircle from './assets/CheckCircle';
interface Props {
  showValidationMessage: boolean;
  isLoading: boolean;
  errorMessages?: Array<string>;
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
        <CheckCircle
          color='var(--chakra-colors-success-600)'
          height={16}
          width={16}
        />
      </div>
    );
  }

  if (!errorMessages.length) {
    return null;
  }

  return (
    <Popover trigger='hover'>
      <PopoverTrigger>
        <IconButton
          size='sm'
          variant='flushed'
          aria-label='Show validationresults'
          icon={<div className={styles.validationSignal} />}
        />
      </PopoverTrigger>
      <Portal>
        <PopoverContent>
          <PopoverArrow />
          <PopoverBody>
            {errorMessages.map((data) => (
              <Text color='gray.600' key={data.split(' ').join('-')}>
                {data}
              </Text>
            ))}
          </PopoverBody>
        </PopoverContent>
      </Portal>
    </Popover>
  );
};
