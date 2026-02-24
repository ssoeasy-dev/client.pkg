import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@ssoeasy-dev/react';

export const LoginPage = () => {
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    auth.handleRedirectCallback()
      .then(({ redirectTo }) => {
        navigate(redirectTo, { replace: true });
      })
      .catch((err) => {
        console.error('Authentication failed:', err);
        navigate('/');
      });
  }, [auth, navigate]);

  return <div>Processing login...</div>;
};
