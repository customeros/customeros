import { observer } from 'mobx-react-lite';
import { AddWebsiteToTrackUsecase } from '@domain/usecases/agents/capabilities/add-website-to-track.usecase';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
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
  return (
    <>
      <div>
        <h1 className='text-sm font-medium mb-4'>
          Track and identify website visitors
        </h1>

        <div className='flex flex-col gap-1'>
          <label className='text-sm font-medium'>Websites to track</label>
          <Button
            variant='ghost'
            className='w-fit'
            onClick={() => usecase.open()}
          >
            Add website
          </Button>
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
                variant='unstyled'
                placeholder='website'
                onChange={(e) => usecase.setWebsite(e.target.value)}
              />
            </AlertDialogBody>
            <AlertDialogFooter>
              <Button onClick={usecase.close}>Cancel</Button>
              <Button colorScheme='primary' onClick={usecase.execute}>
                Add website
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
});
