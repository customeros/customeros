import { memo, useRef, useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { InputGroup } from '@ui/form/InputGroup';
import { Input, InputProps } from '@ui/form/Input';

import { formatSocialUrl } from '../../../shared/util';

interface UrlInputProps extends InputProps {
  value: string;
}

export const UrlInput = memo(({ value, onBlur, ...rest }: UrlInputProps) => {
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
    <InputGroup
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <Input
        value={value}
        ref={inputRef}
        onFocus={handleFocus}
        onBlur={handleBlur}
        {...rest}
      />
      {!isFocused && !!value && (
        <Flex
          bg='gray.25'
          h='100%'
          align='center'
          position='absolute'
          w='calc(100% - 40px)'
        >
          <Text mr='1' cursor='auto' onClick={handleFocus}>
            {formattedUrl}
          </Text>
          {isHovered && (
            <IconButton
              size='xs'
              as={Link}
              href={href}
              target='_blank'
              variant='ghost'
              aria-label='social link'
              icon={<Icons.LinkExternal2 color='gray.500' boxSize='4' />}
            />
          )}
        </Flex>
      )}
    </InputGroup>
  );
});
