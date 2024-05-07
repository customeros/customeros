import { useNavigate } from 'react-router-dom';

import { Button } from '@ui/form/Button/Button';

export const FailurePage = () => {
  const navigate = useNavigate();

  return (
    <div className='h-screen w-screen'>
      <div className='flex flex-col items-center justify-center gap-4'>
        <p>Authenthication failed.</p>
        <Button onClick={() => navigate('/auth/signin')}>Go back</Button>
      </div>
    </div>
  );
};
