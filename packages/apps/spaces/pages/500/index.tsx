import React from 'react';
import { ErrorPage } from '@spaces/organisms/error-page';

export const ServerError: React.FC = () => {
  return (
    <ErrorPage
      imageSrc={`/backgrounds/blueprint/server-error-1.webp`}
      blurredSrc={`/backgrounds/blueprint/server-error-1-blur.webp`}
      title='Oops!'
    >
      <>
        <p>We are sorry, but something went wrong on our end.</p>
        <p>
          Our team has been notified of the issue and is working to fix it as
          soon as possible.
        </p>
        <p>
          In the meantime, please try again later or contact our support team if
          the issue persists.
        </p>
        <p> Thank you for your patience and understanding.</p>
      </>
    </ErrorPage>
  );
};

export default ServerError;
