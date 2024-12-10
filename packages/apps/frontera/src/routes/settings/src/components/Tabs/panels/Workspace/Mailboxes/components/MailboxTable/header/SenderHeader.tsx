import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';

export const SenderHeader = () => {
  const { open, onOpen, onClose } = useDisclosure();

  return (
    <>
      <div className='flex items-center gap-1'>
        <span className='text-sm'>Sender</span>
        <InfoCircle
          onClick={onOpen}
          className='size-3 text-gray-500 cursor-pointer hover:text-gray-700'
        />
      </div>

      <InfoDialog
        isOpen={open}
        label='Senders'
        onClose={onClose}
        onConfirm={onClose}
        confirmButtonLabel='Got it'
        body={
          <div className='text-sm'>
            We’ll use the sender’s first and last name to represent the mailbox.
            They’ll receive notifications when someone replies, and their
            responses will appear in the relevant organizational timeline
          </div>
        }
      />
    </>
  );
};
