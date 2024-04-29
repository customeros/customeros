import { useEffect } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';

export const Auth = () => {
  const navigate = useNavigate();

  useEffect(() => {
    navigate('/auth/signin');
  }, []);

  return <Outlet />;
};
