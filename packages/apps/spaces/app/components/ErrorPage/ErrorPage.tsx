import { PropsWithChildren } from 'react';
import Image from 'next/image';

import { LinkButton } from '@spaces/ui/form/LinkButton/LinkButton';

import '../../../styles/overwrite.scss';
import '../../../styles/normalization.scss';
import '../../../styles/theme.css';
import '../../../styles/globals.scss';

import styles from './ErrorPage.module.scss';

interface ErrorPageProps {
  imageSrc: string;
  blurredSrc: string;
  title: string;
}
export const ErrorPage = ({
  imageSrc,
  blurredSrc,
  title,
  children,
}: PropsWithChildren<ErrorPageProps>) => {
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

        <LinkButton href='/' colorScheme='primary'>
          Go back to main page
        </LinkButton>
      </div>
    </div>
  );
};
