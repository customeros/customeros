/* eslint-disable react/no-unescaped-entities */
import React from 'react';
import { ErrorPage } from '../../components';

const blurredSrc = 'data:image/webp;base64,UklGRhwCAABXRUJQVlA4WAoAAAAgAAAABQAAAwAASUNDUMgBAAAAAAHIAAAAAAQwAABtbnRyUkdCIFhZWiAH4AABAAEAAAAAAABhY3NwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAA9tYAAQAAAADTLQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAlkZXNjAAAA8AAAACRyWFlaAAABFAAAABRnWFlaAAABKAAAABRiWFlaAAABPAAAABR3dHB0AAABUAAAABRyVFJDAAABZAAAAChnVFJDAAABZAAAAChiVFJDAAABZAAAAChjcHJ0AAABjAAAADxtbHVjAAAAAAAAAAEAAAAMZW5VUwAAAAgAAAAcAHMAUgBHAEJYWVogAAAAAAAAb6IAADj1AAADkFhZWiAAAAAAAABimQAAt4UAABjaWFlaIAAAAAAAACSgAAAPhAAAts9YWVogAAAAAAAA9tYAAQAAAADTLXBhcmEAAAAAAAQAAAACZmYAAPKnAAANWQAAE9AAAApbAAAAAAAAAABtbHVjAAAAAAAAAAEAAAAMZW5VUwAAACAAAAAcAEcAbwBvAGcAbABlACAASQBuAGMALgAgADIAMAAxADZWUDggLgAAANABAJ0BKgYABAADgFoliAJ0APSMe7cAAP6I+kEz7c6Pts7AUGA/7+fd60VkAAA=';
export const NotFound: React.FC = () => {
  return (
    <ErrorPage
      imageSrc={`/backgrounds/blueprint/not-found-4.webp`}
      blurredSrc={blurredSrc}
      title='404'
    >
      <>
        <p>We're sorry, but the page you're trying to access doesn't exist.</p>
        <p>
          This might be because you've entered an incorrect URL or the page was
          recently removed.
        </p>
        <p>
          You can try checking the spelling of the URL or navigate back to home
          page.
        </p>
      </>
    </ErrorPage>
  );
};

export default NotFound;
