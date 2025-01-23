import { useRef } from 'react';

import { observer } from 'mobx-react-lite';
import { AddWebsiteToTrackUsecase } from '@domain/usecases/agents/capabilities/add-website-to-track.usecase';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { Copy03 } from '@ui/media/icons/Copy03';
import { IconButton } from '@ui/form/IconButton';
import { XCircle } from '@ui/media/icons/XCircle';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
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

const usecase = new AddWebsiteToTrackUsecase();

export const WebsiteTrackerCapability = observer(() => {
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <>
      <div>
        <h2 className='text-sm font-medium mb-4'>
          Track and identify website visitors
        </h2>

        <div className='flex flex-col gap-1'>
          <p className='text-sm font-medium'>Websites to track</p>

          {usecase.websites.map((website) => (
            <div key={website} className='flex items-center group'>
              <div className='mx-2'>
                <DotSingle className='text-grayModern-500' />
              </div>
              <p className='text-sm cursor-default'>{website}</p>

              <Menu modal={false}>
                <MenuButton asChild>
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    aria-label='more'
                    icon={<DotsVertical />}
                    className='ml-2 invisible group-hover:visible'
                  />
                </MenuButton>
                <MenuList>
                  <MenuItem onClick={() => usecase.removeWebsite(website)}>
                    <XCircle className='text-grayModern-500' />
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
            leftIcon={<PlusCircle />}
          >
            Add website
          </Button>
        </div>

        <div className='w-full h-[1px] bg-grayModern-200 my-4' />

        <div>
          <h2 className='text-sm font-medium mb-1'>Code snippet</h2>
          <p className='pb-2'>{`Place the following code in the <HEAD> section of your website:`}</p>
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
              icon={<Copy03 />}
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
                placeholder='website'
                invalid={usecase.isInvalid}
                onChange={(e) => usecase.setWebsite(e.target.value)}
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
                  usecase.execute({
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
