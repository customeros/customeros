import {
  useRef,
  useState,
  useEffect,
  forwardRef,
  ChangeEvent,
  useImperativeHandle,
} from 'react';

import { Input } from '@ui/form/Input/Input';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';
import { SearchSm } from '@ui/media/icons/SearchSm';
import {
  InputGroup,
  LeftElement,
  RightElement,
} from '@ui/form/InputGroup/InputGroup';

interface DebouncedInputProps {
  value: string;
  placeholder?: string;
  onChange: (value: string) => void;
  onDisplayChange?: (value: string) => void;
}

export const DebouncedSearchInput = forwardRef<
  HTMLInputElement,
  DebouncedInputProps
>(({ value, onChange, onDisplayChange, placeholder }, ref) => {
  const timeout = useRef<NodeJS.Timeout>();
  const innerRef = useRef<HTMLInputElement>(null);
  const [displayValue, setDisplayValue] = useState(value);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    onDisplayChange?.(e.target.value);
    setDisplayValue(e.target.value);

    if (timeout.current) {
      clearTimeout(timeout.current);
    }

    timeout.current = setTimeout(() => {
      onChange(e.target.value);
    }, 250);
  };

  const handleClear = () => {
    onChange('');
    onDisplayChange?.('');
    setDisplayValue('');
    innerRef?.current?.focus();
  };

  useImperativeHandle(ref, () => innerRef.current!, []);

  useEffect(() => {
    return () => {
      if (timeout.current) {
        clearTimeout(timeout.current);
      }
    };
  }, []);

  useEffect(() => {
    innerRef.current?.focus();
  }, []);

  return (
    <InputGroup>
      <LeftElement>
        <SearchSm color='gray.500' />
      </LeftElement>
      <Input
        size='sm'
        ref={innerRef}
        variant='flushed'
        autoComplete='off'
        value={displayValue}
        onChange={handleChange}
        placeholder={placeholder ?? 'Search'}
        className='border-transparent focus:border-transparent focus:hover:border-transparent hover:border-transparent'
      />
      {value.length && (
        <RightElement>
          <IconButton
            size='xs'
            variant='ghost'
            onClick={handleClear}
            aria-label='search organization'
            icon={<Delete color='gray.500' />}
          />
        </RightElement>
      )}
    </InputGroup>
  );
});
