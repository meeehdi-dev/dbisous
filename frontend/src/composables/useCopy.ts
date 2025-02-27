export function useCopy() {
  const toast = useToast();

  async function copy(text: string) {
    await navigator.clipboard.writeText(text);
    toast.add({
      title: "Successfully copied to clipboard!",
      description: text,
    });
  }

  return { copy };
}
