<script setup lang="ts">
declare global {
  interface Window {
    runtime: { BrowserOpenURL: (url: string) => void };
  }
}

import { onMounted, ref } from "vue";
import MarkdownIt from "markdown-it";

const packageVersion = import.meta.env.PACKAGE_VERSION as string;

interface Release {
  tag_name: string;
  name: string;
  body: string;
  html_url: string;
  published_at: string;
}

const releases = ref<Release[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

const md = new MarkdownIt({ linkify: true });
md.renderer.rules.link_open = function (tokens, idx, options, _env, self) {
  const aIndex = tokens[idx].attrIndex("href");
  if (aIndex < 0) {
    return self.renderToken(tokens, idx, options);
  }

  const href = tokens[idx].attrs?.[aIndex][1];
  if (!href) {
    return self.renderToken(tokens, idx, options);
  }

  tokens[idx].attrPush([
    "onclick",
    `window.runtime.BrowserOpenURL('${href}'); return false;`,
  ]);

  return self.renderToken(tokens, idx, options);
};

onMounted(async () => {
  try {
    const response = await fetch(
      "https://api.github.com/repos/meeehdi-dev/dbisous/releases",
    );
    releases.value = (await response.json()) as Release[];
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (_) {
    error.value = "Failed to fetch releases.";
  } finally {
    loading.value = false;
  }
});

function openLink(url: string) {
  window.runtime.BrowserOpenURL(url);
}
</script>

<template>
  <div class="flex flex-initial">
    <div class="flex flex-auto items-center justify-center p-2">
      <UModal :ui="{ content: 'max-w-4xl' }">
        <UButton
          icon="simple-icons:git"
          color="neutral"
          variant="soft"
          size="sm"
          :label="`v${packageVersion}`"
        />

        <template #content>
          <div class="p-4">
            <h2 class="mb-4 text-lg font-semibold">Latest Releases</h2>
            <div v-if="loading">Loading...</div>
            <div v-else-if="error" class="text-red-500">{{ error }}</div>
            <ul v-else class="max-h-96 space-y-4 overflow-y-auto">
              <li v-for="release in releases" :key="release.tag_name">
                <div class="flex items-center gap-2">
                  <h3 class="text-2xl font-medium">
                    {{ release.name }}
                  </h3>
                  <template
                    v-if="release.tag_name.replace('v', '') > packageVersion"
                  >
                    <UBadge
                      color="warning"
                      icon="lucide:link"
                      class="cursor-pointer"
                      @click="openLink(release.html_url)"
                    >
                      NEW
                    </UBadge>
                  </template>
                </div>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  Published on:
                  {{ new Date(release.published_at).toLocaleDateString() }}
                </p>
                <div
                  class="prose prose-sm dark:prose-invert mt-2"
                  v-html="md.render(release.body)"
                ></div>
                <USeparator :ui="{ root: 'mt-4' }" />
              </li>
            </ul>
          </div>
        </template>
      </UModal>
    </div>
  </div>
</template>
