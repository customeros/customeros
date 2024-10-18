import React, { useRef, useMemo, useState } from 'react';

import { Input } from '@ui/form/Input';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';

import { SocialIcon } from './SocialIcons';
import { SocialInput } from './SocialInput';

interface SocialIconInputProps {
  name?: string;
  dataTest?: string;
  isReadOnly?: boolean;
  placeholder?: string;
  leftElement?: React.ReactNode;
  onCreate?: (value: string) => void;
  value?: { label: string; value: string }[];
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (
    e: React.ChangeEvent<HTMLInputElement>,
    newInputRef: React.RefObject<HTMLInputElement>,
  ) => void;
  onKeyDown?: (
    e: React.KeyboardEvent<HTMLInputElement>,
    newInputRef: React.RefObject<HTMLInputElement>,
  ) => void;
}

export const SocialIconInput = ({
  value,
  name = 'socialMedia',
  leftElement,
  isReadOnly,
  onBlur,
  dataTest,
  onChange,
  onCreate,
  onKeyDown,
  ...rest
}: SocialIconInputProps) => {
  const [socialIconValue, setSocialIconValue] = useState('');
  const _leftElement = useMemo(() => leftElement, [leftElement]);
  const newInputRef = useRef<HTMLInputElement>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange?.(e);
  };

  const handleBlur = (e: React.ChangeEvent<HTMLInputElement>) => {
    onBlur?.(e, newInputRef);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    onKeyDown?.(e, newInputRef);
  };

  const handleNewSocial = () => {
    const value = newInputRef.current?.value;

    if (!value) return;
    onCreate?.(value);
    newInputRef.current!.value = '';
    setSocialIconValue('');
  };

  return (
    <>
      {value?.map(({ value: v, label: l }) => (
        <SocialInput
          id={v}
          key={v}
          value={l}
          name={name}
          dataTest={dataTest}
          onBlur={handleBlur}
          onChange={handleChange}
          isReadOnly={isReadOnly}
          onKeyDown={handleKeyDown}
          leftElement={_leftElement}
        />
      ))}

      {!isReadOnly && (
        <InputGroup>
          {leftElement && (
            <LeftElement>
              <SocialIcon url={socialIconValue}>{leftElement}</SocialIcon>
            </LeftElement>
          )}
          <Input
            name={name}
            ref={newInputRef}
            dataTest={dataTest}
            onBlur={handleNewSocial}
            onChange={(e) => {
              setSocialIconValue(e.target.value);
            }}
            onKeyDown={(e) => {
              e.stopPropagation();

              if (e.key === 'Enter') {
                handleNewSocial?.();
              }
            }}
            className='border-b border-transparent hover:border-transparent hover:border-b-none text-md focus:hover:border-b focus:hover:border-transparent focus:border-b focus:border-transparent'
            {...rest}
          />
        </InputGroup>
      )}
    </>
  );
};
