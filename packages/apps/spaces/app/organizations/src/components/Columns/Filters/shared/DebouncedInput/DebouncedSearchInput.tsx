'use client';
import {
  memo,
  useRef,
  useState,
  useEffect,
  forwardRef,
  ChangeEvent,
} from 'react';

import { Input } from '@ui/form/Input';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';
import { SearchSm } from '@ui/media/icons/SearchSm';
import {
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from '@ui/form/InputGroup';

interface DebouncedInputProps {
  value: string;
  onChange: (value: string) => void;
}

export const DebouncedSearchInput = memo(
  forwardRef<HTMLInputElement, DebouncedInputProps>(
    ({ value: _value, onChange }, ref) => {
      const timeout = useRef<NodeJS.Timeout>();
      const [value, setValue] = useState(() => _value);

      useEffect(() => {
        return () => {
          if (timeout.current) {
            clearTimeout(timeout.current);
          }
        };
      }, []);

      const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        setValue(e.target.value);

        if (timeout.current) {
          clearTimeout(timeout.current);
        }

        timeout.current = setTimeout(() => {
          onChange(e.target.value);
        }, 250);
      };

      const handleClear = () => {
        setValue('');
        onChange('');
      };

      return (
        <InputGroup>
          <InputLeftElement w='fit-content'>
            <SearchSm color='gray.500' />
          </InputLeftElement>
          <Input
            pl='6'
            ref={ref}
            value={value}
            variant='flushed'
            placeholder='Search'
            onChange={handleChange}
          />
          <InputRightElement w='fit-content'>
            <IconButton
              size='xs'
              variant='ghost'
              onClick={handleClear}
              aria-label='search organization'
              icon={<Delete color='gray.500' />}
            />
          </InputRightElement>
        </InputGroup>
      );
    },
  ),
);
