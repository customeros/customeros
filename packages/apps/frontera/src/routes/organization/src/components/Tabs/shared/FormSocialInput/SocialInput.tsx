import { Link } from 'react-router-dom';
import { memo, useRef, useState } from 'react';

import { cn } from '@ui/utils/cn';
import { Input, InputProps } from '@ui/form/Input';
import { formatSocialUrl } from '@ui/form/UrlInput/util';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';

import { SocialIcon } from './SocialIcons';

interface SocialInputGroupProps extends InputProps {
  bg?: string;
  value: string;
  index?: number;
  isReadOnly?: boolean;
  leftElement?: React.ReactNode;
}

export const SocialInput = memo(
  ({
    bg,
    value,
    onBlur,
    leftElement,
    isReadOnly,
    ...rest
  }: SocialInputGroupProps) => {
    const [isFocused, setIsFocused] = useState(false);
    const [isHovered, setIsHovered] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);

    const href = value?.startsWith('http') ? value : `https://${value}`;
    const formattedUrl = formatSocialUrl(value);

    const handleFocus = () => {
      if (isReadOnly) return;
      setIsFocused(true);
      inputRef.current?.focus();
    };

    const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
      onBlur?.(e);
      setIsFocused(false);
    };

    return (
      <>
        <div
          className='w-full '
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
        >
          <div className='relative w-full'>
            <InputGroup
              className={cn(
                isHovered
                  ? 'border-b border-transparent hover:border-transparent hover:border-b-none text-md focus-whithin:hover:border-b focus-whithin:hover:border-transparent focus-whithin:border-b focus-whithin:border-transparent'
                  : '',
              )}
            >
              {leftElement && (
                <LeftElement className='mb-[2px]'>
                  <SocialIcon url={value}>{leftElement}</SocialIcon>
                </LeftElement>
              )}
              <Input
                readOnly={isReadOnly}
                value={isFocused ? value : ''}
                ref={inputRef}
                onFocus={handleFocus}
                onBlur={handleBlur}
                className={
                  'border-b border-transparent hover:border-transparent hover:border-b-none text-md focus:hover:border-b focus:hover:border-transparent focus:border-b focus:border-transparent'
                }
                {...rest}
              />
            </InputGroup>

            {!isFocused && !!value && (
              <div className='w-full h-full'>
                <div
                  className={
                    'w-[calc(100%-30px] items-center absolute h-full top-[6px] left-7  hover:outline-none border-b border-transparent'
                  }
                >
                  <p
                    className='top-0 text-base cursor-auto' //try to align the text with the icon
                    onClick={handleFocus}
                  >
                    {formattedUrl}
                  </p>
                  {isHovered && (
                    <Link
                      to={href}
                      target='_blank'
                      className='cursor-pointer absolute top-0 -right-[30px] text-gray-500'
                    >
                      <IconButton
                        size='xs'
                        className='hover:bg-gray-200 '
                        variant='ghost'
                        colorScheme='gray'
                        aria-label='social link'
                        icon={<LinkExternal02 className='text-gray-500' />}
                      />
                    </Link>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </>
    );
  },
);
