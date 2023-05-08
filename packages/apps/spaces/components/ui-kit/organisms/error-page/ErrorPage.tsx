import React, { ReactNode } from 'react';
import Image from 'next/image';
import { Button } from '../../atoms';
import { useRouter } from 'next/router';

import styles from './error-page.module.scss';

interface ErrorPageProps {
  imageSrc: string;
  title: string;
  children: ReactNode;
}
export const ErrorPage: React.FC<ErrorPageProps> = ({
  imageSrc,
  title,
  children,
}) => {
  const router = useRouter();

  return (
    <div className={styles.errorPage}>
      <Image
        alt=''
        src={imageSrc}
        fill
        priority={false}
        sizes='100vw'
        style={{
          objectFit: 'cover',
        }}
      />
      <div className={styles.content}>
        <div className='shine' />
        <h1>{title}</h1>
        <div>{children}</div>

        <Button mode='primary' onClick={() => router.push('/')}>
          Go back to main page
        </Button>
      </div>
    </div>
  );
};
