import React from 'react';
import Image from 'next/image';
import { LoginPanel } from '@spaces/organisms/login-panel';

export async function getServerSideProps() {
  const backgroundImageUrlNumber = `${
    Math.floor(Math.random() * 3) + 1
  }`.padStart(2, '0');

  return {
    props: {
      image: `/backgrounds/blueprint/background-000${backgroundImageUrlNumber}`,
    },
  };
}
const Login = ({ image }: { image: string }) => {
  return (
    <>
      <Image
        alt=''
        src={`${image}.webp`}
        fill
        placeholder='blur'
        blurDataURL={`${image}-blur.webp`}
        sizes='100vw'
        style={{
          objectFit: 'cover',
        }}
      />
      <div className='flex flex-row h-full'>
        <div className='login-panel flex-grow-1'>
          <LoginPanel />
        </div>
      </div>
    </>
  );
};

export default Login;
