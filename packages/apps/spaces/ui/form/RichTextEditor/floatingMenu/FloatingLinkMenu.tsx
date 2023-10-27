import {
  FloatingWrapper,
  useAttrs,
  useChainedCommands,
  useCurrentSelection,
} from '@remirror/react';
import { ChangeEvent, useCallback, useEffect, useRef, useState } from 'react';
import { ShortcutHandlerProps } from '@remirror/extension-link';
import { Input } from '@ui/form/Input';
import { Flex } from '@ui/layout/Flex';
import { Divider } from '@ui/presentation/Divider';
import { Link01 } from '@ui/media/icons/Link01';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Check } from '@ui/media/icons/Check';
import { IconButton } from '@ui/form/IconButton';
import { getExternalUrl } from '@spaces/utils/getExternalLink';

export const FloatingLinkToolbar = () => {
  const ref = useRef<HTMLDivElement>(null);
  const [linkShortcut] = useState<ShortcutHandlerProps | undefined>();
  const [isEditing, setIsEditing] = useState(false);
  const chain = useChainedCommands();
  const { to, from } = useCurrentSelection();
  const url = (useAttrs().link()?.href as string) ?? '';
  const [href, setHref] = useState<string>(url);
  const cancelHref = useCallback(() => {
    const range = linkShortcut;
    setHref('');
    chain.focus(range?.to ?? to).run();
  }, [chain, linkShortcut, to]);
  const onRemove = useCallback(
    () => chain.removeLink().focus().run(),
    [chain, cancelHref],
  );

  useEffect(() => {
    if (isEditing && from === to) {
      setIsEditing(false);
    }
    if (!isEditing && from !== to) {
      setIsEditing(true);
    }
  }, [from, to, isEditing]);

  useEffect(() => {
    setHref(url as string);
  }, [url]);

  const submitHref = useCallback(() => {
    const range = linkShortcut;

    if (href === '') {
      chain.removeLink();
    } else {
      chain.updateLink({ href, auto: false }, range);
    }

    chain.focus(range?.to ?? to).run();
  }, [href, chain, setIsEditing, linkShortcut, to]);
  useEffect(() => {
    const el = ref.current;
    const parentEl = el?.parentElement?.parentElement?.parentElement;

    if (el && parentEl) {
      let elRect = el.getBoundingClientRect();
      const parentElRect = parentEl.getBoundingClientRect();
      let xValue = 0;
      let yValue = 0;

      // check and fix overflow at the top
      while (elRect.top < parentElRect.top) {
        yValue++;
        el.style.transform = `translateY(${yValue}px)`;
        elRect = el.getBoundingClientRect();
      }

      // check and fix overflow at the bottom
      while (elRect.bottom > parentElRect.bottom) {
        yValue--;
        el.style.transform = `translateY(${yValue}px)`;
        elRect = el.getBoundingClientRect();
      }

      // check and fix overflow on the left
      while (elRect.left < parentElRect.left) {
        xValue++;
        el.style.transform = `translateX(${xValue}px)`;
        elRect = el.getBoundingClientRect();
      }

      // check and fix overflow on the right
      while (elRect.right > parentElRect.right) {
        xValue--;
        el.style.transform = `translateX(${xValue}px)`;
        elRect = el.getBoundingClientRect();
      }
    }
  }, [isEditing, ref, from, to]);
  return (
    <>
      <FloatingWrapper
        positioner='selection'
        placement='auto'
        enabled={isEditing}
      >
        {isEditing && (
          <Flex
            ref={ref}
            className='test'
            alignItems='center'
            sx={{
              '&': {
                position: 'relative',
                paddingY: 2,
                paddingX: 3,
                borderRadius: '8px',
                bg: 'gray.700',
              },
            }}
          >
            <IconButton
              size='xs'
              variant='ghost'
              aria-label='Go to url'
              disabled={!href}
              onClick={() => {
                window.open(
                  getExternalUrl(href),
                  '_blank',
                  'noopener noreferrer',
                );
              }}
              icon={<Link01 color='gray.25' />}
              mr={2}
              borderRadius='sm'
              _hover={{ background: 'gray.600', color: 'gray.25' }}
            />

            <Input
              style={{
                background: 'gray.700',
              }}
              sx={{
                textOverflow: 'ellipsis',
                overflow: 'hidden',
                whiteSpace: 'nowrap',
                background: 'gray.700',
                fontSize: 'sm',
                color: 'gray.25',
                '&:focus-visible': { outline: 'none' },
                '&::placeholder': { color: 'gray.400' },
              }}
              tabIndex={1}
              placeholder='Paste or enter a link'
              onChange={(event: ChangeEvent<HTMLInputElement>) =>
                setHref(event.target.value)
              }
              value={href}
              onKeyDown={(event) => {
                const { key } = event;
                if (key === 'Enter') {
                  submitHref();
                  setHref('');
                }

                if (key === 'Escape') {
                  cancelHref();
                }
              }}
            />

            {href && (
              <Flex alignItems='center'>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Save'
                  onClick={submitHref}
                  color='gray.400'
                  icon={<Check color='inherit' />}
                  mr={2}
                  ml={2}
                  borderRadius='sm'
                  _hover={{ background: 'gray.600', color: 'gray.25' }}
                />

                <Divider
                  orientation='vertical'
                  borderLeft='1px solid'
                  borderLeftColor='gray.400 !important'
                  height='14px'
                />

                <IconButton
                  ml={2}
                  borderRadius='sm'
                  size='xs'
                  variant='ghost'
                  aria-label='Remove link'
                  onClick={() => {
                    onRemove();
                    cancelHref();
                  }}
                  color='gray.400'
                  icon={<Trash01 color='inherit' />}
                  _hover={{ background: 'gray.600', color: 'gray.25' }}
                />
              </Flex>
            )}
          </Flex>
        )}
      </FloatingWrapper>
    </>
  );
};
