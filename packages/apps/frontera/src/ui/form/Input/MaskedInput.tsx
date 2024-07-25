import { IMaskInput, IMaskMixinProps } from 'react-imask';

import { MaskElement } from 'imask';
import { twMerge } from 'tailwind-merge';

import { InputProps, inputVariants } from './Input';

type MaskedInputProps = IMaskMixinProps<MaskElement> & InputProps;

export const MaskedInput = ({
  size,
  variant,
  className,
  ...rest
}: MaskedInputProps) => {
  return (
    // @ts-expect-error types in this library are just confusing and not worth the effort
    <IMaskInput
      {...rest}
      data-1p-ignore
      className={twMerge(inputVariants({ className, size, variant }))}
    />
  );
};

MaskedInput.displayName = 'MaskedInput';
