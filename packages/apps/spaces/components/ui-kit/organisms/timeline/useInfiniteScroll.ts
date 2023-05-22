import { useEffect, useRef } from 'react';
// TODO remove when list is virtualised
interface Options {
  callback: () => void;
  element: any;
  isFetching: boolean;
}

export const useInfiniteScroll = ({
  callback,
  element,
  isFetching,
}: Options) => {
  const observer = useRef<IntersectionObserver>();

  useEffect(() => {
    if (!element) {
      return;
    }

    observer.current = new IntersectionObserver(
      (entries) => {
        if (!isFetching && entries[0].isIntersecting) {
          callback();
        }
      },
      {
        rootMargin: '50%',
        threshold: [1, 1],
      },
    );
    observer.current?.observe(element.current);

    return () => observer.current?.disconnect();
  }, [callback, isFetching, element]);
};
