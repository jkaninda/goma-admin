<template>
  <Modal
    :show="show"
    :title="middleware ? 'Edit Middleware' : 'New Middleware'"
    size="xl"
    @close="emit('close')"
  >
    <div class="modal-body">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">Name</label>
          <input
            v-model="form.name"
            required
            class="form-input"
            placeholder="rate-limit"
          />
        </div>
        <div class="form-group">
          <label class="form-label">Type</label>
          <input
            v-model="form.type"
            required
            class="form-input"
            placeholder="rateLimit"
          />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">Config (YAML)</label>
        <CodeEditor v-model="configYaml" language="yaml" min-height="260px" />
        <div v-if="configError" class="form-error">{{ configError }}</div>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" @click="emit('close')">
          Cancel
        </button>
        <button
          type="button"
          class="btn btn-primary"
          :disabled="saving"
          @click="handleSubmit"
        >
          {{ saving ? 'Saving...' : middleware ? 'Update' : 'Create' }}
        </button>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { middlewaresApi, type Middleware, type MiddlewareCreateRequest } from '@/api/middlewares'
import Modal from '@/components/Modal.vue'
import CodeEditor from '@/components/CodeEditor.vue'
import { parseYaml, toYaml } from '@/utils/yaml'

const props = defineProps<{
  show: boolean
  /** Middleware being edited, or null to create a new one. */
  middleware: Middleware | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved', middleware: Middleware): void
}>()

const saving = ref(false)
const configYaml = ref('')
const configError = ref('')
const form = reactive({ name: '', type: '' })

const defaultConfigYaml = `paths:
  - /api/*
rule:
  requestsPerUnit: 100
  unit: minute
  burst: 10`

// Re-initialize the form whenever the modal opens.
watch(
  () => props.show,
  (open) => {
    if (open) initForm()
  },
  { immediate: true }
)

function initForm() {
  configError.value = ''
  if (props.middleware) {
    form.name = props.middleware.name
    form.type = props.middleware.type
    configYaml.value =
      props.middleware.config && Object.keys(props.middleware.config).length > 0
        ? toYaml(props.middleware.config)
        : ''
  } else {
    form.name = ''
    form.type = ''
    configYaml.value = defaultConfigYaml
  }
}

function buildPayload(): MiddlewareCreateRequest | null {
  configError.value = ''

  if (!form.name.trim()) {
    configError.value = 'Name is required.'
    return null
  }
  if (!form.type.trim()) {
    configError.value = 'Type is required.'
    return null
  }

  let config: Record<string, unknown> = {}
  try {
    if (configYaml.value.trim()) {
      config = parseYaml(configYaml.value)
    }
  } catch {
    configError.value = 'Invalid YAML in config. Please check your syntax.'
    return null
  }

  return { name: form.name.trim(), type: form.type.trim(), config }
}

async function handleSubmit() {
  const data = buildPayload()
  if (!data) return

  saving.value = true
  try {
    const res = props.middleware
      ? await middlewaresApi.update(props.middleware.id, data)
      : await middlewaresApi.create(data)
    emit('saved', res.data)
    emit('close')
  } catch (e: any) {
    configError.value = e.response?.data?.message || 'Failed to save middleware.'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 20px;
}
</style>
