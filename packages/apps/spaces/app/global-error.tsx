'use client';

export default function GlobalError({
  error,
  reset,
}: {
  reset: () => void;
  error: Error & { digest?: string };
}) {
  return (
    <html>
      <body>
        <h2>Something went wrong!</h2>
        <button onClick={() => reset()}>Try again</button>
      </body>
    </html>
  );
}
