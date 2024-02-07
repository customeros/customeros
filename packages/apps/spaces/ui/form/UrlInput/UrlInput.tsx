import React, { memo, useRef, useState } from 'react';

import { FormLabel, FormLabelProps, VisuallyHidden } from '@chakra-ui/react';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { InputGroup } from '@ui/form/InputGroup';
import { formatSocialUrl } from '@ui/form/UrlInput/util';
import { Input, InputProps, FormControl } from '@ui/form/Input';

interface UrlInputProps extends InputProps {
  value: string;
  label?: string;
  isLabelVisible?: boolean;
  labelProps?: FormLabelProps;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
}

export const UrlInput = memo(
  ({
    value,
    onBlur,
    isLabelVisible,
    labelProps,
    label,
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
      <FormControl>
        {isLabelVisible ? (
          <FormLabel {...labelProps}>{label}</FormLabel>
        ) : (
          <VisuallyHidden>
            <FormLabel>{label}</FormLabel>
          </VisuallyHidden>
        )}

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
              w='100%'
              position='absolute'
              _hover={{
                '&:after': {
                  content: "''",
                  position: 'absolute',
                  width: '100%',
                  height: '1px',
                  bg: 'gray.300',
                  bottom: 0,
                },
              }}
            >
              <Text
                mr='-2px'
                mt='-2px'
                cursor='auto'
                onClick={handleFocus}
                noOfLines={1}
                w='100%'
              >
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
      </FormControl>
    );
  },
);
