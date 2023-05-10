import React, { ReactNode } from 'react';
import Image from 'next/image';
import { Button } from '@spaces/atoms/button/Button';
import { useRouter } from 'next/router';
import styles from './error-page.module.scss';

interface ErrorPageProps {
  imageSrc: string;
  blurredSrc: string;
  title: string;
  children: ReactNode;
}
export const ErrorPage: React.FC<ErrorPageProps> = ({
  imageSrc,
  blurredSrc,
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
        placeholder='blur'
        blurDataURL={blurredSrc}
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
