import React, { useState } from 'react';
import { useRouter } from 'next/router';

export const NotFound = ({ error, componentStack, resetError }: any) => {
  const router = useRouter();
  return (
    <>
      <div>404 - not found</div>
      <div>{error.toString()}</div>
      <div>{componentStack}</div>
      <button onClick={() => router.push('/')}></button>
    </>
  );
};
