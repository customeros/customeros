export type UnwrapOptions<E = unknown> = {
  defaultValue?: unknown; // Default value to return on error
  onError?: (error: E) => void; // Hook to log or handle errors
};

export type UnwrapResult<T, E = unknown> = [T | null, E | null];

export function unwrap<T, E = unknown>(
  promise: Promise<T>,
  options: UnwrapOptions<E> = {},
): Promise<[T | null, E | null]> {
  const { onError, defaultValue } = options;

  return promise
    .then((data) => [data, null] as [T, null])
    .catch((error: E) => {
      if (onError) {
        onError(error);
      }

      return [defaultValue || null, error] as [T | null, E];
    });
}
