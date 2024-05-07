import { useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';

const childrenPaths = ['/signin', '/success', '/failure'];

export const Auth = () => {
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (childrenPaths.some((p) => location.pathname.includes(p))) return;

    navigate('/auth/signin');
  }, []);

  return <Outlet />;
};
