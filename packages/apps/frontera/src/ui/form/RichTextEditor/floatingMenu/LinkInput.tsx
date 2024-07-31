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
import { Trash01 } from '@ui/media/icons/Trash01';
import { getExternalUrl } from '@utils/getExternalLink';
import { Divider } from '@ui/presentation/Divider/Divider';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';

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
      <Input
        tabIndex={1}
        value={href}
        placeholder='Paste or enter a link'
        onChange={(event: ChangeEvent<HTMLInputElement>) =>
          setHref(event.target.value)
        }
        className='text-ellipsis overflow-hidden whitespace-nowrap bg-gray-700 !text-gray-25 focus-visible:outline-none placeholder:text-gray-400'
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
        <div className='flex items-center '>
          <Divider className='transform  border-l-[1px] border-gray-400 h-[14px] mr-0.5' />

          <IconButton
            size='xs'
            variant='ghost'
            isDisabled={!href}
            aria-label='Go to url'
            className='hover:bg-gray-600 hover:text-gray-25'
            icon={<LinkExternal02 className='text-gray-400' />}
            onClick={() => {
              window.open(
                getExternalUrl(href),
                '_blank',
                'noopener noreferrer',
              );
            }}
          />
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Save'
            onClick={submitHref}
            icon={<Check className='text-gray-400' />}
            className='hover:bg-gray-600 hover:text-gray-25'
          />

          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Remove link'
            icon={<Trash01 className='text-gray-400' />}
            className='mr-0 hover:bg-gray-600 hover:text-gray-25'
            onClick={() => {
              onRemove();
              cancelHref();
            }}
          />
        </div>
      )}
    </div>
  );
};
