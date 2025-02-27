class WailsError extends Error {}

export const useWails = () => {
  const toast = useToast();

  async function wails<T>(fn: () => PromiseLike<T>) {
    try {
      const res = await fn();
      return res;
    } catch (err) {
      const message = typeof err === "string" ? err : (err as Error).message;
      toast.add({
        title: fn.name,
        description: message,
        color: "error",
      });
      return new WailsError(message);
    }
  }

  return wails;
};
