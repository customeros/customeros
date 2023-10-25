import {
  FloatingWrapper,
  useAttrs,
  useChainedCommands,
  useCurrentSelection,
} from '@remirror/react';
import { ChangeEvent, useCallback, useEffect, useState } from 'react';
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
  const [linkShortcut] = useState<ShortcutHandlerProps | undefined>();
  const [isEditing, setIsEditing] = useState(false);
  const chain = useChainedCommands();
  const { to, from } = useCurrentSelection();
  const url = (useAttrs().link()?.href as string) ?? '';
  const [href, setHref] = useState<string>(url);

const onRemove = useCallback(() => chain.removeLink().focus().run(), [chain, cancelHref]);

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

const cancelHref = useCallback(() => {
  const range = linkShortcut;
  setHref('');
  chain.focus(range?.to ?? to).run();
}, [chain, linkShortcut, to]);

  return (
    <>
      <FloatingWrapper
        positioner='selection'
        placement='bottom-start'
        enabled={isEditing}
        displayArrow
        hideWhenInvisible
      >
        {isEditing && (
          <Flex
            style={{ zIndex: 999, position: 'absolute' }}
            alignItems='center'
            sx={{
              '&': {
                position: 'relative',
                paddingY: 2,
                paddingX: 3,
                borderRadius: '8px',
                bg: 'gray.700',
              },

              '&:before': {
                content: "''",
                width: 0,
                height: 0,
                borderStyle: 'solid',
                borderWidth: '9px 10px 9px 0',
                borderColor: 'transparent #344054 transparent transparent',
                display: 'inline-block',
                verticalAlign: 'middle',
                marginRight: '5px',
                position: 'absolute',
                top: '-10px',
                transform: 'rotate(91deg)',
                left: '25px',
              },
            }}
          >
            <IconButton
              size='xs'
              variant='ghost'
              aria-label='Go to url'
              disabled={!href}
              onClick={() => {
window.open(getExternalUrl(href), '_blank', 'noopener noreferrer');
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
                  icon={<Check color='gray.400' />}
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
                  icon={<Trash01 color='gray.400' />}
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
