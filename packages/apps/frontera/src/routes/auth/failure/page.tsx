import { useNavigate } from 'react-router-dom';

import { Button } from '@ui/form/Button/Button';

export const FailurePage = () => {
  const navigate = useNavigate();
  const params = new URLSearchParams(window.location.search);
  const message = params.get('message');

  return (
    <div className='flex items-center justify-center h-screen w-screen'>
      <div className='flex flex-col items-center justify-center gap-4'>
        <p>Authenthication failed.</p>
        {message && <p>{message}</p>}
        <Button onClick={() => navigate('/auth/signin')}>Go back</Button>
      </div>
    </div>
  );
};
