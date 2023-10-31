import React, {
  ChangeEvent,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react';
import { Input } from '@ui/form/Input';
import { Flex } from '@ui/layout/Flex';
import { Divider } from '@ui/presentation/Divider';
import { Link01 } from '@ui/media/icons/Link01';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Check } from '@ui/media/icons/Check';
import { IconButton } from '@ui/form/IconButton';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import {
  useAttrs,
  useChainedCommands,
  useCurrentSelection,
} from '@remirror/react';
import { ShortcutHandlerProps } from '@remirror/extension-link';

function getTransformValues(transformStyle: string) {
  if (!transformStyle) return [0, 0];
  const match = transformStyle.match(/translate\((.*)px, (.*)px\)/);
  return match ? [Number(match[1]), Number(match[2])] : [0, 0];
}
interface LinkComponentProps {
  // Declare the type for other required props here. Replace "any" with the correct type.
  isEditing: boolean;
}

export const LinkComponent: React.FC<LinkComponentProps> = ({ isEditing }) => {
  const [linkShortcut] = useState<ShortcutHandlerProps | undefined>();

  const ref = useRef<HTMLDivElement>(null);
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
  }, [href, chain, linkShortcut, to]);

  useEffect(() => {
    const el = ref.current;
    const parentEl =
      el?.parentElement?.parentElement?.parentElement?.parentElement;
    if (el && parentEl) {
      const elRect = el.getBoundingClientRect();
      const parentElRect = parentEl.getBoundingClientRect();
      // eslint-disable-next-line prefer-const
      let [xValue, yValue] = getTransformValues(el.style.transform);

      if (elRect.top < parentElRect.top) {
        yValue = 8;
        el.style.transform = `translate(${xValue}px, ${yValue}px)`;
      }
      if (elRect.top > parentElRect.top) {
        yValue = 0;
        el.style.transform = `translate(${xValue}px, ${yValue}px)`;
      }
    }
  }, [isEditing, ref, from, to]);
  return (
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
          window.open(getExternalUrl(href), '_blank', 'noopener noreferrer');
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
  );
};
