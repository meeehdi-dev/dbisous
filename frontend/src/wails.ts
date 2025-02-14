import { Effect } from "effect";

class WailsError extends Error {
  readonly _tag = "WailsError";
}

export const useWails = () => {
  const toast = useToast();

  const wails = <T>(fn: () => PromiseLike<T>) =>
    Effect.tryPromise({
      try: fn,
      catch: (err: unknown) => {
        const error = typeof err === "string" ? err : (err as Error).message;
        toast.add({
          title: fn.name,
          description: error,
          color: "error",
        });
        return new WailsError(error);
      },
    });

  return wails;
};
