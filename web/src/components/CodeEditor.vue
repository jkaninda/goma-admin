<template>
  <div class="code-editor-wrapper" ref="wrapperRef"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { EditorView, keymap, placeholder as cmPlaceholder } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { basicSetup } from 'codemirror'
import { json } from '@codemirror/lang-json'
import { yaml } from '@codemirror/lang-yaml'
import { oneDark } from '@codemirror/theme-one-dark'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    language?: 'json' | 'yaml'
    readonly?: boolean
    minHeight?: string
  }>(),
  {
    modelValue: '',
    language: 'json',
    readonly: false,
    minHeight: '200px',
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const wrapperRef = ref<HTMLElement>()
let view: EditorView | null = null
let isExternalUpdate = false

function getLanguageExtension() {
  return props.language === 'yaml' ? yaml() : json()
}

function createState(doc: string) {
  return EditorState.create({
    doc,
    extensions: [
      basicSetup,
      getLanguageExtension(),
      oneDark,
      EditorState.readOnly.of(props.readonly),
      EditorView.theme({
        '&': { minHeight: props.minHeight },
        '.cm-scroller': { overflow: 'auto' },
      }),
      EditorView.updateListener.of((update) => {
        if (update.docChanged && !isExternalUpdate) {
          emit('update:modelValue', update.state.doc.toString())
        }
      }),
    ],
  })
}

onMounted(() => {
  if (!wrapperRef.value) return

  view = new EditorView({
    state: createState(props.modelValue),
    parent: wrapperRef.value,
  })
})

onUnmounted(() => {
  view?.destroy()
  view = null
})

watch(
  () => props.modelValue,
  (newValue) => {
    if (!view) return
    const current = view.state.doc.toString()
    if (newValue !== current) {
      isExternalUpdate = true
      view.dispatch({
        changes: {
          from: 0,
          to: view.state.doc.length,
          insert: newValue,
        },
      })
      isExternalUpdate = false
    }
  }
)

watch(
  () => [props.language, props.readonly],
  () => {
    if (!view) return
    const doc = view.state.doc.toString()
    view.setState(createState(doc))
  }
)
</script>
