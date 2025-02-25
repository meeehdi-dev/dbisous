import { Effect } from "effect";
import { TaggedError } from "effect/Data";

class WailsError extends TaggedError("WailsError")<{ message: string }> {}

export const useWails = () => {
  const toast = useToast();

  const wails = <T>(fn: () => PromiseLike<T>) =>
    Effect.tryPromise({
      try: fn,
      catch: (err: unknown) => {
        const message = typeof err === "string" ? err : (err as Error).message;
        toast.add({
          title: fn.name,
          description: message,
          color: "error",
        });
        return new WailsError({ message });
      },
    });

  return wails;
};
