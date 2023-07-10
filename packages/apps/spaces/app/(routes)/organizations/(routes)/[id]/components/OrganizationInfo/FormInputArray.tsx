import { useState, useRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  InputGroup,
  InputLeftElement,
  InputGroupProps,
} from '@ui/form/InputGroup';
import { Flex } from '@ui/layout/Flex';
import { Input, InputProps } from '@ui/form/Input';
import { Icons } from '@ui/media/Icon';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';

interface FormInputArrayProps extends InputGroupProps {
  name: string;
  formId: string;
  leftElement?: React.ReactNode;
}

export const FormInputArray = ({
  name,
  formId,
  leftElement,
  ...rest
}: FormInputArrayProps) => {
  const { getInputProps } = useField(name, formId);
  const { value: values, onChange, onBlur } = getInputProps();

  const newInputRef = useRef<HTMLInputElement>(null);
  const [newValue, setNewValue] = useState('');

  const handleChange =
    (index: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
      const next = [...values];
      next[index] = e.target.value;
      onChange(next);
    };

  const handleBlur =
    (index: number) => (e: React.FocusEvent<HTMLInputElement>) => {
      if (!e.target.value) {
        const next = [...values];
        next.splice(index, 1);
        onBlur?.(next);
      } else {
        onBlur?.(values);
      }
    };

  const handleRemoveKeyDown =
    (index: number) => (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Backspace' && !values[index]) {
        const next = [...values];
        next.splice(index, 1);
        onBlur?.(next);
        newInputRef.current?.focus();
      }
    };

  const handleAddKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      if (newValue) {
        onBlur([...values, newValue]);
        setNewValue('');
      }
    }
  };

  const handleAdd = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewValue(e.target.value);
  };

  const handleAddBlur = () => {
    if (newValue) {
      onBlur?.([...values, newValue]);
      setNewValue('');
    }
  };

  return (
    <>
      {(values as string[])?.map((v, index) => (
        <SocialInput
          key={index}
          value={v}
          index={index}
          leftElement={leftElement}
          onBlur={handleBlur(index)}
          onChange={handleChange(index)}
          onKeyDown={handleRemoveKeyDown(index)}
        />
      ))}

      <InputGroup {...rest}>
        {leftElement && (
          <InputLeftElement>
            <SocialIcon url={newValue}>{leftElement}</SocialIcon>
          </InputLeftElement>
        )}
        <Input
          value={newValue}
          ref={newInputRef}
          onChange={handleAdd}
          onBlur={handleAddBlur}
          onKeyDown={handleAddKeyDown}
          {...rest}
        />
      </InputGroup>
    </>
  );
};

const SocialIcon = ({
  children,
  url,
}: React.PropsWithChildren<{ url: string }>) => {
  const knownUrl = isKnownUrl(url);

  if (knownUrl === 'twitter')
    return <Icons.Twitter viewBox='0 0 32 32' strokeWidth='0' />;
  if (knownUrl === 'linkedin')
    return <Icons.Linkedin viewBox='0 0 32 32' strokeWidth='0' />;
  return children;
};

interface SocialInputGroupProps extends InputProps {
  index: number;
  value: string;
  leftElement?: React.ReactNode;
}

const SocialInput = ({
  value,
  onBlur,
  leftElement,
  ...rest
}: SocialInputGroupProps) => {
  const [isFocused, setIsFocused] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const href = value.startsWith('http') ? value : `https://${value}`;
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
    <InputGroup>
      {leftElement && (
        <InputLeftElement>
          <SocialIcon url={value}>{leftElement}</SocialIcon>
        </InputLeftElement>
      )}
      <Input value={value} ref={inputRef} onBlur={handleBlur} {...rest} />
      {!isFocused && (
        <Flex
          bg='white'
          w={'full'}
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
};

function isKnownUrl(input: string) {
  const url = input.trim().toLowerCase();
  if (url.includes('twitter')) return 'twitter';
  if (url.includes('linkedin')) return 'linkedin';
}

function formatSocialUrl(value: string) {
  let url = value;

  if (url.startsWith('http')) {
    url = url.replace('https://', '');
  }
  if (url.startsWith('www')) {
    url = url.replace('www.', '');
  }
  if (url.includes('twitter')) {
    url = url.replace('twitter.com', '');
  }
  if (url.includes('linkedin')) {
    url = url.replace('linkedin.com/in', '');
  }

  return url;
}
