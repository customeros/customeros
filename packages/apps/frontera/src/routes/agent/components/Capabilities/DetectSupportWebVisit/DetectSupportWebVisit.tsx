import { useRef, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { DetectSupportWebVistUsecase } from '@domain/usecases/agents/capabilities/detect-support-web-vist';

import { Icon } from '@ui/media/Icon';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { AgentType } from '@shared/types/__generated__/graphql.types';
import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  AlertDialogCloseIconButton,
} from '@ui/overlay/AlertDialog/AlertDialog';
export const DetectSupportWebVisit = observer(() => {
  const store = useStore();
  const navigate = useNavigate();

  const { id } = useParams<{ id: string }>();
  const inputRef = useRef<HTMLInputElement>(null);
  const usecase = useMemo(() => new DetectSupportWebVistUsecase(id!), [id]);

  const agentByType = store.agents.getFirstAgentByType(
    AgentType.WebVisitIdentifier,
  );

  const handleClick = () => {
    if (!agentByType) {
      usecase.executeCreateWebVisitIdentifierAgent(
        AgentType.WebVisitIdentifier,
      );
      navigate(`/agents/${usecase.webVisitIdentifierAgentId}`);
    } else {
      navigate(`/agents/${agentByType?.id}`);
    }
  };

  return (
    <>
      <div className='flex flex-col gap-4'>
        <div>
          <p className='font-semibold text-sm'>Detect support web visit</p>

          {usecase.listenerErrors && (
            <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] mb-4'>
              <Icon stroke='none' className='mr-2' name='dot-single' />
              <span className='text-sm'>{usecase.listenerErrors}</span>
            </div>
          )}
          <p className='text-sm'>
            This agent’s goal relies on the results from the{' '}
            <span
              onClick={handleClick}
              className='font-semibold underline cursor-pointer'
            >
              Website visitor identifier
            </span>{' '}
            agent
          </p>
        </div>
        <div>
          <p className='font-semibold text-sm'>Pages to track</p>
          <p className='text-sm my-1'>
            Add pages you want to track and use to identify companies, such as
            support, documentation, or forum pages.
          </p>
          {usecase.supportUrls.map((url) => (
            <div key={url} className='flex items-center group'>
              <div className='mx-2'>
                <Icon
                  stroke='none'
                  name='dot-single'
                  className='text-grayModern-500'
                />
              </div>
              <p className='text-sm cursor-default'>{url}</p>
              <IconButton
                size='xxs'
                variant='ghost'
                aria-label='more'
                icon={<Icon name='x-circle' />}
                onClick={() => usecase.executeRemove(url)}
                className='ml-2 invisible group-hover:visible'
              />
            </div>
          ))}
          <Button
            size='xs'
            variant='ghost'
            colorScheme='primary'
            className='w-fit mt-1'
            onClick={usecase.open}
            leftIcon={<Icon name='plus-circle' />}
          >
            Page or keyword
          </Button>
        </div>
      </div>

      <AlertDialog isOpen={usecase.isOpen} onClose={usecase.close}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogCloseIconButton />
            <AlertDialogHeader className='font-medium'>
              Add a page or keyword
            </AlertDialogHeader>
            <AlertDialogBody>
              <p className='text-sm mb-2'>
                Add a specific help page URL you want to track or enter a
                keyword like ‘support’ to automatically include any page
                containing that word.
              </p>
              <Input
                size='sm'
                autoFocus
                ref={inputRef}
                variant='unstyled'
                invalid={usecase.isInvalid}
                placeholder='E.g pinkvoid.com or keyword like ‘support’'
                onChange={(e) => usecase.setPageToTrack(e.target.value)}
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
                Add
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
});
