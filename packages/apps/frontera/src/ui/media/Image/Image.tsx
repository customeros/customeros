import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

interface ImageProps
  extends Omit<React.ImgHTMLAttributes<HTMLImageElement>, 'src'> {
  src?: string | null;
  fallbackSrc?: string;
}

export const Image = observer(({ src, fallbackSrc, ...props }: ImageProps) => {
  const store = useStore();

  useEffect(() => {
    if (
      !src ||
      src?.startsWith('http') ||
      src?.startsWith('blob') ||
      src?.startsWith('/')
    )
      return;

    store.files.download(src);

    () => {
      store.files.clear(src);
    };
  }, [src]);

  if (
    src?.startsWith('http') ||
    src?.startsWith('blob') ||
    src?.startsWith('/')
  ) {
    return <img src={src ?? fallbackSrc} {...props} />;
  }

  return (
    <img
      src={src ? store.files.values.get(src) ?? fallbackSrc : fallbackSrc}
      {...props}
    />
  );
});
