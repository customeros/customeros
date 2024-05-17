import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

export const Image = observer(
  ({ src, ...props }: React.ImgHTMLAttributes<HTMLImageElement>) => {
    const store = useStore();

    useEffect(() => {
      if (!src || src?.startsWith('http') || src?.startsWith('blob')) return;

      // console.log('descarc', src);
      store.files.download(src);

      () => {
        store.files.clear(src);
      };
    }, [src]);

    if (src?.startsWith('http') || src?.startsWith('blob')) {
      return <img src={src} {...props} />;
    }

    return (
      <img src={src ? store.files.values.get(src) : undefined} {...props} />
    );
  },
);
