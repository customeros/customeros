import * as ScrollArea from '@radix-ui/react-scroll-area';

export const ScrollAreaRoot = ScrollArea.Root;
ScrollAreaRoot.defaultProps = {
  className: 'w-full h-full overflow-hidden',
};

export type ScrollAreaRootProps = ScrollArea.ScrollAreaViewportProps;
export const ScrollAreaViewport = ScrollArea.Viewport;
ScrollAreaViewport.defaultProps = {
  className: 'h-full w-full',
};

export type ScrollAreaViewportProps = ScrollArea.ScrollAreaScrollbarProps;
export const ScrollAreaScrollbar = ScrollArea.Scrollbar;
ScrollAreaScrollbar.defaultProps = {
  className:
    'flex select-none touch-none p-0.5 bg-gray-100 transition-colors duration-[160ms] ease-out hover:bg-gray-200 data-[orientation=vertical]:w-2.5 data-[orientation=horizontal]:flex-col data-[orientation=horizontal]:h-2.5',
};

export type ScrollAreaThumbProps = ScrollArea.ScrollAreaThumbProps;
export const ScrollAreaThumb = ScrollArea.Thumb;
ScrollAreaThumb.defaultProps = {
  className:
    'flex-1 bg-gray-500 rounded-[10px] relative before:content-[""] before:absolute before:top-1/2 before:left-1/2 before:-translate-x-1/2 before:-translate-y-1/2 before:w-full before:h-full before:min-w-[44px] before:min-h-[44px]',
};
