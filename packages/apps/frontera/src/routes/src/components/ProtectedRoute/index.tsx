import { Navigate } from 'react-router-dom';

export const ProtectedRoute = ({
  children,
  condition,
  fallback,
}: {
  fallback: string;
  children: React.ReactNode;
  condition?: boolean | null;
}) => {
  if (!condition) {
    return <Navigate replace to={fallback} />;
  }

  return <>{children}</>;
};
