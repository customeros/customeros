import React, {
  useRef,
  useState,
  useEffect,
  ChangeEvent,
  useCallback,
} from 'react';

import { ShortcutHandlerProps } from '@remirror/extension-link';
import {
  useAttrs,
  useChainedCommands,
  useCurrentSelection,
} from '@remirror/react';

import { Input } from '@ui/form/Input/Input';
import { Check } from '@ui/media/icons/Check';
import { Link01 } from '@ui/media/icons/Link01';
import { Trash01 } from '@ui/media/icons/Trash01';
import { getExternalUrl } from '@utils/getExternalLink';
import { Divider } from '@ui/presentation/Divider/Divider';
import { IconButton } from '@ui/form/IconButton/IconButton';

function getTransformValues(transformStyle: string) {
  if (!transformStyle) return [0, 0];
  const match = transformStyle.match(/translate\((.*)px, (.*)px\)/);

  return match ? [Number(match[1]), Number(match[2])] : [0, 0];
}
interface LinkComponentProps {
  // Declare the type for other required props here. Replace "any" with the correct type.
  isEditing: boolean;
}

//TODO:before merge check if the design is correct
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
    <div
      ref={ref}
      className='flex items-center relative py-0 px-3 rounded-lg bg-gray-700'
    >
      <IconButton
        className='mr-2 rounded-sm hover:bg-gray-600 hover:text-gray-25'
        size='xs'
        variant='ghost'
        aria-label='Go to url'
        isDisabled={!href}
        onClick={() => {
          window.open(getExternalUrl(href), '_blank', 'noopener noreferrer');
        }}
        icon={<Link01 className='text-inherit' />}
      />
      <Input
        style={{
          background: 'gray.700',
        }}
        className='text-ellipsis overflow-hidden whitespace-nowrap bg-gray-700 !text-gray-25 focus-visible:outline-none placeholder:text-gray-400'
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
        <div className='flex items-center'>
          <IconButton
            className='mr-2 ml-2 hover:text-gray-25 hover:bg-gray-600 hover:text-gray-25'
            size='xs'
            variant='ghost'
            aria-label='Save'
            onClick={submitHref}
            icon={<Check className='text-inherit' />}
          />

          <Divider className='transform  border-l-[1px] border-gray-400 h-[14px]' />

          <IconButton
            className='ml-2 hover:text-gray-25 hover:bg-gray-600 hover:text-gray-25'
            size='xs'
            variant='ghost'
            aria-label='Remove link'
            onClick={() => {
              onRemove();
              cancelHref();
            }}
            icon={<Trash01 className='text-inherit' />}
          />
        </div>
      )}
    </div>
  );
};
