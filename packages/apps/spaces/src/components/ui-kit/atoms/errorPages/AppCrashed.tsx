import React from 'react';

export const AppCrashed = ({ error, componentStack, resetError }: any) => (
  <>
    <div>Oops! App crashed!</div>
    <div>{error.toString()}</div>
    <div>{componentStack}</div>
    <button onClick={resetError}>Try again</button>
  </>
);
