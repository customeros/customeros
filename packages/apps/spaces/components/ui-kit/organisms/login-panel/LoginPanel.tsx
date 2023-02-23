import { InputText } from 'primereact/inputtext';
import React, { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import styles from './login-panel.module.scss';
import { NextRouter, useRouter } from 'next/router';
import {
  Configuration,
  FrontendApi,
  LoginFlow,
  UpdateLoginFlowBody,
} from '@ory/client';
import { AxiosError } from 'axios';
import { edgeConfig } from '@ory/integrations/next';
import { Flow } from './ui';
import Image from 'next/image';
import { Button, Input } from '../../atoms';

// interface Props {}

// A small function to help us deal with errors coming from fetching a flow.
export function handleFlowError<S>(
  router: NextRouter,
  flowType: 'login' | 'registration' | 'settings' | 'recovery' | 'verification',
  resetFlow: Dispatch<SetStateAction<S | undefined>>,
) {
  return async (err: AxiosError) => {
    console.log(err);
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //@ts-expect-error
    switch (err.response?.data.error?.id) {
      case 'session_aal2_required':
        // 2FA is enabled and enforced, but user did not perform 2fa yet!
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        //@ts-expect-error
        window.location.href = err.response?.data.redirect_browser_to;
        return;
      case 'session_already_available':
        // User is already signed in, let's redirect them home!
        await router.push('/');
        return;
      case 'session_refresh_required':
        // We need to re-authenticate to perform this action
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        //@ts-expect-error
        window.location.href = err.response?.data.redirect_browser_to;
        return;
      case 'self_service_flow_return_to_forbidden':
        // The flow expired, let's request a new one.
        toast.error('The return_to address is not allowed.');
        resetFlow(undefined);
        await router.push('/' + flowType);
        return;
      case 'self_service_flow_expired':
        // The flow expired, let's request a new one.
        toast.error(
          'Your interaction expired, please fill out the form again.',
        );
        resetFlow(undefined);
        await router.push('/' + flowType);
        return;
      case 'security_csrf_violation':
        // A CSRF violation occurred. Best to just refresh the flow!
        toast.error(
          'A security violation was detected, please fill out the form again.',
        );
        resetFlow(undefined);
        await router.push('/' + flowType);
        return;
      case 'security_identity_mismatch':
        // The requested item was intended for someone else. Let's request a new flow...
        resetFlow(undefined);
        await router.push('/' + flowType);
        return;
      case 'browser_location_change_required':
        // Ory Kratos asked us to point the user to this URL.
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        //@ts-expect-error
        window.location.href = err.response.data.redirect_browser_to;
        return;
    }

    switch (err.response?.status) {
      case 410:
        // The flow expired, let's request a new one.
        resetFlow(undefined);
        await router.push('/' + flowType);
        return;
    }

    // We are not able to handle the error? Return it.
    return Promise.reject(err);
  };
}

const ory = new FrontendApi(new Configuration(edgeConfig));

export const LoginPanel: React.FC = () => {
  const [loginForm, setLoginForm] = useState('login');

  const waitlist = () => setLoginForm('waitlist');
  const login = () => setLoginForm('login');
  const forgotPassword = () => setLoginForm('forgotPassword');

  const INIT = 'INIT';
  const SUBMITTING = 'SUBMITTING';
  const SUCCESS = 'SUCCESS';

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [forgottenPasswordEmail, setForgottenPasswordEmail] = useState('');
  const [formState, setFormState] = useState(INIT);
  const [errorMessage, setErrorMessage] = useState('');

  const SignUpFormError = (errorMessage?: string) =>
    toast.error(errorMessage || 'Oops! Something went wrong, please try again');

  const isValidEmail = (email: string) => {
    return /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(
      email,
    );
    // return true;
  };

  /**
   * Rate limit the number of submissions allowed
   * @returns {boolean} true if the form has been successfully submitted in the past minute
   */
  const hasRecentSubmission = () => {
    const time = new Date();
    const timestamp = time.valueOf();
    const previousTimestamp = localStorage.getItem('loops-form-timestamp');

    // Indicate if the last sign up was less than a minute ago
    if (
      previousTimestamp &&
      Number(previousTimestamp) + 60 * 1000 > timestamp
    ) {
      setErrorMessage('Too many signups, please try again in a little while');
      SignUpFormError(errorMessage);
      return true;
    }

    localStorage.setItem('loops-form-timestamp', timestamp.toString());
    return false;
  };

  const resetForm = () => {
    setEmail('');
    setFormState(INIT);
    setErrorMessage('');
  };

  const handleSubmit = (event: any) => {
    // Prevent the default form submission
    event.preventDefault();

    // boundary conditions for submission
    if (formState !== INIT) return;
    if (!isValidEmail(email)) {
      setErrorMessage('Please enter a valid email');
      SignUpFormError(errorMessage);
      return;
    }
    if (hasRecentSubmission()) return;
    setFormState(SUBMITTING);

    // build body
    const formBody = `userGroup=${encodeURIComponent(
      'Waitlist-Login',
    )}&email=${encodeURIComponent(email)}&firstName=${encodeURIComponent(
      firstName,
    )}&lastName=${encodeURIComponent(lastName)}`;

    // API request to add user to newsletter
    fetch(
      `https://app.loops.so/api/newsletter-form/cl7hzfqge458409jvsbqy93u9`,
      {
        method: 'POST',
        body: formBody,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    )
      .then((res: any) => {
        if (res) {
          resetForm();
          setFormState(SUCCESS);
        } else {
          setErrorMessage(res.statusText);
          localStorage.setItem('loops-form-timestamp', '');
          SignUpFormError(errorMessage);
        }
      })
      .catch((error) => {
        // check for cloudflare error
        if (error.message === 'Failed to fetch') {
          setErrorMessage(
            'Too many signups, please try again in a little while',
          );
          SignUpFormError(errorMessage);
        } else if (error.message) {
          setErrorMessage(error.message);
          SignUpFormError(errorMessage);
        }
        localStorage.setItem('loops-form-timestamp', '');
      });
  };

  const handleForgotPassword = (event: any) => {
    event.preventDefault();

    if (forgottenPasswordEmail === '' || null) {
      SignUpFormError('Please enter your email');
      return;
    }
    if (!isValidEmail(forgottenPasswordEmail)) {
      SignUpFormError('Please enter a valid email');
      return;
    } else {
      toast.success('Please check your email for a password reset link!');
      return;
    }
  };

  const router = useRouter();
  const {
    return_to: returnTo,
    flow: flowId,
    // Refresh means we want to refresh the session. This is needed, for example, when we want to update the password
    // of a user.
    refresh,
    // AAL = Authorization Assurance Level. This implies that we want to upgrade the AAL, meaning that we want
    // to perform two-factor authentication/verification.
    aal,
  } = router.query;

  const [flow, setFlow] = useState<LoginFlow>();

  useEffect(() => {
    // If the router is not ready yet, or we already have a flow, do nothing.
    if (!router.isReady || flow) {
      return;
    }

    // If ?flow=.. was in the URL, we fetch it
    if (flowId) {
      ory
        .getLoginFlow({ id: String(flowId) })
        .then(({ data }) => {
          setFlow(data);
        })
        .catch(handleFlowError(router, 'login', setFlow));
      return;
    }

    // Otherwise we initialize it
    ory
      .createBrowserLoginFlow({
        refresh: Boolean(refresh),
        aal: aal ? String(aal) : undefined,
        returnTo: returnTo ? String(returnTo) : undefined,
      })
      .then(({ data }) => {
        setFlow(data);
      })
      .catch(handleFlowError(router, 'login', setFlow));
  }, [flowId, router, router.isReady, aal, refresh, returnTo, flow]);

  const handleLogin = (values: UpdateLoginFlowBody): Promise<void> => {
    return (
      router
        // On submission, add the flow ID to the URL but do not navigate. This prevents the user loosing
        // his data when she/he reloads the page.
        .push(`/login?flow=${flow?.id}`, undefined, { shallow: true })
        .then(() =>
          ory
            .updateLoginFlow({
              flow: String(flow?.id),
              updateLoginFlowBody: values,
            })
            // We logged in successfully! Let's bring the user home.
            .then(() => {
              if (flow?.return_to) {
                window.location.href = flow?.return_to;
                return;
              }
              router.push('/');
            })
            .catch(handleFlowError(router, 'login', setFlow))
            .catch((err: AxiosError) => {
              // If the previous handler did not catch the error it's most likely a form validation error
              if (err.response?.status === 400) {
                // Yup, it is!
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                //@ts-expect-error
                setFlow(err.response?.data);
                return;
              }

              return Promise.reject(err);
            }),
        )
    );
  };

  return (
    <>
      <div className={styles.loginPanel}>
        {loginForm === 'login' && (
          <>
            <Image
              className={styles.logo}
              src='logos/openline.svg'
              alt='Openline'
              height={60}
              width={170}
            />

            <p>Don&apos;t have an account?</p>
            <Button mode='link' onClick={waitlist}>
              Join the waitlist!
            </Button>

            <Flow flow={flow} onSubmit={handleLogin} />
            <div className={styles.oryInfoSection}>
              <span
                className='font-medium line-height-3 text-sm'
                style={{ color: '#9E9E9E' }}
              >
                Protected by{' '}
              </span>
              <Image
                className={styles.oryLogo}
                src='logos/ory-small.svg'
                alt='Ory'
                height={30}
                width={30}
                style={{ verticalAlign: 'middle' }}
              />
            </div>
          </>
        )}

        {loginForm === 'waitlist' && (
          <>
            <div className='text-center mb-5'>
              <Image
                className={styles.logo}
                src='logos/openline.svg'
                alt='Openline'
                height={50}
                width={50}
              />

              <div>
                <span className='text-600 font-medium line-height-3 text-sm'>
                  Already have an account?
                </span>
                <a
                  className='font-medium no-underline ml-2 text-blue-500 cursor-pointer text-sm'
                  onClick={() => login()}
                >
                  Login now!
                </a>
              </div>
            </div>

            {formState === SUCCESS && (
              <>
                <div className='text-800 font-medium line-height-3 text-center py-8'>
                  Thanks for joining the waitlist - you should have a welcome
                  email in your inbox already!
                </div>
                <div className='pt-5 text-center'>
                  <a
                    className='font-medium no-underline ml-2 text-blue-500 cursor-pointer text-sm'
                    href='https://www.openline.ai'
                  >
                    Head back to the Openline website!
                  </a>
                </div>
              </>
            )}

            {formState === INIT && (
              <>
                <form onSubmit={handleSubmit}>
                  <label
                    htmlFor='firstName'
                    className='block text-600 font-medium mb-2 text-sm'
                  >
                    First Name
                  </label>
                  <InputText
                    type='firstName'
                    className='w-full mb-3'
                    onChange={(e) => setFirstName(e.target.value)}
                  />

                  <label
                    htmlFor='lastName'
                    className='block text-600 font-medium mb-2 text-sm'
                  >
                    Last Name
                  </label>
                  <InputText
                    type='lastName'
                    className='w-full mb-3'
                    onChange={(e) => setLastName(e.target.value)}
                  />

                  <label
                    htmlFor='email'
                    className='block text-600 font-medium mb-2 text-sm'
                  >
                    Email
                  </label>
                  <InputText
                    id='email'
                    type='text'
                    className='w-full mb-6'
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </form>
                <div className={styles.oryInfoSection}>
                  <a
                    href='https://www.openline.ai'
                    style={{ color: '#9E9E9E', textDecoration: 'none' }}
                  >
                    <span className='font-medium mr-1 cursor-pointer text-sm'>
                      Powered by
                    </span>
                    <Image
                      src='logos/openline_gray.svg'
                      alt='Ory'
                      height={30}
                      width={30}
                      style={{ verticalAlign: 'middle' }}
                    />
                  </a>
                </div>
              </>
            )}
          </>
        )}

        {loginForm === 'forgotPassword' && (
          <>
            <div className='text-center mb-5'>
              <Image
                className={styles.logo}
                src='logos/openline.svg'
                alt='Openline'
                height={50}
                width={50}
              />

              <div>
                <span className='text-600 font-medium line-height-3 text-sm'>
                  Remembered already?
                </span>
                <a
                  className='font-medium no-underline ml-2 text-blue-500 cursor-pointer text-sm'
                  onClick={() => login()}
                >
                  Login now!
                </a>
              </div>
            </div>

            <div>
              <form onSubmit={handleForgotPassword}>
                <label
                  htmlFor='email'
                  className='block text-600 font-medium mb-3 text-sm'
                >
                  Enter your email here for a password reset
                </label>
                <Input
                  id='email'
                  type='text'
                  label='Email'
                  autocomplete='username'
                  className='w-full mb-5'
                  onChange={(e) => setForgottenPasswordEmail(e.target.value)}
                />

                <Button className='w-full p-button-secondary' type='submit'>
                  Reset Password
                </Button>
              </form>

              <div className={styles.oryInfoSection}>
                <span
                  className='font-medium line-height-3 text-sm'
                  style={{ color: '#9E9E9E' }}
                >
                  Protected by{' '}
                </span>
                <Image
                  className={styles.oryLogo}
                  src='logos/ory-small.svg'
                  alt='Ory'
                  height={30}
                  width={30}
                  style={{ verticalAlign: 'middle' }}
                />
              </div>
            </div>
          </>
        )}
      </div>
    </>
  );
};
