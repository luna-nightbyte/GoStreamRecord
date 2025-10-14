<template>
  <form @submit.prevent="save" class="col" style="gap:.75rem;">
    <div class="col">
      <label for="name">Name</label>
      <input class="rdp_input" id="name" v-model.trim="draft.name" placeholder="EXAMPLE_CLIENT" required />
    </div>
    <div class="col">
      <label for="id">ID (unique)</label>
      <input class="rdp_input" id="id" v-model.trim="draft.id" placeholder="20250101-xxxx" />
    </div>
    <div class="col">
      <label for="url">URL</label>
      <input class="rdp_input" id="url" v-model.trim="draft.url" placeholder="http://CLIENT-IP:PORT" required />
    </div>
    <div class="row" style="gap:.5rem; margin-top:.5rem;">
      <button class="btn" type="submit">Save</button>
      <button class="btn ghost" type="button" @click="$emit('cancel')">Cancel</button>
      <div class="spacer"></div>
      <slot name="danger"></slot>
    </div>
  </form>
</template>

<script setup>
import { defineEmits,defineProps,reactive, watchEffect } from 'vue'

const props = defineProps({ modelValue: { type: Object, required: true } })
const emit = defineEmits(['update:modelValue','save','cancel'])

const draft = reactive({ ...props.modelValue })
watchEffect(()=> emit('update:modelValue', draft))

function save(){ emit('save') }
</script>
