/* eslint-disable react/no-unescaped-entities */
import React from 'react';
import { ErrorPage } from '../../src/components';

export const NotFound: React.FC = () => {
  return (
    <ErrorPage imageSrc={`/backgrounds/blueprint/not-found-4.webp`} title='404'>
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
