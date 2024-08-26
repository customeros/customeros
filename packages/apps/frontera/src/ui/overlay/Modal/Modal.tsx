import type {
  DialogProps,
  DialogTitleProps,
  DialogCloseProps,
  DialogContentProps,
  DialogTriggerProps,
  DialogOverlayProps,
} from '@radix-ui/react-alert-dialog';

import { forwardRef } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva } from 'class-variance-authority';
import * as Dialog from '@radix-ui/react-dialog';

import { cn } from '@ui/utils/cn';
import {
  FeaturedIcon,
  FeaturedIconStyleProps,
} from '@ui/media/Icon/FeaturedIcon';
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

export const ModalOverlay = forwardRef<HTMLDivElement, DialogOverlayProps>(
  ({ className, ...props }, ref) => {
    return (
      <Dialog.Overlay
        ref={ref}
        className={twMerge(
          'z-10 backdrop-brightness-[.55] data-[state=open]:animate-overlayShow fixed inset-0 cursor-pointer overflow-y-auto top-0 left-0 bottom-0 right-0 h-[100vh]',
          className,
        )}
        {...props}
      />
    );
  },
);

export const ModalHeader = ({
  children,
  className,
  ...props
}: DialogTitleProps) => {
  return (
    <Dialog.Title className={twMerge('px-6 pt-6 pb-4', className)} {...props}>
      {children}
    </Dialog.Title>
  );
};

export const ModalClose = (props: DialogCloseProps) => {
  return <Dialog.Close {...props} />;
};

export const ModalCloseButton = (props: DialogCloseProps) => {
  return <Dialog.Close asChild {...props}></Dialog.Close>;
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

export const ModalContent = forwardRef<
  HTMLDivElement,
  DialogContentProps & { placement?: 'center' | 'top' }
>(({ children, className, placement = 'top', ...props }, ref) => {
  return (
    <Dialog.Content
      ref={ref}
      className={twMerge(modalContentVariant({ placement, className }))}
      {...props}
    >
      {children}
    </Dialog.Content>
  );
});

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
    <div className={twMerge('p-6', className)} {...props}>
      {children}
    </div>
  );
};

export const ModalTrigger = (props: DialogTriggerProps) => {
  return <Dialog.Trigger {...props} />;
};

export const ModalFeaturedContent = forwardRef<
  HTMLDivElement,
  DialogTitleProps
>(({ className, ...props }, ref) => {
  return (
    <ModalContent
      ref={ref}
      className={cn(`rounded-2xl `, className)}
      style={{
        backgroundPositionX: '1px',
        backgroundPositionY: '-7px',
      }}
      {...props}
    >
      {props.children}
    </ModalContent>
  );
});

export const ModalFeaturedHeader = ({
  featuredIcon,
  featuredIconProps,
  ...props
}: DialogTitleProps & {
  featuredIcon?: React.ReactElement;
  featuredIconProps?: FeaturedIconStyleProps;
}) => {
  return (
    <ModalHeader {...props}>
      <FeaturedIcon
        size={featuredIconProps?.size ?? 'lg'}
        colorScheme={featuredIconProps?.colorScheme ?? 'primary'}
        className={cn('ml-[12px] mt-1 mb-[31px]', featuredIconProps?.className)}
      >
        {featuredIcon}
      </FeaturedIcon>
      {props.children}
    </ModalHeader>
  );
};
