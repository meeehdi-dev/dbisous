import { Effect } from "effect";

export const useWails = () => {
  const toast = useToast();

  const wails = <T>(fn: () => PromiseLike<T>) =>
    Effect.tryPromise({
      try: fn,
      catch: (err: unknown): Effect.Effect<void, string> => {
        const error = typeof err === "string" ? err : (err as Error).message;
        toast.add({
          title: fn.name,
          description: err as string,
          color: "error",
        });
        return Effect.fail(error);
      },
    });

  return wails;
};
