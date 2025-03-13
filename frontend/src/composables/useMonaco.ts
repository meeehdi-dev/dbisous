import { editor } from "monaco-editor";
import { Ref } from "vue";

export const useMonaco = () => {
  function create(el: HTMLDivElement, value: Ref<string>) {
    const e = editor.create(el, {
      language: "sql",
      minimap: { enabled: false },
      lineNumbers: "off",
      lineDecorationsWidth: 0,
      folding: false,
      contextmenu: false,
      value: value.value,
      theme: "vs-dark",
      scrollBeyondLastLine: false,
      wordWrap: "on",
    });
    e.onEndUpdate(() => {
      value.value = e.getValue();
    });

    return e;
  }

  return { create };
};
