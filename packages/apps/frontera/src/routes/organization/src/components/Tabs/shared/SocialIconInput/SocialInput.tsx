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
  dataTest?: string;
  isReadOnly?: boolean;
  leftElement?: React.ReactNode;
}

export const SocialInput = memo(
  ({
    bg,
    value,
    onBlur,
    dataTest,
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
                ref={inputRef}
                onBlur={handleBlur}
                data-test={dataTest}
                readOnly={isReadOnly}
                onFocus={handleFocus}
                value={isFocused ? value : ''}
                className={
                  'border-b border-transparent hover:border-transparent hover:border-b-none text-md focus:hover:border-b focus:hover:border-transparent focus:border-b focus:border-transparent'
                }
                {...rest}
              />
            </InputGroup>

            {!isFocused && !!value && (
              <div className='h-full '>
                <div
                  className={
                    'items-center absolute w-[stretch] h-full top-[6px] left-7 hover:outline-none border-b border-transparent whitespace-nowrap'
                  }
                >
                  <p
                    onClick={handleFocus}
                    className='top-0 text-base cursor-auto overflow-hidden overflow-ellipsis'
                  >
                    {formattedUrl}
                  </p>
                  {isHovered && (
                    <Link
                      to={href}
                      target='_blank'
                      className='cursor-pointer absolute top-0 -right-[24px] text-gray-500'
                    >
                      <IconButton
                        size='xs'
                        variant='ghost'
                        colorScheme='gray'
                        aria-label='social link'
                        className='hover:bg-gray-200 '
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
