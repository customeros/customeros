import { Link } from 'react-router-dom';
import React, { memo, useRef, useState } from 'react';

import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';

import { Input } from '../Input/Input';
import { formatSocialUrl } from './util';
import { IconButton } from '../IconButton/IconButton';

export interface UrlInputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  value: string;
  label?: string;
  isLabelVisible?: boolean;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

export const UrlInput = memo(
  ({
    value,
    onBlur,
    isLabelVisible,
    label,
    labelProps,
    ...rest
  }: UrlInputProps) => {
    const [isFocused, setIsFocused] = useState(false);
    const [isHovered, setIsHovered] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);

    const href = value?.startsWith('http') ? value : `https://${value}`;

    const formattedUrl = formatSocialUrl(value);

    const handleFocus = () => {
      setIsFocused(true);
      inputRef.current?.focus();
    };

    const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
      onBlur?.(e);
      setIsFocused(false);
    };

    return (
      <div className='w-full'>
        {isLabelVisible ? (
          <label {...labelProps}>{label}</label>
        ) : (
          <label className='sr-only'>{label}</label>
        )}

        <div className='relative'>
          <Input
            value={value}
            ref={inputRef}
            onFocus={handleFocus}
            onBlur={handleBlur}
            size='xs'
            variant='unstyled'
            className='border border-transparent text-md'
            {...rest}
          />
          {!isFocused && !!value && (
            <div
              className='bg-gray-25 w-full absolute top-[1px] hover:border-gray-300 hover:border-b border-b border-transparent'
              onMouseEnter={() => setIsHovered(true)}
              onMouseLeave={() => setIsHovered(false)}
            >
              <p
                className='text-gray-700 top-0 truncate text-base'
                onClick={handleFocus}
              >
                {formattedUrl}
              </p>
              {isHovered && (
                <Link
                  to={href}
                  target='_blank'
                  rel='noopener noreferrer'
                  className='absolute -top-[1px] right-0 flex items-center text-gray-500 hover:text-gray-900'
                >
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    colorScheme='gray'
                    aria-label='social link'
                    className='mt-[5px]'
                    icon={<LinkExternal02 />}
                  />
                </Link>
              )}
            </div>
          )}
        </div>
      </div>
    );
  },
);
