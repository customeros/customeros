import type {
  DialogProps,
  DialogTitleProps,
  DialogCloseProps,
  DialogContentProps,
  DialogTriggerProps,
  DialogOverlayProps,
} from '@radix-ui/react-alert-dialog';

import { twMerge } from 'tailwind-merge';
import { cva } from 'class-variance-authority';
import * as Dialog from '@radix-ui/react-dialog';

import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton/IconButton';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
  type ScrollAreaViewportProps,
} from '@ui/utils/ScrollArea';

export const Modal = (props: DialogProps) => {
  return <Dialog.Root {...props} />;
};

export const ModalPortal = Dialog.Portal;

export const ModalOverlay = ({ className, ...props }: DialogOverlayProps) => {
  return (
    <Dialog.Overlay
      className={twMerge(
        'backdrop-brightness-[.55] data-[state=open]:animate-overlayShow fixed inset-0 z-10 cursor-pointer overflow-y-auto top-0 left-0 bottom-0 right-0 h-[100vh]',
        className,
      )}
      {...props}
    />
  );
};

export const ModalHeader = ({
  children,
  className,
  ...props
}: DialogTitleProps) => {
  return (
    <Dialog.Title className={twMerge('px-6 pt-6 pb-5', className)} {...props}>
      {children}
    </Dialog.Title>
  );
};

export const ModalClose = (props: DialogCloseProps) => {
  return <Dialog.Close {...props} />;
};

export const ModalCloseButton = (props: DialogCloseProps) => {
  return (
    <Dialog.Close asChild {...props}>
      <IconButton
        size='lg'
        variant='ghost'
        colorScheme='gray'
        className='absolute top-4 right-4'
        icon={<XClose boxSize={5} className='w-5 h-5' />}
        aria-label='Close modal'
      />
    </Dialog.Close>
  );
};

const modalContentVariant = cva(
  'z-10 fixed left-[50%] w-[90vw] max-w-[450px] translate-x-[-50%] rounded-[6px] bg-white shadow-xl focus:outline-none data-[state=open]:will-change-auto',
  {
    variants: {
      placement: {
        center: [
          'top-[50%]',
          'translate-y-[-50%]',
          'data-[state=open]:animate-contentShowCenter',
        ],
        top: ['data-[state=open]:animate-contentShowTop', 'top-[4%]'],
      },
    },
    defaultVariants: {
      placement: 'top',
    },
  },
);

export const ModalContent = ({
  children,
  className,
  placement = 'top',
  ...props
}: DialogContentProps & { placement?: 'center' | 'top' }) => {
  return (
    <Dialog.Content
      className={twMerge(modalContentVariant({ placement, className }))}
      {...props}
    >
      {children}
    </Dialog.Content>
  );
};

export const ModalBody = ({
  children,
  className,
}: React.HTMLAttributes<HTMLDivElement>) => {
  return <div className={twMerge('px-6', className)}>{children}</div>;
};

export const ModalScrollBody = ({
  children,
  className,
  ...props
}: ScrollAreaViewportProps) => {
  return (
    <ScrollAreaRoot>
      <ScrollAreaViewport className={twMerge('px-6', className)} {...props}>
        {children}
      </ScrollAreaViewport>
      <ScrollAreaScrollbar orientation='vertical'>
        <ScrollAreaThumb />
      </ScrollAreaScrollbar>
    </ScrollAreaRoot>
  );
};

export const ModalFooter = ({
  children,
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => {
  return (
    <div className={twMerge('px-6 pb-6 pt-8', className)} {...props}>
      {children}
    </div>
  );
};

export const ModalTrigger = (props: DialogTriggerProps) => {
  return <Dialog.Trigger {...props} />;
};
