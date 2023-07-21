import { useState, useRef, memo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Input, InputProps } from '@ui/form/Input';
import { InputGroup, InputLeftElement } from '@ui/form/InputGroup';

import { SocialIcon } from './SocialIcons';
import { formatSocialUrl } from '../util';

interface SocialInputGroupProps extends InputProps {
  index: number;
  value: string;
  leftElement?: React.ReactNode;
}

export const SocialInput = memo(
  ({ value, onBlur, leftElement, ...rest }: SocialInputGroupProps) => {
    const [isFocused, setIsFocused] = useState(false);
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
      <InputGroup onFocus={handleFocus}>
        {leftElement && (
          <InputLeftElement>
            <SocialIcon url={value}>{leftElement}</SocialIcon>
          </InputLeftElement>
        )}
        <Input value={value} ref={inputRef} onBlur={handleBlur} {...rest} />
        {!isFocused && !!value && (
          <Flex
            bg='white'
            w='calc(100% - 40px)'
            left='40px'
            align='center'
            position='absolute'
            h='calc(100% - 1px)'
          >
            <Text mr='1' cursor='auto' onClick={handleFocus}>
              {formattedUrl}
            </Text>
            <IconButton
              size='sm'
              as={Link}
              href={href}
              target='_blank'
              variant='ghost'
              aria-label='social link'
              icon={<Icons.LinkExternal2 color='gray.500' />}
            />
          </Flex>
        )}
      </InputGroup>
    );
  },
);
