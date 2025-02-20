export function useCopy() {
  const toast = useToast();

  function copy(text: string) {
    navigator.clipboard.writeText(text);
    toast.add({
      title: "Successfully copied to clipboard!",
      description: text,
    });
  }

  return { copy };
}
