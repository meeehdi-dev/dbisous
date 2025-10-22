<script setup vapor lang="ts">
import { ref } from "vue";
import Button from "@/components/button.vue";
import { DbService } from "../../bindings/dbisous/app";

const tmp = ref("");

async function test() {
  tmp.value = await DbService.GetDb();
}
</script>

<template>
  <Button label="Go to about" to="/about" />
  <Button label="test" @click="test()" />
  <span>{{ tmp }}</span>
</template>
