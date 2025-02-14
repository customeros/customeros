import { useRef, useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { AddWebsiteToTrackUsecase } from '@domain/usecases/agents/capabilities/add-website-to-track.usecase';

import { Icon } from '@ui/media/Icon';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  AlertDialogCloseIconButton,
} from '@ui/overlay/AlertDialog/AlertDialog';

const SCRIPT = `<script id="customeros-tracker" type="text/javascript">
  (function (c, u, s, t, o, m, e, r, O, S) { 
  var customerOS = document.createElement(s); 
  customerOS.src = u; 
  customerOS.async = true; 
  (document.body || document.head).appendChild(customerOS); 
})(window, "https://app.customeros.ai/analytics-0.1.js", "script");</script>`;

export const NewWebSessionListener = observer(() => {
  const { id } = useParams<{ id: string }>();

  const usecase = useMemo(() => new AddWebsiteToTrackUsecase(id!), [id]);
  const inputRef = useRef<HTMLInputElement>(null);

  const [_, copyToClipboard] = useCopyToClipboard();

  useKey('Escape', usecase.close, {
    when: usecase.isOpen,
  });

  return (
    <>
      <div>
        <h2 className='text-sm font-medium mb-4'>
          Track and identify website visitors
        </h2>

        {usecase.listenerErrors && (
          <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] mb-4'>
            <Icon stroke='none' className='mr-2' name='dot-single' />
            <span className='text-sm'>{usecase.listenerErrors}</span>
          </div>
        )}

        <div className='flex flex-col gap-1'>
          <p className='text-sm font-medium'>Websites to track</p>

          {usecase.websites.map((website) => (
            <div key={website} className='flex items-center group'>
              <div className='mx-2'>
                <Icon
                  stroke='none'
                  name='dot-single'
                  className='text-grayModern-500'
                />
              </div>
              <p className='text-sm cursor-default'>{website}</p>

              <Menu modal={false}>
                <MenuButton asChild>
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    aria-label='more'
                    icon={<Icon name='dots-vertical' />}
                    className='ml-2 invisible group-hover:visible'
                  />
                </MenuButton>
                <MenuList>
                  <MenuItem
                    onClick={() => {
                      usecase.executeRemove(website);
                    }}
                  >
                    <Icon name='x-circle' className='text-grayModern-500' />
                    Remove
                  </MenuItem>
                </MenuList>
              </Menu>
            </div>
          ))}

          <Button
            size='xs'
            variant='ghost'
            className='w-fit'
            onClick={usecase.open}
            colorScheme={'primary'}
            leftIcon={<Icon name='plus-circle' />}
          >
            Add website
          </Button>
          {usecase.websitesError.length > 0 && (
            <p className='text-sm text-error-500 ml-2'>
              {usecase.websitesError}
            </p>
          )}
        </div>

        <div className='w-full h-[1px] bg-grayModern-200 my-4' />

        <div>
          <h2 className='text-sm font-medium mb-1'>Code snippet</h2>
          <p className='pb-2 text-sm'>{`Place the following code in the <HEAD> section of your website:`}</p>
          <div className='px-3 py-2 rounded-md bg-grayModern-100 flex items-baseline'>
            <pre className='text-sm font-sticky whitespace-pre-wrap'>
              {`<script id="customeros-tracker" type="text/javascript">`}
              <br />
              &nbsp;&nbsp;{`(function (c, u, s, t, o, m, e, r, O, S) {`} <br />
              &nbsp;&nbsp;{`var customerOS = document.createElement(s);`} <br />
              &nbsp;&nbsp;{`customerOS.src = u;`} <br />
              &nbsp;&nbsp;{`customerOS.async = true;`} <br />
              &nbsp;&nbsp;
              {`(document.body || document.head).appendChild(customerOS);`}{' '}
              <br />
              {`})(window, "https://app.customeros.ai/analytics-0.1.js", "script");`}
              {`</script>`}
            </pre>
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='copy'
              icon={<Icon name='copy-03' />}
              onClick={() => copyToClipboard(SCRIPT, 'Code copied')}
            />
          </div>
        </div>
      </div>

      <AlertDialog isOpen={usecase.isOpen} onClose={usecase.close}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogCloseIconButton />
            <AlertDialogHeader className='font-medium'>
              Add a website to track
            </AlertDialogHeader>
            <AlertDialogBody>
              <p className='text-sm mb-2'>{`Once added, remember to place the tracking code in the <HEAD> section of this website`}</p>
              <Input
                size='sm'
                autoFocus
                ref={inputRef}
                variant='unstyled'
                placeholder='Website'
                invalid={usecase.isInvalid}
                onChange={(e) => usecase.setWebsite(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    usecase.executeAdd({
                      onInvalid: () => inputRef.current?.focus(),
                    });
                  }

                  if (e.key === 'Escape') {
                    usecase.close();
                  }
                }}
              />
              {usecase.isInvalid && (
                <p className='text-xs text-error-500'>
                  {usecase.validationError}
                </p>
              )}
            </AlertDialogBody>
            <AlertDialogFooter>
              <Button onClick={usecase.close}>Cancel</Button>
              <Button
                colorScheme='primary'
                onClick={() =>
                  usecase.executeAdd({
                    onInvalid: () => inputRef.current?.focus(),
                  })
                }
              >
                Add website
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
});
